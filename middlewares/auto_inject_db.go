package middlewares

import (
	"github.com/fachrunwira/basic-go-api-template/db"
	"github.com/labstack/echo/v4"
)

func InjectDB(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := db.InjectDB(c.Request().Context())
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		return next(c)
	}
}
