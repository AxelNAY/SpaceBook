package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetAdminNotifications(c echo.Context) error {
	var notifications []models.Notification

	// Latest notifications first
	config.DB.Order("created_at desc").Find(&notifications)

	return c.JSON(http.StatusOK, notifications)
}

func MarkNotificationAsRead(c echo.Context) error {
	id := c.Param("id")

	result := config.DB.Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true)

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Notification introuvable",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Notification marquée comme lue",
	})
}
