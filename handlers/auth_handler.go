package handlers

import (
	"net/http"

	"spacebook/config"
	"spacebook/middleware"
	"spacebook/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
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
			"error": "Invalid payload",
		})
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email, username and password are required",
		})
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "Email already registered",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to hash password",
		})
	}

	user := models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid payload",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email and password are required",
		})
	}

	// Find user
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}
