package handlers

import (
	"net/http"
	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetResources(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Non authentifié"})
	}

	var placeUsers []models.PlaceUser
	config.DB.Where("user_id = ?", userID).Find(&placeUsers)

	placeIDs := make([]uuid.UUID, 0, len(placeUsers))
	for _, pu := range placeUsers {
		placeIDs = append(placeIDs, pu.PlaceID)
	}

	if len(placeIDs) == 0 {
		return c.JSON(http.StatusOK, []models.Resource{})
	}

	var resources []models.Resource
	config.DB.Preload("Place").Preload("Category").
		Where("place_id IN ?", placeIDs).
		Find(&resources)

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

func getResourceCompanyID(resource models.Resource) (uuid.UUID, error) {
	if resource.Place == nil {
		var place models.Place
		if err := config.DB.First(&place, "id = ?", resource.PlaceID).Error; err != nil {
			return uuid.UUID{}, err
		}
		return place.CompanyID, nil
	}
	return resource.Place.CompanyID, nil
}

func AdminGetResources(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	var resources []models.Resource
	if err := config.DB.Preload("Place").Preload("Category").
		Joins("JOIN places p ON p.id = resources.place_id").
		Where("p.company_id = ?", companyIDStr).
		Find(&resources).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la récupération des ressources"})
	}
	return c.JSON(http.StatusOK, resources)
}

func AdminCreateResource(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}
	companyID, _ := uuid.Parse(companyIDStr)

	var resource models.Resource
	if err := c.Bind(&resource); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Données invalides"})
	}

	var place models.Place
	if err := config.DB.First(&place, "id = ?", resource.PlaceID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Lieu introuvable"})
	}
	if place.CompanyID != companyID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à ce lieu"})
	}

	if err := config.DB.Create(&resource).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la création de la ressource"})
	}

	config.DB.Create(&models.Notification{Type: "resource", Message: "Une nouvelle ressource a été créée"})
	return c.JSON(http.StatusCreated, resource)
}

func AdminUpdateResource(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}
	companyID, _ := uuid.Parse(companyIDStr)

	id := c.Param("id")
	var resource models.Resource
	if err := config.DB.Preload("Place").First(&resource, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Ressource introuvable"})
	}

	resourceCompanyID, err := getResourceCompanyID(resource)
	if err != nil || resourceCompanyID != companyID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à cette ressource"})
	}

	if err := c.Bind(&resource); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Données invalides"})
	}
	if err := config.DB.Save(&resource).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la mise à jour de la ressource"})
	}
	return c.JSON(http.StatusOK, resource)
}

func AdminDeleteResource(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}
	companyID, _ := uuid.Parse(companyIDStr)

	id := c.Param("id")
	var resource models.Resource
	if err := config.DB.Preload("Place").First(&resource, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Ressource introuvable"})
	}

	resourceCompanyID, err := getResourceCompanyID(resource)
	if err != nil || resourceCompanyID != companyID {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Accès refusé à cette ressource"})
	}

	var count int64
	config.DB.Model(&models.Reservation{}).Where("ressource_id = ?", id).Count(&count)
	if count > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "La ressource ne peut pas être supprimée car elle a des réservations"})
	}

	if err := config.DB.Delete(&models.Resource{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la suppression de la ressource"})
	}
	return c.NoContent(http.StatusNoContent)
}
