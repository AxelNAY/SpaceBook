package handlers

import (
	"net/http"
	"time"

	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetUserReservations(c echo.Context) error {
	userId := c.QueryParam("userId")

	if userId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "userId est requis",
		})
	}

	var reservations []models.Reservation

	if err := config.DB.
		Preload("Resource").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&reservations).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la récupération des réservations",
		})
	}

	return c.JSON(http.StatusOK, reservations)
}

func CreateReservation(c echo.Context) error {
	var reservation models.Reservation

	if err := c.Bind(&reservation); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Corps de requête invalide",
		})
	}

	// Validation des dates
	if reservation.StartAt.After(reservation.EndAt) || reservation.StartAt.Equal(reservation.EndAt) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "La date de début doit être antérieure à la date de fin",
		})
	}

	// Récupérer la ressource pour connaître sa capacité
	var resource models.Resource
	if err := config.DB.First(&resource, "id = ?", reservation.ResourceID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Ressource introuvable",
		})
	}

	// Compter les réservations qui chevauchent ce créneau (non rejetées)
	var overlappingCount int64
	if err := config.DB.Model(&models.Reservation{}).
		Where("resource_id = ?", reservation.ResourceID).
		Where("status != ?", "rejected").
		Where("start_at < ? AND end_at > ?", reservation.EndAt, reservation.StartAt).
		Count(&overlappingCount).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la vérification de disponibilité",
		})
	}

	// Vérifier si la capacité est atteinte
	if int(overlappingCount) >= resource.Capacity {
		return c.JSON(http.StatusConflict, echo.Map{
			"error":     "Ressource complète pour ce créneau horaire",
			"capacity":  resource.Capacity,
			"booked":    overlappingCount,
			"available": 0,
		})
	}

	reservation.ID = uuid.New()
	reservation.Status = "pending"

	if err := config.DB.Create(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la création de la réservation",
		})
	}

	// Récupérer l'utilisateur pour le message de notification
	var user models.User
	config.DB.First(&user, "id = ?", reservation.UserID)

	// Notification pour les admins (sans UserID = visible par tous les admins)
	notification := models.Notification{
		Type:    "reservation",
		Message: "Nouvelle demande de réservation de " + user.Username + " pour " + resource.Name,
		IsRead:  false,
	}
	config.DB.Create(&notification)

	return c.JSON(http.StatusCreated, reservation)
}


/*
GET /admin/reservations
Admin only – list all reservations with User + Resource
*/
func GetAdminReservations(c echo.Context) error {
	var reservations []models.Reservation

	if err := config.DB.
		Preload("User").
		Preload("Resource").
		Order("created_at DESC").
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
		Preload("Resource").
		First(&reservation, "id = ?", reservationID).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Réservation introuvable",
		})
	}

	// Update status
	reservation.Status = "approved"
	reservation.UpdatedAt = time.Now()

	if err := config.DB.Save(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de l'approbation de la réservation",
		})
	}

	// Notification for user (UUID pointer)
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
(Admin optional)
*/
func RejectReservation(c echo.Context) error {
	idParam := c.Param("id")

	reservationID, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "ID de réservation invalide",
		})
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, "id = ?", reservationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Réservation introuvable",
		})
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
