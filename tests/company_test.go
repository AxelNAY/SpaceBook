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

func TestGetCompanies(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("get_all_companies", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "GetCompanies Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

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

	t.Run("successful_creation", func(t *testing.T) {
		id := uuid.New().String()[:8]
		payload := map[string]string{"name": "New Company " + id}
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
		t.Cleanup(func() { config.DB.Delete(&company) })

		if company.Name != "New Company "+id {
			t.Errorf("Expected name 'New Company %s', got '%s'", id, company.Name)
		}
		if company.ID == (uuid.UUID{}) {
			t.Error("Expected a non-empty ID")
		}
	})
}

func TestDeleteCompany(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("delete_existing_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "Company To Delete "+id)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&company) })

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

func TestUpdateCompany(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "Old Company Name "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		payload := map[string]string{"name": "New Company Name " + id}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/admin/companies/"+company.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(company.ID.String())

		err := handlers.UpdateCompany(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.Company
		config.DB.First(&updated, "id = ?", company.ID)
		if updated.Name != "New Company Name "+id {
			t.Errorf("Expected name 'New Company Name %s', got '%s'", id, updated.Name)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		payload := map[string]string{"name": "Test"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/admin/companies/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		handlers.UpdateCompany(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
