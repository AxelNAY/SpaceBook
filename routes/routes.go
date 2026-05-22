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

	e.GET("/companies", handlers.GetCompanies)
	e.GET("/places", handlers.GetPlaces)
	e.GET("/categories", handlers.GetCategories)

	// =====================
	// Protected routes (authenticated users)
	// =====================

	protected := e.Group("")
	protected.Use(middleware.JWTAuth)

	protected.GET("/resources", handlers.GetResources)
	protected.POST("/reservation", handlers.CreateSimpleReservation)
	protected.POST("/reservation_group", handlers.CreateGroupReservation)
	protected.GET("/reservations", handlers.GetUserReservations)
	protected.GET("/notifications", handlers.GetUserNotifications)

	// =====================
	// Superadmin routes (authenticated + superadmin role)
	// =====================

	superadmin := e.Group("/superadmin")
	superadmin.Use(middleware.JWTAuth)
	superadmin.Use(middleware.SuperadminOnly)

	// Companies
	superadmin.POST("/companies", handlers.CreateCompany)
	superadmin.PUT("/companies/:id", handlers.UpdateCompany)
	superadmin.DELETE("/companies/:id", handlers.DeleteCompany)

	// Places
	superadmin.POST("/places", handlers.CreatePlace)
	superadmin.PUT("/places/:id", handlers.UpdatePlace)
	superadmin.DELETE("/places/:id", handlers.DeletePlace)

	// =====================
	// Admin routes (authenticated + admin role, scopé à l'entreprise)
	// =====================

	admin := e.Group("/admin")
	admin.Use(middleware.JWTAuth)
	admin.Use(middleware.AdminOnly)

	// Places (scopées à l'entreprise de l'admin)
	admin.GET("/places", handlers.AdminGetPlaces)
	admin.GET("/places/:id", handlers.AdminGetPlaceDetails)
	admin.POST("/places", handlers.AdminCreatePlace)
	admin.PUT("/places/:id", handlers.AdminUpdatePlace)
	admin.DELETE("/places/:id", handlers.AdminDeletePlace)

	// Categories
	admin.POST("/categories", handlers.CreateCategory)
	admin.PUT("/categories/:id", handlers.UpdateCategory)
	admin.DELETE("/categories/:id", handlers.DeleteCategory)

	// Resources (scopées à l'entreprise de l'admin)
	admin.GET("/resources", handlers.AdminGetResources)
	admin.POST("/resources", handlers.AdminCreateResource)
	admin.PUT("/resources/:id", handlers.AdminUpdateResource)
	admin.DELETE("/resources/:id", handlers.AdminDeleteResource)

	// Reservations (scopées à l'entreprise de l'admin)
	admin.GET("/reservations", handlers.GetAdminReservations)
	admin.PUT("/reservations/:id/approve", handlers.ApproveReservation)
	admin.PUT("/reservations/:id/reject", handlers.RejectReservation)
	admin.PUT("/reservations/:id/resources/approve", handlers.ApproveGroupReservationResources)
	admin.PUT("/reservations/:id/resources/reject", handlers.RejectGroupReservationResources)

	// Users (scopés à l'entreprise de l'admin)
	admin.GET("/users", handlers.GetUsers)
	admin.PUT("/users/accept/:id", handlers.AcceptUser)
	admin.PUT("/users/refuse/:id", handlers.RefuseUser)
	admin.PUT("/users/promote/:id", handlers.PromoteUser)
	admin.PUT("/users/:id/place", handlers.AssignUserPlace)
	admin.DELETE("/user/:id", handlers.DeleteUser)

	// Notifications
	admin.GET("/notifications", handlers.GetAdminNotifications)
	admin.PUT("/notifications/:id/read", handlers.MarkNotificationAsRead)
}
