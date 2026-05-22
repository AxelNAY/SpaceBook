package handlers

import (
	"net/http"

	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

/*
GET /admin/users
Admin only – liste les utilisateurs classiques de l'entreprise de l'admin
*/
func GetUsers(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Administrateur non associé à une entreprise",
		})
	}

	var users []models.User
	if err := config.DB.
		Joins("JOIN place_users pu ON pu.user_id = users.id").
		Joins("JOIN places p ON p.id = pu.place_id").
		Where("p.company_id = ?", companyIDStr).
		Where("users.role = ?", "user").
		Order("users.created_at DESC").
		Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la récupération des utilisateurs",
		})
	}

	return c.JSON(http.StatusOK, users)
}

func verifyUserInCompany(userID string, companyIDStr string) bool {
	var count int64
	config.DB.Model(&models.PlaceUser{}).
		Joins("JOIN places p ON p.id = place_users.place_id").
		Where("place_users.user_id = ? AND p.company_id = ?", userID, companyIDStr).
		Count(&count)
	return count > 0
}

/*
PUT /admin/users/accept/:id
Admin only – valide un utilisateur en attente
*/
func AcceptUser(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Utilisateur introuvable"})
	}
	if user.Role != "user" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Seuls les utilisateurs classiques peuvent être validés"})
	}
	if user.Status != "En attente" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "L'utilisateur n'est pas en attente de validation"})
	}
	if !verifyUserInCompany(id, companyIDStr) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Cet utilisateur n'appartient pas à votre entreprise"})
	}

	user.Status = "Actif"
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la mise à jour de l'utilisateur"})
	}

	userID := user.ID
	notification := models.Notification{
		UserID:  &userID,
		Type:    "account",
		Message: "Votre compte a été validé, vous pouvez maintenant vous connecter",
		IsRead:  false,
	}
	config.DB.Create(&notification)

	return c.JSON(http.StatusOK, user)
}

/*
DELETE /admin/users/refuse/:id
Admin only – supprime un utilisateur en attente
*/
func RefuseUser(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Utilisateur introuvable"})
	}
	if user.Role != "user" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Seuls les utilisateurs classiques peuvent être refusés"})
	}
	if user.Status != "En attente" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "L'utilisateur n'est pas en attente de validation"})
	}
	if !verifyUserInCompany(id, companyIDStr) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Cet utilisateur n'appartient pas à votre entreprise"})
	}

	if err := config.DB.Delete(&models.Reservation{}, "user_id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la suppression des réservations de l'utilisateur"})
	}
	if err := config.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la suppression de l'utilisateur"})
	}

	return c.NoContent(http.StatusNoContent)
}

/*
PUT /admin/users/promote/:id
Admin only – promeut un utilisateur au rang d'admin
*/
func PromoteUser(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Utilisateur introuvable"})
	}
	if user.Role != "user" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Seuls les utilisateurs classiques peuvent être promus"})
	}
	if user.Status != "Actif" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Seuls les utilisateurs actifs peuvent être promus"})
	}
	if !verifyUserInCompany(id, companyIDStr) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Cet utilisateur n'appartient pas à votre entreprise"})
	}

	user.Role = "admin"
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de la promotion de l'utilisateur"})
	}

	userID := user.ID
	notification := models.Notification{
		UserID:  &userID,
		Type:    "account",
		Message: "Vous avez été promu administrateur",
		IsRead:  false,
	}
	config.DB.Create(&notification)

	return c.JSON(http.StatusOK, user)
}

/*
PUT /admin/users/:id/place
Admin only – transfère ou ajoute un utilisateur à un lieu de la même entreprise
Body: { place_id: string, transfer: bool }
  - transfer=true  : retire l'utilisateur de ses lieux actuels dans l'entreprise puis l'affecte au nouveau lieu
  - transfer=false : ajoute le lieu en plus des lieux existants
*/
func AssignUserPlace(c echo.Context) error {
	companyIDStr, _ := c.Get("company_id").(string)
	if companyIDStr == "" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Administrateur non associé à une entreprise"})
	}

	id := c.Param("id")

	var req struct {
		PlaceID  string `json:"place_id"`
		Transfer bool   `json:"transfer"`
	}
	if err := c.Bind(&req); err != nil || req.PlaceID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "place_id requis"})
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Utilisateur introuvable"})
	}
	if user.Role != "user" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Seuls les utilisateurs classiques peuvent être réaffectés"})
	}
	if !verifyUserInCompany(id, companyIDStr) {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Cet utilisateur n'appartient pas à votre entreprise"})
	}

	// Vérifier que le lieu cible appartient à la même entreprise
	var targetPlace models.Place
	if err := config.DB.First(&targetPlace, "id = ?", req.PlaceID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Lieu introuvable"})
	}
	if targetPlace.CompanyID.String() != companyIDStr {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Ce lieu n'appartient pas à votre entreprise"})
	}

	placeID, err := uuid.Parse(req.PlaceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Identifiant de lieu invalide"})
	}

	if req.Transfer {
		// Supprimer les affectations existantes dans cette entreprise
		config.DB.
			Where("user_id = ? AND place_id IN (SELECT id FROM places WHERE company_id = ?)", id, companyIDStr).
			Delete(&models.PlaceUser{})
	} else {
		// Vérifier que l'utilisateur n'est pas déjà dans ce lieu
		var count int64
		config.DB.Model(&models.PlaceUser{}).Where("user_id = ? AND place_id = ?", id, placeID).Count(&count)
		if count > 0 {
			return c.JSON(http.StatusConflict, echo.Map{"error": "L'utilisateur est déjà affecté à ce lieu"})
		}
	}

	placeUser := models.PlaceUser{
		ID:      uuid.New(),
		PlaceID: placeID,
		UserID:  user.ID,
	}
	if err := config.DB.Create(&placeUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Échec de l'affectation au lieu"})
	}

	return c.JSON(http.StatusOK, user)
}

/*
DELETE /admin/user/:id
Admin only – supprime un utilisateur classique
*/
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Utilisateur introuvable",
		})
	}
	if user.Role != "user" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "Seuls les utilisateurs classiques peuvent être supprimés",
		})
	}

	if err := config.DB.Delete(&models.Reservation{}, "user_id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression des réservations de l'utilisateur",
		})
	}

	if err := config.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression de l'utilisateur",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
