package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/models"

	"github.com/labstack/echo/v4"
)

func TestGetResources(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Create test resource
	resource := models.Resource{
		Name:     "Test Room",
		Type:     "room",
		Category: "none",
		Capacity: 10,
		Status:   "available",
	}
	config.DB.Create(&resource)

	defer config.DB.Where("name = ?", "Test Room").Delete(&models.Resource{})

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
	setupTestDB()

	e := echo.New()

	defer config.DB.Where("name = ?", "New Test Resource").Delete(&models.Resource{})

	t.Run("create resource", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":     "New Test Resource",
			"type":     "equipment",
			"category": "printer",
			"capacity": 3,
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

		if resource.Capacity != 3 {
			t.Errorf("Expected capacity 3, got %d", resource.Capacity)
		}
	})

	t.Run("create resource with default capacity", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "Default Capacity Resource",
			"type": "room",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/admin/resources", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateResource(c)

		var resource models.Resource
		json.Unmarshal(rec.Body.Bytes(), &resource)

		if resource.Capacity != 1 {
			t.Errorf("Expected default capacity 1, got %d", resource.Capacity)
		}

		// Cleanup
		config.DB.Where("name = ?", "Default Capacity Resource").Delete(&models.Resource{})
	})
}
