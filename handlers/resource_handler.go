package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetResources(c echo.Context) error {
	placeID := c.QueryParam("place_id")

	var resources []models.Resource
	query := config.DB.Preload("Place").Preload("Category")
	if placeID != "" {
		query = query.Where("place_id = ?", placeID)
	}
	query.Find(&resources)
	return c.JSON(http.StatusOK, resources)
}

func CreateResource(c echo.Context) error {
	var resource models.Resource
	if err := c.Bind(&resource); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Create(&resource).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Échec de la création de la ressource",
		})
	}

	notification := models.Notification{
		Type:    "resource",
		Message: "Une nouvelle ressource a été créée",
	}
	config.DB.Create(&notification)
	return c.JSON(http.StatusCreated, resource)
}

func UpdateResource(c echo.Context) error {
	id := c.Param("id")

	var resource models.Resource
	if err := config.DB.First(&resource, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Ressource introuvable",
		})
	}

	if err := c.Bind(&resource); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Save(&resource).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la mise à jour de la ressource",
		})
	}

	return c.JSON(http.StatusOK, resource)
}

func DeleteResource(c echo.Context) error {
	id := c.Param("id")

	var count int64
	config.DB.Model(&models.Reservation{}).
		Where("ressource_id = ?", id).
		Count(&count)

	if count > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "La ressource ne peut pas être supprimée car elle a des réservations",
		})
	}

	if err := config.DB.Delete(&models.Resource{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression de la ressource",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
