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

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func setupTestDB() {
	godotenv.Load("../.env")
	if config.DB == nil {
		config.ConnectDatabase()
	}
}

func TestRegister(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Test successful registration
	t.Run("successful registration", func(t *testing.T) {
		payload := map[string]string{
			"email":    "testregister@test.com",
			"username": "testregister",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.Register(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["token"] == nil {
			t.Error("Expected token in response")
		}

		// Cleanup
		config.DB.Where("email = ?", "testregister@test.com").Delete(&models.User{})
	})

	// Test registration with missing fields
	t.Run("missing fields", func(t *testing.T) {
		payload := map[string]string{
			"email": "test@test.com",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Register(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	// Test registration with duplicate email
	t.Run("duplicate email", func(t *testing.T) {
		// Create first user
		payload := map[string]string{
			"email":    "duplicate@test.com",
			"username": "duplicate1",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handlers.Register(c)

		// Try to create second user with same email
		payload["username"] = "duplicate2"
		body, _ = json.Marshal(payload)

		req = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		handlers.Register(c)

		if rec.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, rec.Code)
		}

		// Cleanup
		config.DB.Where("email = ?", "duplicate@test.com").Delete(&models.User{})
	})
}

func TestLogin(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Create test user
	payload := map[string]string{
		"email":    "testlogin@test.com",
		"username": "testlogin",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handlers.Register(c)

	// Test successful login
	t.Run("successful login", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "testlogin@test.com",
			"password": "password123",
		}
		body, _ := json.Marshal(loginPayload)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.Login(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["token"] == nil {
			t.Error("Expected token in response")
		}
	})

	// Test login with wrong password
	t.Run("wrong password", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "testlogin@test.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginPayload)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Login(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	// Test login with non-existent user
	t.Run("non-existent user", func(t *testing.T) {
		loginPayload := map[string]string{
			"email":    "nonexistent@test.com",
			"password": "password123",
		}
		body, _ := json.Marshal(loginPayload)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.Login(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})

	// Cleanup
	config.DB.Where("email = ?", "testlogin@test.com").Delete(&models.User{})
}
