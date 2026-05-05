package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func GetCompanies(c echo.Context) error {
	var companies []models.Company
	config.DB.Find(&companies)
	return c.JSON(http.StatusOK, companies)
}

func CreateCompany(c echo.Context) error {
	var company models.Company
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Create(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la création de l'entreprise",
		})
	}

	return c.JSON(http.StatusCreated, company)
}

func UpdateCompany(c echo.Context) error {
	id := c.Param("id")

	var company models.Company
	if err := config.DB.First(&company, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Entreprise introuvable",
		})
	}

	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Données invalides",
		})
	}

	if err := config.DB.Save(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la mise à jour de l'entreprise",
		})
	}

	return c.JSON(http.StatusOK, company)
}

func DeleteCompany(c echo.Context) error {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Company{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression de l'entreprise",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
