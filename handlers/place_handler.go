package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
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

	config.DB.Preload("Company").First(&place, "id = ?", place.ID)
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

func getAdminCompanyID(c echo.Context) (uuid.UUID, error) {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusForbidden, "Administrateur non associé à une entreprise")
	}
	return uuid.Parse(companyIDStr)
}

type PlaceWithDetails struct {
	models.Place
	Users     []models.User     `json:"users"`
	Resources []models.Resource `json:"resources"`
}

func buildPlaceDetails(place models.Place) PlaceWithDetails {
	var placeUsers []models.PlaceUser
	config.DB.Preload("User").Where("place_id = ?", place.ID).Find(&placeUsers)
	users := make([]models.User, 0, len(placeUsers))
	for _, pu := range placeUsers {
		users = append(users, pu.User)
	}

	var resources []models.Resource
	config.DB.Where("place_id = ?", place.ID).Find(&resources)

	return PlaceWithDetails{Place: place, Users: users, Resources: resources}
}

func AdminGetPlaces(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	var places []models.Place
	if err := config.DB.Preload("Company").Where("company_id = ?", companyIDStr).Find(&places).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la récupération des lieux"})
	}

	result := make([]PlaceWithDetails, 0, len(places))
	for _, p := range places {
		result = append(result, buildPlaceDetails(p))
	}
	return c.JSON(http.StatusOK, result)
}

func AdminGetPlaceDetails(c echo.Context) error {
	companyID, err := getAdminCompanyID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")
	var place models.Place
	if err := config.DB.Preload("Company").First(&place, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Lieu introuvable"})
	}
	if place.CompanyID != companyID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à ce lieu"})
	}

	return c.JSON(http.StatusOK, buildPlaceDetails(place))
}

func AdminCreatePlace(c echo.Context) error {
	companyID, err := getAdminCompanyID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	var place models.Place
	if err := c.Bind(&place); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Données invalides"})
	}
	place.CompanyID = companyID

	if err := config.DB.Create(&place).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la création du lieu"})
	}

	// Affecter automatiquement l'admin au nouveau lieu
	adminUserID, _ := c.Get("user_id").(uuid.UUID)
	if adminUserID != (uuid.UUID{}) {
		placeUser := models.PlaceUser{
			ID:      uuid.New(),
			PlaceID: place.ID,
			UserID:  adminUserID,
		}
		config.DB.Create(&placeUser)
	}

	config.DB.Preload("Company").First(&place, "id = ?", place.ID)
	return c.JSON(http.StatusCreated, place)
}

func AdminUpdatePlace(c echo.Context) error {
	companyID, err := getAdminCompanyID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")
	var place models.Place
	if err := config.DB.First(&place, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Lieu introuvable"})
	}
	if place.CompanyID != companyID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à ce lieu"})
	}

	if err := c.Bind(&place); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Données invalides"})
	}
	place.CompanyID = companyID

	if err := config.DB.Save(&place).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la mise à jour du lieu"})
	}
	return c.JSON(http.StatusOK, place)
}

func AdminDeletePlace(c echo.Context) error {
	companyID, err := getAdminCompanyID(c)
	if err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")
	var place models.Place
	if err := config.DB.First(&place, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Lieu introuvable"})
	}
	if place.CompanyID != companyID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à ce lieu"})
	}

	if err := config.DB.Delete(&place).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la suppression du lieu"})
	}
	return c.NoContent(http.StatusNoContent)
}
