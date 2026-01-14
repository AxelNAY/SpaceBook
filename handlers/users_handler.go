package handlers

import (
	"net/http"

	"spacebook/config"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)


/*
GET /admin/users
Admin only – list all users with User
*/
func GetUsers(c echo.Context) error {
	var users []models.User

	if err := config.DB.
		Order("created_at DESC").
		Find(&users).Error; err != nil {

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la récupération des utilisateurs",
		})
	}

	return c.JSON(http.StatusOK, users)
}

/*
GET /admin/user/:id
Admin only – Delete a user with User
*/
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	var count int64
	config.DB.Model(&models.Reservation{}).
		Where("user_id = ?", id).
		Count(&count)

	if count > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "L'utilisateur ne peut pas être supprimé car il a des réservations",
		})
	}

	if err := config.DB.Delete(&models.User{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Échec de la suppression de l'utilisateur",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
