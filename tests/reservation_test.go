package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func createTestUser(t *testing.T) models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		ID:       uuid.New(),
		Email:    "reservationtest@test.com",
		Username: "reservationtest",
		Password: hashedPassword,
		Role:     "user",
	}
	config.DB.Create(&user)
	return user
}

func createTestResource(t *testing.T, capacity int) models.Resource {
	resource := models.Resource{
		ID:       uuid.New().String(),
		Name:     "Test Resource",
		Type:     "equipment",
		Category: "printer",
		Capacity: capacity,
		Status:   "available",
	}
	config.DB.Create(&resource)
	return resource
}

func cleanupTestData(userEmail string, resourceName string) {
	config.DB.Where("resource_id IN (SELECT id FROM resources WHERE name = ?)", resourceName).Delete(&models.Reservation{})
	config.DB.Where("name = ?", resourceName).Delete(&models.Resource{})
	config.DB.Where("email = ?", userEmail).Delete(&models.User{})
}

func TestCreateReservation(t *testing.T) {
	setupTestDB()

	e := echo.New()

	user := createTestUser(t)
	resource := createTestResource(t, 2)

	defer cleanupTestData("reservationtest@test.com", "Test Resource")

	// Test successful reservation
	t.Run("successful reservation", func(t *testing.T) {
		startAt := time.Now().Add(24 * time.Hour)
		endAt := startAt.Add(1 * time.Hour)

		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": resource.ID,
			"start_at":    startAt.Format(time.RFC3339),
			"end_at":      endAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateReservation(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, rec.Code, rec.Body.String())
		}

		var response models.Reservation
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response.Status != "pending" {
			t.Errorf("Expected status 'pending', got '%s'", response.Status)
		}
	})

	// Test reservation with invalid dates
	t.Run("invalid dates - start after end", func(t *testing.T) {
		startAt := time.Now().Add(48 * time.Hour)
		endAt := time.Now().Add(24 * time.Hour) // End before start

		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": resource.ID,
			"start_at":    startAt.Format(time.RFC3339),
			"end_at":      endAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateReservation(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	// Test reservation with non-existent resource
	t.Run("non-existent resource", func(t *testing.T) {
		startAt := time.Now().Add(24 * time.Hour)
		endAt := startAt.Add(1 * time.Hour)

		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": uuid.New().String(),
			"start_at":    startAt.Format(time.RFC3339),
			"end_at":      endAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateReservation(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestCapacityCheck(t *testing.T) {
	setupTestDB()

	e := echo.New()

	user := createTestUser(t)
	// Resource with capacity of 2
	resource := createTestResource(t, 2)

	defer cleanupTestData("reservationtest@test.com", "Test Resource")

	startAt := time.Now().Add(72 * time.Hour)
	endAt := startAt.Add(1 * time.Hour)

	// Create first reservation
	t.Run("first reservation - should succeed", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": resource.ID,
			"start_at":    startAt.Format(time.RFC3339),
			"end_at":      endAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateReservation(c)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}
	})

	// Create second reservation - same time slot
	t.Run("second reservation - should succeed (capacity 2)", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": resource.ID,
			"start_at":    startAt.Format(time.RFC3339),
			"end_at":      endAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateReservation(c)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}
	})

	// Create third reservation - should fail (capacity exceeded)
	t.Run("third reservation - should fail (capacity exceeded)", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": resource.ID,
			"start_at":    startAt.Format(time.RFC3339),
			"end_at":      endAt.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateReservation(c)

		if rec.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusConflict, rec.Code, rec.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		if response["error"] != "Resource fully booked for this time slot" {
			t.Errorf("Expected capacity error message, got: %v", response["error"])
		}
	})

	// Reservation at different time slot should succeed
	t.Run("reservation at different time - should succeed", func(t *testing.T) {
		differentStart := startAt.Add(2 * time.Hour)
		differentEnd := differentStart.Add(1 * time.Hour)

		payload := map[string]interface{}{
			"user_id":     user.ID.String(),
			"resource_id": resource.ID,
			"start_at":    differentStart.Format(time.RFC3339),
			"end_at":      differentEnd.Format(time.RFC3339),
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.CreateReservation(c)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}
	})
}
