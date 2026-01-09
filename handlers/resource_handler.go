package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetResources(c echo.Context) error {
	var resources []models.Resource
	config.DB.Find(&resources)
	return c.JSON(http.StatusOK, resources)
}

func CreateResource(c echo.Context) error {
	var resource models.Resource
	if err := c.Bind(&resource); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid payload",
		})
	}

	if resource.Type == "room" {
		resource.Capacity = 1
		resource.Category = "none"
	}
	config.DB.Create(&resource)
	notification := models.Notification{
		Type:    "resource",
		Message: "A new resource has been created",
	}
	config.DB.Create(&notification)
	return c.JSON(http.StatusCreated, resource)
}

func DeleteResource(c echo.Context) error {
	id := c.Param("id")

	var count int64
	config.DB.Model(&models.Reservation{}).
		Where("resource_id = ?", id).
		Count(&count)

	if count > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Resource cannot be deleted because it has reservations",
		})
	}

	if err := config.DB.Delete(&models.Resource{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to delete resource",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
