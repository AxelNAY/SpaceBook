package tests

import (
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

func TestGetUserNotifications(t *testing.T) {
	helpers.SetupTestDB()

	e := echo.New()

	user := helpers.CreateTestUser(t, "notiftest@test.com", "notiftest")
	notification := helpers.CreateTestNotification(t, &user.ID, "Test notification")
	defer func() {
		config.DB.Delete(&notification)
		config.DB.Delete(&user)
	}()

	t.Run("get user notifications", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notifications?userId="+user.ID.String(), nil)
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
	helpers.SetupTestDB()

	e := echo.New()

	notification := helpers.CreateTestNotification(t, nil, "Admin test notification")
	defer config.DB.Delete(&notification)

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
	helpers.SetupTestDB()

	e := echo.New()

	notification := helpers.CreateTestNotification(t, nil, "Mark as read test")
	defer config.DB.Delete(&notification)

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
