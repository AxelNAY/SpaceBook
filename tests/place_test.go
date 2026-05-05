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

func TestGetPlaces(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	company := helpers.CreateTestCompany(t, "Places Test Company")
	place := helpers.CreateTestPlace(t, company.ID, "Test Place", "A test place")
	defer func() {
		config.DB.Delete(&place)
		config.DB.Delete(&company)
	}()

	t.Run("get all places", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/places", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetPlaces(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var places []models.Place
		json.Unmarshal(rec.Body.Bytes(), &places)

		if len(places) == 0 {
			t.Error("Expected at least one place")
		}
	})

	t.Run("filter by company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/places?company_id="+company.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.QueryParams().Set("company_id", company.ID.String())

		err := handlers.GetPlaces(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var places []models.Place
		json.Unmarshal(rec.Body.Bytes(), &places)

		for _, p := range places {
			if p.CompanyID != company.ID {
				t.Errorf("Expected all places to belong to company %s, got %s", company.ID, p.CompanyID)
			}
		}
	})
}

func TestCreatePlace(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	company := helpers.CreateTestCompany(t, "Create Place Company")
	defer config.DB.Delete(&company)

	t.Run("successful creation", func(t *testing.T) {
		payload := map[string]interface{}{
			"company_id":  company.ID.String(),
			"name":        "New Place",
			"description": "A brand new place",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/admin/places", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreatePlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var place models.Place
		json.Unmarshal(rec.Body.Bytes(), &place)

		if place.Name != "New Place" {
			t.Errorf("Expected name 'New Place', got '%s'", place.Name)
		}
		if place.CompanyID != company.ID {
			t.Errorf("Expected company_id '%s', got '%s'", company.ID, place.CompanyID)
		}

		defer config.DB.Delete(&place)
	})
}

func TestDeletePlace(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	t.Run("delete existing place", func(t *testing.T) {
		company := helpers.CreateTestCompany(t, "Delete Place Company")
		place := helpers.CreateTestPlace(t, company.ID, "Place To Delete", "")
		defer config.DB.Delete(&company)

		req := httptest.NewRequest(http.MethodDelete, "/admin/places/"+place.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())

		err := handlers.DeletePlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		var count int64
		config.DB.Model(&models.Place{}).Where("id = ?", place.ID).Count(&count)
		if count != 0 {
			t.Error("Expected place to be deleted")
		}
	})
}
