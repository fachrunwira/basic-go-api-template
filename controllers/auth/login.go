package auth

import (
	"database/sql"
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/db/query"
	"github.com/fachrunwira/basic-go-api-template/lib/response"
	"github.com/fachrunwira/basic-go-api-template/lib/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	loginDTO struct {
		Email    string `json:"email" validate:"required,email,max_string=120"`
		Password string `json:"password" validate:"required,min_string=8,max_string=64"`
	}
)

func (h *authHandler) Login(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var dto loginDTO

	if err := c.Bind(&dto); err != nil {
		h.AppLogger.Printf("LoginDTO, Error: %s", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid body request.",
			"errors":  nil,
		})
	}

	if err = c.Validate(&dto); err != nil {
		errMsg := validation.Errors(err.(validator.ValidationErrors)).(string)
		return response.FailedValidation(c, errMsg, nil)
	}

	user, err := query.NewQuery(ctx).
		Table("users").
		Select("passwords as pass", "age", "name").
		Where("email = ?", dto.Email).
		First()

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Email atau Password salah.",
			"errors":  nil,
		})
	} else if err != nil {
		h.AppLogger.Printf("AuthHandler, Login: %v", err)
		return response.InternalError(c, "Gagal mengambil data", "internal server error")
	}

	return response.Success(c, "ditemukan", user)
}
