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

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestGetCategories(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		category := helpers.CreateTestCategory(t, "GetCategories Category "+id)
		t.Cleanup(func() { config.DB.Delete(&category) })

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetCategories(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var categories []models.Category
		json.Unmarshal(rec.Body.Bytes(), &categories)
		if len(categories) == 0 {
			t.Error("Expected at least one category")
		}
	})
}

func TestCreateCategory(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		payload := map[string]string{"name": "New Category " + id}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/admin/categories", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateCategory(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var category models.Category
		json.Unmarshal(rec.Body.Bytes(), &category)
		t.Cleanup(func() { config.DB.Delete(&category) })

		if category.Name != "New Category "+id {
			t.Errorf("Expected name 'New Category %s', got '%s'", id, category.Name)
		}
	})
}

func TestUpdateCategory(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		category := helpers.CreateTestCategory(t, "Old Category Name "+id)
		t.Cleanup(func() { config.DB.Delete(&category) })

		payload := map[string]string{"name": "New Category Name " + id}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/admin/categories/"+category.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(category.ID.String())

		err := handlers.UpdateCategory(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.Category
		config.DB.First(&updated, "id = ?", category.ID)
		if updated.Name != "New Category Name "+id {
			t.Errorf("Expected name 'New Category Name %s', got '%s'", id, updated.Name)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		payload := map[string]string{"name": "Test"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/admin/categories/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		handlers.UpdateCategory(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestDeleteCategory(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		category := helpers.CreateTestCategory(t, "Category To Delete "+id)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&category) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/categories/"+category.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(category.ID.String())

		err := handlers.DeleteCategory(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		var count int64
		config.DB.Model(&models.Category{}).Where("id = ?", category.ID).Count(&count)
		if count != 0 {
			t.Error("Expected category to be deleted")
		}
	})
}
