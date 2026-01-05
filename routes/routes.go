package routes

import (
	"spacebook/handlers"
	"spacebook/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "API is running",
		})
	})

	// =====================
	// Public / User routes
	// =====================

	e.GET("/resources", handlers.GetResources)
	e.POST("/reservations", handlers.CreateReservation)
	e.GET("/notifications", handlers.GetUserNotifications)

	// =====================
	// Admin routes
	// =====================

	admin := e.Group("/admin")
	admin.Use(middleware.AdminOnly)

	// Resources
	admin.POST("/resources", handlers.CreateResource)
	admin.DELETE("/resources/:id", handlers.DeleteResource)

	// Reservations
	admin.GET("/reservations", handlers.GetAdminReservations)
	admin.PUT("/reservations/:id/approve", handlers.ApproveReservation)

	// Notifications
	admin.GET("/notifications", handlers.GetAdminNotifications)
	admin.PUT("/notifications/:id/read", handlers.MarkNotificationAsRead)
}
