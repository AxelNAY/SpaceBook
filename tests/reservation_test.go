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
	"spacebook/helpers"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestCreateReservation(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	user := helpers.CreateTestUser(t, "reservationtest@test.com", "reservationtest")
	resource := helpers.CreateTestResource(t, "Test Resource")
	defer func() {
		config.DB.Where("ressource_id = ?", resource.ID).Delete(&models.Reservation{})
		config.DB.Delete(&resource)
		config.DB.Delete(&user)
	}()

	t.Run("successful reservation", func(t *testing.T) {
		startDatetime := time.Now().Add(24 * time.Hour)
		endDatetime := startDatetime.Add(1 * time.Hour)

		payload := map[string]interface{}{
			"user_id":        user.ID.String(),
			"ressource_id":   resource.ID.String(),
			"start_datetime": startDatetime.Format(time.RFC3339),
			"end_datetime":   endDatetime.Format(time.RFC3339),
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

	t.Run("invalid dates - start after end", func(t *testing.T) {
		startDatetime := time.Now().Add(48 * time.Hour)
		endDatetime := time.Now().Add(24 * time.Hour)

		payload := map[string]interface{}{
			"user_id":        user.ID.String(),
			"ressource_id":   resource.ID.String(),
			"start_datetime": startDatetime.Format(time.RFC3339),
			"end_datetime":   endDatetime.Format(time.RFC3339),
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

	t.Run("non-existent resource", func(t *testing.T) {
		startDatetime := time.Now().Add(24 * time.Hour)
		endDatetime := startDatetime.Add(1 * time.Hour)

		payload := map[string]interface{}{
			"user_id":        user.ID.String(),
			"ressource_id":   uuid.New().String(),
			"start_datetime": startDatetime.Format(time.RFC3339),
			"end_datetime":   endDatetime.Format(time.RFC3339),
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

func TestOverlapCheck(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	user := helpers.CreateTestUser(t, "overlaptest@test.com", "overlaptest")
	resource := helpers.CreateTestResource(t, "Overlap Test Resource")
	defer func() {
		config.DB.Where("ressource_id = ?", resource.ID).Delete(&models.Reservation{})
		config.DB.Delete(&resource)
		config.DB.Delete(&user)
	}()

	startDatetime := time.Now().Add(72 * time.Hour)
	endDatetime := startDatetime.Add(1 * time.Hour)

	t.Run("first reservation - should succeed", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":        user.ID.String(),
			"ressource_id":   resource.ID.String(),
			"start_datetime": startDatetime.Format(time.RFC3339),
			"end_datetime":   endDatetime.Format(time.RFC3339),
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

	t.Run("second overlapping reservation - should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"user_id":        user.ID.String(),
			"ressource_id":   resource.ID.String(),
			"start_datetime": startDatetime.Format(time.RFC3339),
			"end_datetime":   endDatetime.Format(time.RFC3339),
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
	})

	t.Run("reservation at different time - should succeed", func(t *testing.T) {
		differentStart := startDatetime.Add(2 * time.Hour)
		differentEnd := differentStart.Add(1 * time.Hour)

		payload := map[string]interface{}{
			"user_id":        user.ID.String(),
			"ressource_id":   resource.ID.String(),
			"start_datetime": differentStart.Format(time.RFC3339),
			"end_datetime":   differentEnd.Format(time.RFC3339),
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
