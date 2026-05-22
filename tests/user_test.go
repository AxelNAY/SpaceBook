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

func TestGetUsers(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "GetUsers Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "GetUsers Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		user := helpers.CreateTestUser(t, "getusers_"+id+"@test.com", "getusers_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("company_id", company.ID.String())

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

	t.Run("no_company_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.GetUsers(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAcceptUser(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AcceptUser Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AcceptUser Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		pendingUser := helpers.CreateTestPendingUser(t, "accept_"+id+"@test.com", "accept_"+id)
		t.Cleanup(func() { config.DB.Delete(&pendingUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, pendingUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/accept/"+pendingUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(pendingUser.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AcceptUser(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.User
		config.DB.First(&updated, "id = ?", pendingUser.ID)
		if updated.Status != "Actif" {
			t.Errorf("Expected status 'Actif', got '%s'", updated.Status)
		}
	})

	t.Run("user_not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AcceptNotFound Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/accept/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.AcceptUser(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("already_active", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AcceptActive Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AcceptActive Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		activeUser := helpers.CreateTestUser(t, "acceptact_"+id+"@test.com", "acceptact_"+id)
		t.Cleanup(func() { config.DB.Delete(&activeUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, activeUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/accept/"+activeUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(activeUser.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.AcceptUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AcceptWC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AcceptWC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "AcceptWC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		pendingUser := helpers.CreateTestPendingUser(t, "acceptwc_"+id+"@test.com", "acceptwc_"+id)
		t.Cleanup(func() { config.DB.Delete(&pendingUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, pendingUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/accept/"+pendingUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(pendingUser.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.AcceptUser(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestRefuseUser(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RefuseUser Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "RefuseUser Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		pendingUser := helpers.CreateTestPendingUser(t, "refuse_"+id+"@test.com", "refuse_"+id)
		t.Cleanup(func() {
			config.DB.Unscoped().Where("user_id = ?", pendingUser.ID).Delete(&models.PlaceUser{})
			config.DB.Unscoped().Delete(&pendingUser)
		})
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, pendingUser.ID)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/users/refuse/"+pendingUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(pendingUser.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.RefuseUser(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusNoContent, rec.Code, rec.Body.String())
		}

		var count int64
		config.DB.Unscoped().Model(&models.User{}).Where("id = ?", pendingUser.ID).Count(&count)
		if count != 0 {
			t.Error("Expected user to be deleted")
		}
	})

	t.Run("user_not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RefuseNF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/users/refuse/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.RefuseUser(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("not_pending", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "RefuseNP Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "RefuseNP Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		activeUser := helpers.CreateTestUser(t, "refusenp_"+id+"@test.com", "refusenp_"+id)
		t.Cleanup(func() { config.DB.Delete(&activeUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, activeUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/users/refuse/"+activeUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(activeUser.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.RefuseUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestPromoteUser(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "PromoteUser Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "PromoteUser Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		activeUser := helpers.CreateTestUser(t, "promote_"+id+"@test.com", "promote_"+id)
		t.Cleanup(func() { config.DB.Delete(&activeUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, activeUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/promote/"+activeUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(activeUser.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.PromoteUser(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var updated models.User
		config.DB.First(&updated, "id = ?", activeUser.ID)
		if updated.Role != "admin" {
			t.Errorf("Expected role 'admin', got '%s'", updated.Role)
		}
	})

	t.Run("user_not_found", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "PromoteNF Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/promote/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())
		c.Set("company_id", company.ID.String())

		handlers.PromoteUser(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("already_admin", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "PromoteAA Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "PromoteAA Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		admin := helpers.CreateTestAdmin(t, "promoteaa_"+id+"@test.com", "promoteaa_"+id)
		t.Cleanup(func() { config.DB.Delete(&admin) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, admin.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/promote/"+admin.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(admin.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.PromoteUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("not_active", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "PromoteNA Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "PromoteNA Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		pendingUser := helpers.CreateTestPendingUser(t, "promotena_"+id+"@test.com", "promotena_"+id)
		t.Cleanup(func() { config.DB.Delete(&pendingUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, pendingUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/promote/"+pendingUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(pendingUser.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.PromoteUser(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "PromoteWC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "PromoteWC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place := helpers.CreateTestPlace(t, company1.ID, "PromoteWC Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		activeUser := helpers.CreateTestUser(t, "promotewc_"+id+"@test.com", "promotewc_"+id)
		t.Cleanup(func() { config.DB.Delete(&activeUser) })
		placeUser := helpers.CreateTestPlaceUser(t, place.ID, activeUser.ID)
		t.Cleanup(func() { config.DB.Delete(&placeUser) })

		req := httptest.NewRequest(http.MethodPut, "/admin/users/promote/"+activeUser.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(activeUser.ID.String())
		c.Set("company_id", company2.ID.String())

		handlers.PromoteUser(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})
}

func TestAssignUserPlace(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("transfer_success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AssignTransfer Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place1 := helpers.CreateTestPlace(t, company.ID, "AssignTransfer Place1 "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place1) })
		place2 := helpers.CreateTestPlace(t, company.ID, "AssignTransfer Place2 "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place2) })
		user := helpers.CreateTestUser(t, "assign_tr_"+id+"@test.com", "assign_tr_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		pu1 := helpers.CreateTestPlaceUser(t, place1.ID, user.ID)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&pu1) })

		payload := map[string]interface{}{
			"place_id": place2.ID.String(),
			"transfer": true,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/users/"+user.ID.String()+"/place", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AssignUserPlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var pus []models.PlaceUser
		config.DB.Where("user_id = ?", user.ID).Find(&pus)
		if len(pus) != 1 {
			t.Errorf("Expected exactly 1 place assignment, got %d", len(pus))
		}
		if len(pus) == 1 && pus[0].PlaceID != place2.ID {
			t.Errorf("Expected assignment to place2, got %s", pus[0].PlaceID)
		}
		t.Cleanup(func() {
			config.DB.Unscoped().Where("user_id = ?", user.ID).Delete(&models.PlaceUser{})
		})
	})

	t.Run("add_success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AssignAdd Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place1 := helpers.CreateTestPlace(t, company.ID, "AssignAdd Place1 "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place1) })
		place2 := helpers.CreateTestPlace(t, company.ID, "AssignAdd Place2 "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place2) })
		user := helpers.CreateTestUser(t, "assign_add_"+id+"@test.com", "assign_add_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		pu1 := helpers.CreateTestPlaceUser(t, place1.ID, user.ID)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&pu1) })

		payload := map[string]interface{}{
			"place_id": place2.ID.String(),
			"transfer": false,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/users/"+user.ID.String()+"/place", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())
		c.Set("company_id", company.ID.String())

		err := handlers.AssignUserPlace(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var pus []models.PlaceUser
		config.DB.Where("user_id = ?", user.ID).Find(&pus)
		if len(pus) != 2 {
			t.Errorf("Expected 2 place assignments, got %d", len(pus))
		}
		t.Cleanup(func() {
			config.DB.Unscoped().Where("user_id = ?", user.ID).Delete(&models.PlaceUser{})
		})
	})

	t.Run("already_in_place", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AssignDup Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AssignDup Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		user := helpers.CreateTestUser(t, "assign_dup_"+id+"@test.com", "assign_dup_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		pu := helpers.CreateTestPlaceUser(t, place.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&pu) })

		payload := map[string]interface{}{
			"place_id": place.ID.String(),
			"transfer": false,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/users/"+user.ID.String()+"/place", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.AssignUserPlace(c)

		if rec.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, rec.Code)
		}
	})

	t.Run("target_place_wrong_company", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company1 := helpers.CreateTestCompany(t, "AssignWC Company1 "+id)
		t.Cleanup(func() { config.DB.Delete(&company1) })
		company2 := helpers.CreateTestCompany(t, "AssignWC Company2 "+id)
		t.Cleanup(func() { config.DB.Delete(&company2) })
		place1 := helpers.CreateTestPlace(t, company1.ID, "AssignWC Place1 "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place1) })
		place2 := helpers.CreateTestPlace(t, company2.ID, "AssignWC Place2 "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place2) })
		user := helpers.CreateTestUser(t, "assign_wc_"+id+"@test.com", "assign_wc_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		pu := helpers.CreateTestPlaceUser(t, place1.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&pu) })

		payload := map[string]interface{}{
			"place_id": place2.ID.String(),
			"transfer": false,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/users/"+user.ID.String()+"/place", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())
		c.Set("company_id", company1.ID.String())

		handlers.AssignUserPlace(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("missing_place_id", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "AssignMissing Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "AssignMissing Place "+id, "")
		t.Cleanup(func() { config.DB.Delete(&place) })
		user := helpers.CreateTestUser(t, "assign_miss_"+id+"@test.com", "assign_miss_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		pu := helpers.CreateTestPlaceUser(t, place.ID, user.ID)
		t.Cleanup(func() { config.DB.Delete(&pu) })

		payload := map[string]interface{}{
			"transfer": false,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/admin/users/"+user.ID.String()+"/place", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(user.ID.String())
		c.Set("company_id", company.ID.String())

		handlers.AssignUserPlace(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "deluser_ok_"+id+"@test.com", "deluser_ok_"+id)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&user) })

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
		config.DB.Unscoped().Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
		if count != 0 {
			t.Error("Expected user to be deleted")
		}
	})

	t.Run("with_reservations", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "deluser_res_"+id+"@test.com", "deluser_res_"+id)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&user) })
		resource := helpers.CreateTestResource(t, "DeleteUser Resource "+id)
		t.Cleanup(func() { config.DB.Delete(&resource) })
		reservation := helpers.CreateTestReservation(t, user.ID, resource.ID)
		t.Cleanup(func() { config.DB.Unscoped().Delete(&reservation) })

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

		var userCount int64
		config.DB.Unscoped().Model(&models.User{}).Where("id = ?", user.ID).Count(&userCount)
		if userCount != 0 {
			t.Error("Expected user to be deleted")
		}

		var resCount int64
		config.DB.Unscoped().Model(&models.Reservation{}).Where("id = ?", reservation.ID).Count(&resCount)
		if resCount != 0 {
			t.Error("Expected reservation to be deleted")
		}
	})

	t.Run("user_is_admin", func(t *testing.T) {
		id := uuid.New().String()[:8]
		admin := helpers.CreateTestAdmin(t, "deluser_admin_"+id+"@test.com", "deladmin_"+id)
		t.Cleanup(func() { config.DB.Delete(&admin) })

		req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+admin.ID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(admin.ID.String())

		handlers.DeleteUser(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/admin/user/"+uuid.New().String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		handlers.DeleteUser(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
