package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetCategories(c echo.Context) error {
	var categories []models.Category
	config.DB.Find(&categories)
	return c.JSON(http.StatusOK, categories)
}

func CreateCategory(c echo.Context) error {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Create(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la création de la catégorie",
		})
	}

	return c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c echo.Context) error {
	id := c.Param("id")

	var category models.Category
	if err := config.DB.First(&category, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Catégorie introuvable",
		})
	}

	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Save(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la mise à jour de la catégorie",
		})
	}

	return c.JSON(http.StatusOK, category)
}

func DeleteCategory(c echo.Context) error {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression de la catégorie",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
