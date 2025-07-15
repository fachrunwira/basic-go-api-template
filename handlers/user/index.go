package user

import (
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/db/query"
	"github.com/fachrunwira/basic-go-api-template/lib/response"
	"github.com/labstack/echo/v4"
)

type userListDTO struct {
	Page int `query:"page"`
	Size int `query:"size"`
}

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

func (h *userHandler) ListUser(c echo.Context) error {
	ctx := c.Request().Context()
	var listDTO userListDTO

	if err := c.Bind(&listDTO); err != nil {
		h.AppLogger.Printf("UserHandler, ListUser: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid query param",
			"errors":  nil,
		})
	}

	data, err := query.Builder(ctx).
		Table("users").
		Select("id", "name", "email", "age").
		Page(listDTO.Page).
		Limit(listDTO.Size).
		OrderBy("created_at", "desc").
		Paginate(c)

	if err != nil {
		h.AppLogger.Printf("UserHandler, ListUser: %v", err)
		return response.InternalError(c, "error while fetching data.", "internal server error")
	}

	return response.Success(c, "berhasil", data)
}
