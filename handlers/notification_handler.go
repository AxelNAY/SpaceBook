package handlers

import (
	"net/http"

	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetUserNotifications(c echo.Context) error {
	userId := c.QueryParam("userId")

	if userId == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "userId is required",
		})
	}

	var notifications []models.Notification

	if err := config.DB.
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch notifications",
		})
	}

	return c.JSON(http.StatusOK, notifications)
}
