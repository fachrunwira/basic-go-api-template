package user

import (
	"github.com/fachrunwira/basic-go-api-template/db/query"
	"github.com/fachrunwira/basic-go-api-template/lib/response"
	"github.com/labstack/echo/v4"
)

func (h *userHandler) Home(c echo.Context) error {
	ctx := c.Request().Context()

	data, err := query.Builder(ctx).
		Table("users").
		Select("id", "name", "age", "email").
		Get()
	if err != nil {
		h.AppLogger.Printf("UserHandler, GetAll: %v", err)
		return response.InternalError(c, "error while fetching data.", "internal server error")
	}

	if len(data) == 0 {
		return response.Success(c, "Belum ada user", data)
	}

	return response.Success(c, "User ditemukan", data)
}
