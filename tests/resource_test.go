package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/helpers"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func TestGetResources(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	resource := helpers.CreateTestResource(t, "Test Room")
	defer config.DB.Delete(&resource)

	t.Run("get all resources", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resources", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetResources(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var resources []models.Resource
		json.Unmarshal(rec.Body.Bytes(), &resources)

		if len(resources) == 0 {
			t.Error("Expected at least one resource")
		}
	})
}

func TestCreateResource(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	t.Run("create resource", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":   "New Test Resource",
			"status": "available",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/admin/resources", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateResource(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var resource models.Resource
		json.Unmarshal(rec.Body.Bytes(), &resource)

		if resource.Name != "New Test Resource" {
			t.Errorf("Expected name 'New Test Resource', got '%s'", resource.Name)
		}

		defer config.DB.Delete(&resource)
	})
}
