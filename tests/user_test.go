package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestGetUsers(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Create test user
	user := models.User{
		ID:       uuid.New(),
		Email:    "getuserstest@test.com",
		Username: "getuserstest",
		Password: []byte("password"),
		Role:     "user",
	}
	config.DB.Create(&user)

	defer config.DB.Where("email = ?", "getuserstest@test.com").Delete(&models.User{})

	t.Run("get all users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetUsers(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var users []models.User
		json.Unmarshal(rec.Body.Bytes(), &users)

		if len(users) == 0 {
			t.Error("Expected at least one user")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	setupTestDB()

	e := echo.New()

	t.Run("delete user without reservations", func(t *testing.T) {
		// Create test user
		userID := uuid.New()
		user := models.User{
			ID:       userID,
			Email:    "deletetest@test.com",
			Username: "deletetest",
			Password: []byte("password"),
			Role:     "user",
		}
		config.DB.Create(&user)

		req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+userID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(userID.String())

		err := handlers.DeleteUser(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		// Verify user is deleted
		var count int64
		config.DB.Model(&models.User{}).Where("id = ?", userID).Count(&count)
		if count != 0 {
			t.Error("Expected user to be deleted")
		}
	})

	t.Run("delete user with reservations - should fail", func(t *testing.T) {
		// Create test user
		userID := uuid.New()
		user := models.User{
			ID:       userID,
			Email:    "deletewithres@test.com",
			Username: "deletewithres",
			Password: []byte("password"),
			Role:     "user",
		}
		config.DB.Create(&user)

		// Create test resource
		resource := models.Resource{
			ID:       uuid.New().String(),
			Name:     "Delete Test Resource",
			Type:     "room",
			Capacity: 1,
			Status:   "available",
		}
		config.DB.Create(&resource)

		// Create reservation for user
		reservation := models.Reservation{
			ID:         uuid.New(),
			UserID:     userID,
			ResourceID: uuid.MustParse(resource.ID),
			Status:     "pending",
		}
		config.DB.Create(&reservation)

		defer func() {
			config.DB.Where("id = ?", reservation.ID).Delete(&models.Reservation{})
			config.DB.Where("id = ?", resource.ID).Delete(&models.Resource{})
			config.DB.Where("id = ?", userID).Delete(&models.User{})
		}()

		req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+userID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(userID.String())

		handlers.DeleteUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
