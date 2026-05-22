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

func TestGetPlaces(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("get_all", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "GetPlaces Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "GetPlaces Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

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

	t.Run("filter_by_company_id", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "FilterPlaces Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "FilterPlaces Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		req := httptest.NewRequest(http.MethodGet, "/places?company_id="+company.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetPlaces(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
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

func TestAdminGetPlaces(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminGetPlaces Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminGetPlaces Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })
		user := helpers.CreateTestUser(t, "agp_"+id+"@test.com", "agp_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AGP Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		req := httptest.NewRequest(http.MethodGet, "/admin/places", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())

		err := handlers.AdminGetPlaces(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var result []map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &result)
		if len(result) == 0 {
			t.Error("Expected at least one place")
		}

		// Verify places contain users and resources fields
		found := false
		for _, p := range result {
			if p["id"] == place.ID.String() {
				found = true
				if _, ok := p["users"]; !ok {
					t.Error("Expected 'users' field in place details")
				}
				if _, ok := p["resources"]; !ok {
					t.Error("Expected 'resources' field in place details")
				}
			}
		}
		if !found {
			t.Error("Expected to find the created place in the response")
		}
	})

	t.Run("no_company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/places", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.AdminGetPlaces(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminGetPlaceDetails(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminGetPD Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminGetPD Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })
		user := helpers.CreateTestUser(t, "agpd_"+id+"@test.com", "agpd_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AGPD Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		req := httptest.NewRequest(http.MethodGet, "/admin/places/"+place.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AdminGetPlaceDetails(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminGetPD NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodGet, "/admin/places/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.AdminGetPlaceDetails(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AdminGetPD WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AdminGetPD WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "AdminGetPD WC Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		req := httptest.NewRequest(http.MethodGet, "/admin/places/"+place.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.AdminGetPlaceDetails(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminCreatePlace(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminCreatePlace Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		adminUser := helpers.CreateTestAdmin(t, "acp_admin_"+id+"@test.com", "acp_admin_"+id)
		t.Cleanup(func() { config.DB.Delete(&adminUser) })

		payload := map[string]interface{}{
			"name":        "New Admin Place " + id,
			"description": "desc",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/admin/places", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())
		c.Set("user_id", adminUser.ID)

		err := handlers.AdminCreatePlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var place models.Place
		json.Unmarshal(rec.Body.Bytes(), &place)
		t.Cleanup(func() {
			config.DB.Unscoped().Where("user_id = ?", adminUser.ID).Delete(&models.PlaceUser{})
			config.DB.Delete(&place)
		})

		// Verify PlaceUser was created for admin
		var count int64
		config.DB.Model(&models.PlaceUser{}).Where("place_id = ? AND user_id = ?", place.ID, adminUser.ID).Count(&count)
		if count != 1 {
			t.Error("Expected PlaceUser to be created for admin")
		}
	})

	t.Run("no_company_id", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "No Company Place",
			"description": "desc",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/admin/places", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.AdminCreatePlace(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminUpdatePlace(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminUpdatePlace Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "Old Name "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		payload := map[string]interface{}{
			"name":        "New Name " + id,
			"description": "desc",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/places/"+place.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AdminUpdatePlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.Place
		config.DB.First(&updated, "id = ?", place.ID)
		if updated.Name != "New Name "+id {
			t.Errorf("Expected name 'New Name %s', got '%s'", id, updated.Name)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminUpdatePlace NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		payload := map[string]interface{}{"name": "test"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/places/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.AdminUpdatePlace(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AdminUpdatePlace WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AdminUpdatePlace WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "WC Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		payload := map[string]interface{}{"name": "Renamed"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/places/"+place.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.AdminUpdatePlace(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminDeletePlace(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminDeletePlace Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminDeletePlace Place "+id, "desc")
		t.Cleanup(func() { config.DB.Unscoped().Delete(&place) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/places/"+place.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AdminDeletePlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminDeletePlace NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/places/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.AdminDeletePlace(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AdminDeletePlace WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AdminDeletePlace WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "AdminDeletePlace WC Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/places/"+place.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(place.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.AdminDeletePlace(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}
