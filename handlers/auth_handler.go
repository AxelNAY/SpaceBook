package handlers

import (
	"net/http"

	"spacebook/config"
	"spacebook/middleware"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func resolveAdminCompanyID(userID interface{}) string {
	var placeUser models.PlaceUser
	if err := config.DB.Preload("Place").Where("user_id = ?", userID).First(&placeUser).Error; err != nil {
		return ""
	}
	return placeUser.Place.CompanyID.String()
}

type RegisterRequest struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Phone    int     `json:"phone"`
	PlaceID  *string `json:"place_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Données invalides",
		})
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email, nom d'utilisateur et mot de passe requis",
		})
	}

	if req.PlaceID == nil || *req.PlaceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Un lieu est requis pour l'inscription",
		})
	}

	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Cet email est déjà utilisé",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Échec du hachage du mot de passe",
		})
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Role:     "user",
		Status:   "En attente",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Échec de la création de l'utilisateur",
		})
	}

	placeID, err := uuid.Parse(*req.PlaceID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Identifiant de lieu invalide",
		})
	}

	placeUser := models.PlaceUser{
		ID:      uuid.New(),
		PlaceID: placeID,
		UserID:  user.ID,
	}
	config.DB.Create(&placeUser)

	// Notifier les admins du lieu
	var place models.Place
	if err := config.DB.First(&place, "id = ?", placeID).Error; err == nil {
		var adminUsers []models.User
		config.DB.
			Joins("JOIN place_users pu ON pu.user_id = users.id").
			Joins("JOIN places p ON p.id = pu.place_id").
			Where("p.company_id = ? AND users.role = ?", place.CompanyID, "admin").
			Find(&adminUsers)

		for _, admin := range adminUsers {
			adminID := admin.ID
			notification := models.Notification{
				UserID:  &adminID,
				Type:    "user_registration",
				Message: user.Username + " (" + user.Email + ") a demandé à rejoindre votre espace",
				IsRead:  false,
			}
			config.DB.Create(&notification)
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Votre compte a été créé et est en attente de validation par un administrateur",
		"user":    user,
	})
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Données invalides",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email et mot de passe requis",
		})
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Identifiants invalides",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Identifiants invalides",
		})
	}

	if user.Role == "user" && user.Status == "En attente" {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "Votre compte est en attente de validation par un administrateur",
		})
	}

	companyID := ""
	if user.Role == "admin" {
		companyID = resolveAdminCompanyID(user.ID)
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, companyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Échec de la génération du token",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}
