package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Accès administrateur requis",
			})
		}
		return next(c)
	}
}

func SuperadminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "superadmin" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Accès superadmin requis",
			})
		}
		return next(c)
	}
}
