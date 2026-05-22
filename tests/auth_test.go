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

func TestRegister(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "Company "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "Place "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		email := "register_ok_" + id + "@test.com"
		payload := map[string]interface{}{
			"email":    email,
			"username": "regok_" + id,
			"password": "password123",
			"place_id": place.ID.String(),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.Register(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		if _, ok := response["message"]; !ok {
			t.Error("Expected 'message' field in response")
		}
		if _, ok := response["user"]; !ok {
			t.Error("Expected 'user' field in response")
		}

		t.Cleanup(func() {
			config.DB.Where("email = ?", email).Delete(&models.User{})
		})
	})

	t.Run("missing_place_id", func(t *testing.T) {
		id := uuid.New().String()[:8]
		payload := map[string]interface{}{
			"email":    "reg_nopi_" + id + "@test.com",
			"username": "regnopi_" + id,
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Register(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("missing_required_fields", func(t *testing.T) {
		id := uuid.New().String()[:8]
		payload := map[string]interface{}{
			"email":    "reg_nopw_" + id + "@test.com",
			"username": "regnopw_" + id,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Register(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("duplicate_email", func(t *testing.T) {
		id := uuid.New().String()[:8]
		company := helpers.CreateTestCompany(t, "Company Dup "+id)
		t.Cleanup(func() { config.DB.Delete(&company) })
		place := helpers.CreateTestPlace(t, company.ID, "Place Dup "+id, "desc")
		t.Cleanup(func() { config.DB.Delete(&place) })

		existingUser := helpers.CreateTestUser(t, "dup_"+id+"@test.com", "dupuser_"+id)
		t.Cleanup(func() {
			config.DB.Where("user_id = ?", existingUser.ID).Delete(&models.PlaceUser{})
			config.DB.Delete(&existingUser)
		})

		payload := map[string]interface{}{
			"email":    existingUser.Email,
			"username": "other_" + id,
			"password": "password123",
			"place_id": place.ID.String(),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Register(c)

		if rec.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, rec.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("success_active_user", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "login_ok_"+id+"@test.com", "loginok_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		payload := map[string]string{
			"email":    user.Email,
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.Login(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d — body: %s", http.StatusOK, rec.Code, rec.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		if _, ok := response["token"]; !ok {
			t.Error("Expected 'token' field in response")
		}
	})

	t.Run("pending_user_blocked", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestPendingUser(t, "login_pend_"+id+"@test.com", "loginpend_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		payload := map[string]string{
			"email":    user.Email,
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Login(c)

		if rec.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("wrong_password", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "login_wpw_"+id+"@test.com", "loginwpw_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })

		payload := map[string]string{
			"email":    user.Email,
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Login(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("nonexistent_user", func(t *testing.T) {
		payload := map[string]string{
			"email":    "nonexistent_" + uuid.New().String()[:8] + "@test.com",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Login(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("missing_fields", func(t *testing.T) {
		payload := map[string]string{
			"email": "missing@test.com",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Login(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}
