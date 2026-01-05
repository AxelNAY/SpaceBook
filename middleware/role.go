package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Request().Header.Get("X-ROLE")
		if role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Admin access required",
			})
		}
		return next(c)
	}
}
