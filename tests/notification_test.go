package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"spacebook/config"
	"spacebook/handlers"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestGetUserNotifications(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Create test user
	userID := uuid.New()
	user := models.User{
		ID:       userID,
		Email:    "notiftest@test.com",
		Username: "notiftest",
		Password: []byte("password"),
		Role:     "user",
	}
	config.DB.Create(&user)

	// Create test notification for user
	notification := models.Notification{
		UserID:  &userID,
		Type:    "reservation",
		Message: "Test notification",
		IsRead:  false,
	}
	config.DB.Create(&notification)

	defer func() {
		config.DB.Where("user_id = ?", userID).Delete(&models.Notification{})
		config.DB.Where("id = ?", userID).Delete(&models.User{})
	}()

	t.Run("get user notifications", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notifications?userId="+userID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetUserNotifications(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var notifications []models.Notification
		json.Unmarshal(rec.Body.Bytes(), &notifications)

		if len(notifications) == 0 {
			t.Error("Expected at least one notification")
		}
	})

	t.Run("missing userId parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.GetUserNotifications(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetAdminNotifications(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Create test notification without user (admin notification)
	notification := models.Notification{
		Type:    "reservation",
		Message: "Admin test notification",
		IsRead:  false,
	}
	config.DB.Create(&notification)

	defer config.DB.Where("message = ?", "Admin test notification").Delete(&models.Notification{})

	t.Run("get all notifications", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/notifications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetAdminNotifications(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var notifications []models.Notification
		json.Unmarshal(rec.Body.Bytes(), &notifications)

		if len(notifications) == 0 {
			t.Error("Expected at least one notification")
		}
	})
}

func TestMarkNotificationAsRead(t *testing.T) {
	setupTestDB()

	e := echo.New()

	// Create test notification
	notification := models.Notification{
		Type:    "reservation",
		Message: "Mark as read test",
		IsRead:  false,
	}
	config.DB.Create(&notification)

	defer config.DB.Where("message = ?", "Mark as read test").Delete(&models.Notification{})

	t.Run("mark notification as read", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/admin/notifications/"+notification.ID.String()+"/read", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(notification.ID.String())

		err := handlers.MarkNotificationAsRead(c)
		if err != nil {
			t.Fatalf("Handler returned error: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		// Verify notification is marked as read
		var updatedNotif models.Notification
		config.DB.First(&updatedNotif, "id = ?", notification.ID)

		if !updatedNotif.IsRead {
			t.Error("Expected notification to be marked as read")
		}
	})

	t.Run("non-existent notification", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/admin/notifications/"+uuid.New().String()+"/read", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		handlers.MarkNotificationAsRead(c)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}
