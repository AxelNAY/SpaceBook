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
	// Auth routes (public)
	// =====================

	auth := e.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)

	// =====================
	// Public routes
	// =====================

	e.GET("/resources", handlers.GetResources)

	// =====================
	// Protected routes (authenticated users)
	// =====================

	protected := e.Group("")
	protected.Use(middleware.JWTAuth)

	protected.POST("/reservations", handlers.CreateReservation)
	protected.GET("/notifications", handlers.GetUserNotifications)

	// =====================
	// Admin routes (authenticated + admin role)
	// =====================

	admin := e.Group("/admin")
	admin.Use(middleware.JWTAuth)
	admin.Use(middleware.AdminOnly)

	// Resources
	admin.POST("/resources", handlers.CreateResource)
	admin.DELETE("/resources/:id", handlers.DeleteResource)

	// Reservations
	admin.GET("/reservations", handlers.GetAdminReservations)
	admin.PUT("/reservations/:id/approve", handlers.ApproveReservation)
	admin.PUT("/reservations/:id/reject", handlers.RejectReservation)

	// Users
	admin.GET("/users", handlers.GetUsers)
	admin.DELETE("/user/:id", handlers.DeleteUser)

	// Notifications
	admin.GET("/notifications", handlers.GetAdminNotifications)
	admin.PUT("/notifications/:id/read", handlers.MarkNotificationAsRead)
}
