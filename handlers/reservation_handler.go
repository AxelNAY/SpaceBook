package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetUserReservations(c echo.Context) error {
	userId := c.Get("user_id")
	if userId == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Non authentifié"})
	}

	var reservations []models.Reservation

	if err := config.DB.
		Preload("Ressource.Place").
		Preload("ReservationResources.Resource.Place").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la récupération des réservations",
		})
	}

	return c.JSON(http.StatusOK, reservations)
}

type simpleReservationRequest struct {
	RessourceID   uuid.UUID `json:"ressource_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}

func notifyPlaceAdmins(resourceID uuid.UUID, message string) {
	var resource models.Resource
	if err := config.DB.Preload("Place").First(&resource, "id = ?", resourceID).Error; err != nil {
		return
	}
	if resource.Place == nil {
		return
	}
	var admins []models.User
	config.DB.
		Joins("JOIN place_users pu ON pu.user_id = users.id").
		Joins("JOIN places p ON p.id = pu.place_id").
		Where("p.company_id = ? AND users.role = ?", resource.Place.CompanyID, "admin").
		Find(&admins)
	for _, admin := range admins {
		adminID := admin.ID
		config.DB.Create(&models.Notification{
			UserID:  &adminID,
			Type:    "reservation",
			Message: message,
			IsRead:  false,
		})
	}
}

func CreateSimpleReservation(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Non authentifié"})
	}

	var req simpleReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Corps de requête invalide"})
	}

	if req.RessourceID == uuid.Nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ressource_id est requis"})
	}

	if !req.StartDatetime.Before(req.EndDatetime) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "La date de début doit être antérieure à la date de fin"})
	}

	var resource models.Resource
	if err := config.DB.First(&resource, "id = ?", req.RessourceID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Ressource introuvable"})
	}

	var simpleOverlap int64
	config.DB.Model(&models.Reservation{}).
		Where("ressource_id = ? AND status != ? AND start_datetime < ? AND end_datetime > ?",
			req.RessourceID, "rejected", req.EndDatetime, req.StartDatetime).
		Count(&simpleOverlap)

	var groupOverlap int64
	config.DB.Model(&models.Reservation{}).
		Joins("JOIN reservation_resources rr ON rr.reservation_id = reservations.id").
		Where("rr.resource_id = ? AND reservations.status != ? AND reservations.start_datetime < ? AND reservations.end_datetime > ?",
			req.RessourceID, "rejected", req.EndDatetime, req.StartDatetime).
		Count(&groupOverlap)

	if simpleOverlap+groupOverlap > 0 {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Ressource déjà réservée pour ce créneau horaire"})
	}

	resID := req.RessourceID
	reservation := models.Reservation{
		ID:            uuid.New(),
		UserID:        userID,
		RessourceID:   &resID,
		StartDatetime: req.StartDatetime,
		EndDatetime:   req.EndDatetime,
		Status:        "pending",
	}

	if err := config.DB.Create(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la création de la réservation"})
	}

	var user models.User
	config.DB.First(&user, "id = ?", userID)
	notifyPlaceAdmins(req.RessourceID, "Nouvelle demande de réservation de "+user.Username+" pour "+resource.Name)

	return c.JSON(http.StatusCreated, reservation)
}

type groupReservationRequest struct {
	ResourceIDs   []uuid.UUID `json:"resource_ids"`
	StartDatetime time.Time   `json:"start_datetime"`
	EndDatetime   time.Time   `json:"end_datetime"`
}

func CreateGroupReservation(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Non authentifié"})
	}

	var req groupReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Corps de requête invalide"})
	}

	if len(req.ResourceIDs) < 2 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Au moins deux ressources sont requises"})
	}

	if !req.StartDatetime.Before(req.EndDatetime) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "La date de début doit être antérieure à la date de fin"})
	}

	var reservation models.Reservation
	statusCode := http.StatusInternalServerError

	txErr := config.DB.Transaction(func(tx *gorm.DB) error {
		var resourceNames []string

		for _, resourceID := range req.ResourceIDs {
			var resource models.Resource
			if err := tx.First(&resource, "id = ?", resourceID).Error; err != nil {
				statusCode = http.StatusNotFound
				return fmt.Errorf("ressource %s introuvable", resourceID)
			}
			resourceNames = append(resourceNames, resource.Name)

			var simpleOverlap, groupOverlap int64
			tx.Model(&models.Reservation{}).
				Where("ressource_id = ? AND status != ? AND start_datetime < ? AND end_datetime > ?",
					resourceID, "rejected", req.EndDatetime, req.StartDatetime).
				Count(&simpleOverlap)
			tx.Model(&models.Reservation{}).
				Joins("JOIN reservation_resources rr ON rr.reservation_id = reservations.id").
				Where("rr.resource_id = ? AND reservations.status != ? AND reservations.start_datetime < ? AND reservations.end_datetime > ?",
					resourceID, "rejected", req.EndDatetime, req.StartDatetime).
				Count(&groupOverlap)

			if simpleOverlap+groupOverlap > 0 {
				statusCode = http.StatusConflict
				return fmt.Errorf("ressource %s déjà réservée pour ce créneau", resource.Name)
			}
		}

		reservation = models.Reservation{
			ID:            uuid.New(),
			UserID:        userID,
			StartDatetime: req.StartDatetime,
			EndDatetime:   req.EndDatetime,
			Status:        "pending",
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		for _, resourceID := range req.ResourceIDs {
			rid := resourceID
			if err := tx.Create(&models.ReservationResource{
				ID:            uuid.New(),
				ReservationID: reservation.ID,
				ResourceID:    rid,
				Status:        "pending",
			}).Error; err != nil {
				return err
			}
		}

		var user models.User
		tx.First(&user, "id = ?", userID)
		msg := "Nouvelle réservation de groupe de " + user.Username + " (" + strings.Join(resourceNames, ", ") + ")"

		// Collecter les admins uniques de tous les lieux concernés
		seen := map[uuid.UUID]bool{}
		for _, resourceID := range req.ResourceIDs {
			var res models.Resource
			if tx.Preload("Place").First(&res, "id = ?", resourceID).Error != nil || res.Place == nil {
				continue
			}
			var admins []models.User
			tx.Joins("JOIN place_users pu ON pu.user_id = users.id").
				Joins("JOIN places p ON p.id = pu.place_id").
				Where("p.company_id = ? AND users.role = ?", res.Place.CompanyID, "admin").
				Find(&admins)
			for _, admin := range admins {
				if !seen[admin.ID] {
					seen[admin.ID] = true
					adminID := admin.ID
					tx.Create(&models.Notification{
						UserID:  &adminID,
						Type:    "reservation",
						Message: msg,
						IsRead:  false,
					})
				}
			}
		}

		return nil
	})

	if txErr != nil {
		return c.JSON(statusCode, echo.Map{"error": txErr.Error()})
	}

	return c.JSON(http.StatusCreated, reservation)
}

/*
PUT /admin/reservations/:id/resources/approve
Admin only – approuve une ou plusieurs ressources d'une réservation de groupe
*/
func ApproveGroupReservationResources(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	reservationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID de réservation invalide"})
	}

	var req struct {
		ResourceIDs []uuid.UUID `json:"resource_ids"`
	}
	if err := c.Bind(&req); err != nil || len(req.ResourceIDs) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Liste de ressources invalide ou vide"})
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, "id = ?", reservationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Réservation introuvable"})
	}

	if reservation.RessourceID != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Cette réservation n'est pas une réservation de groupe"})
	}

	statusCode := http.StatusInternalServerError

	txErr := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, resourceID := range req.ResourceIDs {
			var resource models.Resource
			if err := tx.Preload("Place").First(&resource, "id = ?", resourceID).Error; err != nil {
				statusCode = http.StatusNotFound
				return fmt.Errorf("ressource %s introuvable", resourceID)
			}
			if resource.Place == nil || resource.Place.CompanyID.String() != companyIDStr {
				statusCode = http.StatusForbidden
				return fmt.Errorf("accès refusé à la ressource %s", resource.Name)
			}

			result := tx.Model(&models.ReservationResource{}).
				Where("reservation_id = ? AND resource_id = ?", reservationID, resourceID).
				Update("status", "approved")
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				statusCode = http.StatusNotFound
				return fmt.Errorf("ressource %s non trouvée dans cette réservation", resource.Name)
			}
		}

		var totalCount, approvedCount int64
		tx.Model(&models.ReservationResource{}).Where("reservation_id = ?", reservationID).Count(&totalCount)
		tx.Model(&models.ReservationResource{}).Where("reservation_id = ? AND status = ?", reservationID, "approved").Count(&approvedCount)

		if totalCount > 0 && totalCount == approvedCount {
			if err := tx.Model(&reservation).Update("status", "approved").Error; err != nil {
				return err
			}
		}

		userID := reservation.UserID
		notification := models.Notification{
			UserID:  &userID,
			Type:    "reservation",
			Message: fmt.Sprintf("%d ressource(s) de votre réservation de groupe ont été approuvées", len(req.ResourceIDs)),
			IsRead:  false,
		}
		tx.Create(&notification)

		return nil
	})

	if txErr != nil {
		return c.JSON(statusCode, echo.Map{"error": txErr.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ressources approuvées avec succès"})
}

/*
PUT /admin/reservations/:id/resources/reject
Admin only – rejette une ou plusieurs ressources d'une réservation de groupe
*/
func RejectGroupReservationResources(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	reservationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ID de réservation invalide"})
	}

	var req struct {
		ResourceIDs []uuid.UUID `json:"resource_ids"`
	}
	if err := c.Bind(&req); err != nil || len(req.ResourceIDs) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Liste de ressources invalide ou vide"})
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, "id = ?", reservationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Réservation introuvable"})
	}

	if reservation.RessourceID != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Cette réservation n'est pas une réservation de groupe"})
	}

	statusCode := http.StatusInternalServerError

	txErr := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, resourceID := range req.ResourceIDs {
			var resource models.Resource
			if err := tx.Preload("Place").First(&resource, "id = ?", resourceID).Error; err != nil {
				statusCode = http.StatusNotFound
				return fmt.Errorf("ressource %s introuvable", resourceID)
			}
			if resource.Place == nil || resource.Place.CompanyID.String() != companyIDStr {
				statusCode = http.StatusForbidden
				return fmt.Errorf("accès refusé à la ressource %s", resource.Name)
			}

			result := tx.Model(&models.ReservationResource{}).
				Where("reservation_id = ? AND resource_id = ?", reservationID, resourceID).
				Update("status", "rejected")
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				statusCode = http.StatusNotFound
				return fmt.Errorf("ressource %s non trouvée dans cette réservation", resource.Name)
			}
		}

		var totalCount, rejectedCount int64
		tx.Model(&models.ReservationResource{}).Where("reservation_id = ?", reservationID).Count(&totalCount)
		tx.Model(&models.ReservationResource{}).Where("reservation_id = ? AND status = ?", reservationID, "rejected").Count(&rejectedCount)

		if totalCount > 0 && totalCount == rejectedCount {
			if err := tx.Model(&reservation).Update("status", "rejected").Error; err != nil {
				return err
			}
		}

		userID := reservation.UserID
		notification := models.Notification{
			UserID:  &userID,
			Type:    "reservation",
			Message: fmt.Sprintf("%d ressource(s) de votre réservation de groupe ont été refusées", len(req.ResourceIDs)),
			IsRead:  false,
		}
		tx.Create(&notification)

		return nil
	})

	if txErr != nil {
		return c.JSON(statusCode, echo.Map{"error": txErr.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Ressources rejetées avec succès"})
}

/*
GET /admin/reservations
Admin only – liste les réservations des ressources de son entreprise
*/
func GetAdminReservations(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Administrateur non associé à une entreprise",
		})
	}

	var reservations []models.Reservation
	if err := config.DB.
		Preload("User").
		Preload("Ressource.Place").
		Preload("ReservationResources.Resource.Place").
		Joins("LEFT JOIN resources res ON res.id = reservations.ressource_id").
		Joins("LEFT JOIN places p ON p.id = res.place_id").
		Joins("LEFT JOIN reservation_resources rr ON rr.reservation_id = reservations.id").
		Joins("LEFT JOIN resources res2 ON res2.id = rr.resource_id").
		Joins("LEFT JOIN places p2 ON p2.id = res2.place_id").
		Where("p.company_id = ? OR p2.company_id = ?", companyIDStr, companyIDStr).
		Distinct("reservations.*").
		Order("reservations.created_at DESC").
		Find(&reservations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la récupération des réservations",
		})
	}

	return c.JSON(http.StatusOK, reservations)
}

/*
PUT /admin/reservations/:id/approve
Admin only – approve reservation + notify user
*/
func ApproveReservation(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	idParam := c.Param("id")

	reservationID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "ID de réservation invalide",
		})
	}

	var reservation models.Reservation
	if err := config.DB.
		Preload("User").
		Preload("Ressource.Place").
		First(&reservation, "id = ?", reservationID).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Réservation introuvable",
		})
	}

	if reservation.Ressource == nil || reservation.Ressource.Place == nil || reservation.Ressource.Place.CompanyID.String() != companyIDStr {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à cette réservation"})
	}

	reservation.Status = "approved"
	reservation.UpdatedAt = time.Now()

	if err := config.DB.Save(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de l'approbation de la réservation",
		})
	}

	userID := reservation.UserID

	notification := models.Notification{
		UserID:  &userID,
		Type:    "reservation",
		Message: "Votre réservation a été approuvée",
		IsRead:  false,
	}

	_ = config.DB.Create(&notification)

	return c.JSON(http.StatusOK, reservation)
}

/*
PUT /admin/reservations/:id/reject
*/
func RejectReservation(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	idParam := c.Param("id")

	reservationID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "ID de réservation invalide",
		})
	}

	var reservation models.Reservation
	if err := config.DB.Preload("Ressource.Place").First(&reservation, "id = ?", reservationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Réservation introuvable",
		})
	}

	if reservation.Ressource == nil || reservation.Ressource.Place == nil || reservation.Ressource.Place.CompanyID.String() != companyIDStr {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à cette réservation"})
	}

	reservation.Status = "rejected"
	reservation.UpdatedAt = time.Now()

	if err := config.DB.Save(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec du rejet de la réservation",
		})
	}

	userID := reservation.UserID

	notification := models.Notification{
		UserID:  &userID,
		Type:    "reservation",
		Message: "Votre réservation a été refusée",
		IsRead:  false,
	}

	_ = config.DB.Create(&notification)

	return c.JSON(http.StatusOK, reservation)
}
