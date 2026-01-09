package handlers

import (
	"net/http"
	"time"

	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateReservation(c echo.Context) error {
	var reservation models.Reservation

	if err := c.Bind(&reservation); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request body",
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
			"error": "Resource not found",
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
			"error": "Failed to check availability",
		})
	}

	// Vérifier si la capacité est atteinte
	if int(overlappingCount) >= resource.Capacity {
		return c.JSON(http.StatusConflict, echo.Map{
			"error":     "Resource fully booked for this time slot",
			"capacity":  resource.Capacity,
			"booked":    overlappingCount,
			"available": 0,
		})
	}

	reservation.ID = uuid.New()
	reservation.Status = "pending"

	if err := config.DB.Create(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to create reservation",
		})
	}

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
			"error": "Failed to fetch reservations",
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
			"error": "Invalid reservation ID",
		})
	}

	var reservation models.Reservation
	if err := config.DB.
		Preload("User").
		Preload("Resource").
		First(&reservation, "id = ?", reservationID).Error; err != nil {

		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Reservation not found",
		})
	}

	// Update status
	reservation.Status = "approved"
	reservation.UpdatedAt = time.Now()

	if err := config.DB.Save(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to approve reservation",
		})
	}

	// Notification for user (UUID pointer)
	userID := reservation.UserID

	notification := models.Notification{
		UserID:  &userID,
		Type:    "reservation",
		Message: "Your reservation has been approved",
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
			"error": "Invalid reservation ID",
		})
	}

	var reservation models.Reservation
	if err := config.DB.First(&reservation, "id = ?", reservationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Reservation not found",
		})
	}

	reservation.Status = "rejected"
	reservation.UpdatedAt = time.Now()

	if err := config.DB.Save(&reservation).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to reject reservation",
		})
	}

	userID := reservation.UserID

	notification := models.Notification{
		UserID:  &userID,
		Type:    "reservation",
		Message: "Your reservation has been rejected",
		IsRead:  false,
	}

	_ = config.DB.Create(&notification)

	return c.JSON(http.StatusOK, reservation)
}
