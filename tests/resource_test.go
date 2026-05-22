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

func TestGetResources(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("returns_only_user_place_resources", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "GetRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "GetRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		user := helpers.CreateTestUser(t, "getres_"+id+"@test.com", "getres_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		pu := helpers.CreateTestPlaceUser(t, place.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&pu) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "GetRes Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		req := httptest.NewRequest(http.MethodGet, "/resources", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

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
		for _, r := range resources {
			if r.PlaceID == nil || *r.PlaceID != place.ID {
				t.Errorf("Expected only resources from place %s", place.ID)
			}
		}
	})

	t.Run("no_auth_returns_401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resources", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.GetResources(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("user_not_in_any_place_returns_empty", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "getres_npl_"+id+"@test.com", "getres_npl_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		req := httptest.NewRequest(http.MethodGet, "/resources", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID)

		handlers.GetResources(c)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}
		var resources []models.Resource
		json.Unmarshal(rec.Body.Bytes(), &resources)
		if len(resources) != 0 {
			t.Error("Expected empty list for user with no place")
		}
	})
}

func TestAdminGetResources(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminGetRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminGetRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AdminGetRes Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		req := httptest.NewRequest(http.MethodGet, "/admin/resources", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())

		err := handlers.AdminGetResources(c)
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

	t.Run("no_company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/resources", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.AdminGetResources(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminCreateResource(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminCreateRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminCreateRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })

		payload := map[string]interface{}{
			"name":     "New Admin Resource " + id,
			"place_id": place.ID.String(),
			"status":   "available",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/admin/resources", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())

		err := handlers.AdminCreateResource(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var resource models.Resource
		json.Unmarshal(rec.Body.Bytes(), &resource)
		t.Cleanup(func() { config.DB.Delete(&resource) })
	})

	t.Run("no_company_id", func(t *testing.T) {
		id := uuid.New().String()[:8]
		payload := map[string]interface{}{
			"name":   "Resource " + id,
			"status": "available",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/admin/resources", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.AdminCreateResource(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("place_not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminCreateRes PNF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		fakeID := uuid.New()
		payload := map[string]interface{}{
			"name":     "Resource PNF " + id,
			"place_id": fakeID.String(),
			"status":   "available",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/admin/resources", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())

		handlers.AdminCreateResource(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("place_wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AdminCreateRes WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AdminCreateRes WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "AdminCreateRes WC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })

		payload := map[string]interface{}{
			"name":     "Resource WC " + id,
			"place_id": place.ID.String(),
			"status":   "available",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/admin/resources", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company2.ID.String())

		handlers.AdminCreateResource(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminUpdateResource(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminUpdateRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminUpdateRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "Old Resource Name "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		payload := map[string]interface{}{
			"name":   "Updated Resource Name " + id,
			"status": "available",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/resources/"+resource.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AdminUpdateResource(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.Resource
		config.DB.First(&updated, "id = ?", resource.ID)
		if updated.Name != "Updated Resource Name "+id {
			t.Errorf("Expected name 'Updated Resource Name %s', got '%s'", id, updated.Name)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminUpdateRes NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		payload := map[string]interface{}{"name": "test"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/resources/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.AdminUpdateResource(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AdminUpdateRes WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AdminUpdateRes WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "AdminUpdateRes WC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AdminUpdateRes WC Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		payload := map[string]interface{}{"name": "Renamed"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/resources/"+resource.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.AdminUpdateResource(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAdminDeleteResource(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminDeleteRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminDeleteRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AdminDeleteRes Resource "+id)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&resource) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/resources/"+resource.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AdminDeleteResource(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminDeleteRes NF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/resources/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.AdminDeleteResource(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("has_reservations", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AdminDeleteRes HasRes Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AdminDeleteRes HasRes Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AdminDeleteRes HasRes Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "adres_hasres_"+id+"@test.com", "adres_hasres_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/resources/"+resource.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.AdminDeleteResource(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AdminDeleteRes WC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AdminDeleteRes WC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "AdminDeleteRes WC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		resource := helpers.CreateTestResourceInPlace(t, place.ID, "AdminDeleteRes WC Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/resources/"+resource.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.AdminDeleteResource(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestDeleteResource(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		resource := helpers.CreateTestResource(t, "DeleteRes Resource "+id)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&resource) })

		req := httptest.NewRequest(http.MethodDelete, "/resources/"+resource.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())

		err := handlers.DeleteResource(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}
	})

	t.Run("has_reservations", func(t *testing.T) {
		id := uuid.New().String()[:8]
		resource := helpers.CreateTestResource(t, "DeleteRes HasRes Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		user := helpers.CreateTestUser(t, "delres_hasres_"+id+"@test.com", "delres_hasres_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Delete(&reservation) })

		req := httptest.NewRequest(http.MethodDelete, "/resources/"+resource.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())

		handlers.DeleteResource(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUpdateResource(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		resource := helpers.CreateTestResource(t, "UpdateRes Old Name "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })

		payload := map[string]interface{}{
			"name":   "UpdateRes New Name " + id,
			"status": "available",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/resources/"+resource.ID.String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(resource.ID.String())

		err := handlers.UpdateResource(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("not_found", func(t *testing.T) {
		payload := map[string]interface{}{"name": "test"}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/resources/"+uuid.New().String(), bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		handlers.UpdateResource(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
