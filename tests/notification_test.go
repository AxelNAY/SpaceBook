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

	t.Run("get_user_notifications", func(t *testing.T) {
		id := uuid.New().String()[:8]
		user := helpers.CreateTestUser(t, "notif_"+id+"@test.com", "notif_"+id)
		t.Cleanup(func() { config.DB.Delete(&user) })
		notification := helpers.CreateTestNotification(t, &user.ID, "Test notification "+id)
		t.Cleanup(func() { config.DB.Delete(&notification) })

		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", user.ID.String())

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

	t.Run("missing_auth_context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handlers.GetUserNotifications(c)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, rec.Code)
		}
	})
}

func TestGetAdminNotifications(t *testing.T) {
	helpers.SetupTestDB()
	e := echo.New()

	t.Run("get_admin_notifications", func(t *testing.T) {
		id := uuid.New().String()[:8]
		admin := helpers.CreateTestAdmin(t, "admin_notif_"+id+"@test.com", "admin_notif_"+id)
		t.Cleanup(func() { config.DB.Delete(&admin) })
		notification := helpers.CreateTestNotification(t, &admin.ID, "Admin test notification "+id)
		t.Cleanup(func() { config.DB.Delete(&notification) })

		req := httptest.NewRequest(http.MethodGet, "/admin/notifications", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", admin.ID)

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

	t.Run("mark_notification_as_read", func(t *testing.T) {
		id := uuid.New().String()[:8]
		notification := helpers.CreateTestNotification(t, nil, "Mark as read test "+id)
		t.Cleanup(func() { config.DB.Delete(&notification) })

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

	t.Run("non_existent_notification", func(t *testing.T) {
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
