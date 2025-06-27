package ipwhitelisting

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IPWhitelist(allowedIPs []string) echo.MiddlewareFunc {
	allowed := make(map[string]bool)
	for _, ip := range allowedIPs {
		allowed[ip] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clientIP := c.RealIP()

			// Allowed if in whitelist
			if allowed[clientIP] {
				return next(c)
			}

			// Block Otherwise
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"status": false,
				"code":   http.StatusForbidden,
				"errors": "Access Denied",
			})
		}
	}
}
