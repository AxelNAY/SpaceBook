package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetPlaces(c echo.Context) error {
	companyID := c.QueryParam("company_id")

	var places []models.Place
	query := config.DB.Preload("Company")
	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}
	query.Find(&places)
	return c.JSON(http.StatusOK, places)
}

func CreatePlace(c echo.Context) error {
	var place models.Place
	if err := c.Bind(&place); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Create(&place).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la création du lieu",
		})
	}

	return c.JSON(http.StatusCreated, place)
}

func UpdatePlace(c echo.Context) error {
	id := c.Param("id")

	var place models.Place
	if err := config.DB.First(&place, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Lieu introuvable",
		})
	}

	if err := c.Bind(&place); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Save(&place).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la mise à jour du lieu",
		})
	}

	return c.JSON(http.StatusOK, place)
}

func DeletePlace(c echo.Context) error {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Place{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression du lieu",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
