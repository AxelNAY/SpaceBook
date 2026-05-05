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

func TestGetCompanies(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	company := helpers.CreateTestCompany(t, "Test Company")
	defer config.DB.Delete(&company)

	t.Run("get all companies", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/companies", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetCompanies(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var companies []models.Company
		json.Unmarshal(rec.Body.Bytes(), &companies)

		if len(companies) == 0 {
			t.Error("Expected at least one company")
		}
	})
}

func TestCreateCompany(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	t.Run("successful creation", func(t *testing.T) {
		payload := map[string]string{"name": "New Company"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/admin/companies", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateCompany(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var company models.Company
		json.Unmarshal(rec.Body.Bytes(), &company)

		if company.Name != "New Company" {
			t.Errorf("Expected name 'New Company', got '%s'", company.Name)
		}
		if company.ID.String() == "" {
			t.Error("Expected a non-empty ID")
		}

		defer config.DB.Delete(&company)
	})
}

func TestDeleteCompany(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	t.Run("delete existing company", func(t *testing.T) {
		company := helpers.CreateTestCompany(t, "Company To Delete")

		req := httptest.NewRequest(http.MethodDelete, "/admin/companies/"+company.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(company.ID.String())

		err := handlers.DeleteCompany(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}

		var count int64
		config.DB.Model(&models.Company{}).Where("id = ?", company.ID).Count(&count)
		if count != 0 {
			t.Error("Expected company to be deleted")
		}
	})
}
