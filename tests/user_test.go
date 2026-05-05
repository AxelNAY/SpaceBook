package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/helpers"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestGetUsers(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	user := helpers.CreateTestUser(t, "getuserstest@test.com", "getuserstest")
	defer config.DB.Delete(&user)

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
	helpers.SetupTestDB()

	e := echo.New()

	t.Run("delete user without reservations", func(t *testing.T) {
		user := helpers.CreateTestUser(t, "deletetest@test.com", "deletetest")

		req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+user.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())

		err := handlers.DeleteUser(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		var count int64
		config.DB.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
		if count != 0 {
			t.Error("Expected user to be deleted")
		}
	})

	t.Run("delete user with reservations - should fail", func(t *testing.T) {
		user := helpers.CreateTestUser(t, "deletewithres@test.com", "deletewithres")
		resource := helpers.CreateTestResource(t, "Delete Test Resource")

		reservation := models.Reservation{
			ID:          uuid.New(),
			UserID:      user.ID,
			RessourceID: resource.ID,
			Status:      "pending",
		}
		config.DB.Create(&reservation)

		defer func() {
			config.DB.Delete(&reservation)
			config.DB.Delete(&resource)
			config.DB.Delete(&user)
		}()

		req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+user.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())

		handlers.DeleteUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
