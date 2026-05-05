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
		Preload("Ressource").
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
	if reservation.StartDatetime.After(reservation.EndDatetime) || reservation.StartDatetime.Equal(reservation.EndDatetime) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "La date de début doit être antérieure à la date de fin",
		})
	}

	// Vérifier que la ressource existe
	var resource models.Resource
	if err := config.DB.First(&resource, "id = ?", reservation.RessourceID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Ressource introuvable",
		})
	}

	// Vérifier s'il y a des réservations qui chevauchent ce créneau (non rejetées)
	var overlappingCount int64
	if err := config.DB.Model(&models.Reservation{}).
		Where("ressource_id = ?", reservation.RessourceID).
		Where("status != ?", "rejected").
		Where("start_datetime < ? AND end_datetime > ?", reservation.EndDatetime, reservation.StartDatetime).
		Count(&overlappingCount).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la vérification de disponibilité",
		})
	}

	if overlappingCount > 0 {
		return c.JSON(http.StatusConflict, echo.Map{
			"error": "Ressource déjà réservée pour ce créneau horaire",
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
		Preload("Ressource").
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
		Preload("Ressource").
		First(&reservation, "id = ?", reservationID).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Réservation introuvable",
		})
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
