package user

import (
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/db/query"
	"github.com/fachrunwira/basic-go-api-template/lib/passwordhash"
	"github.com/fachrunwira/basic-go-api-template/lib/response"
	"github.com/fachrunwira/basic-go-api-template/lib/validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type userRegisterDTO struct {
	Name            string `json:"name" validate:"required,max_string=90"`
	Email           string `json:"email" validate:"required,email,max_string=100"`
	Age             int    `json:"age" validate:"required,number"`
	Password        string `json:"password" validate:"required,min_string=8,max_string=64"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min_string=8,max_string=64,eqfield=Password"`
}

func (h *userHandler) Register(c echo.Context) (err error) {
	var registerDTO userRegisterDTO
	ctx := c.Request().Context()

	if err := c.Bind(&registerDTO); err != nil {
		h.AppLogger.Printf("RegisterDTO, Error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
			"errors":  nil,
		})
	}

	if err = c.Validate(&registerDTO); err != nil {
		errMsg := validation.Errors(err.(validator.ValidationErrors)).(string)
		return response.FailedValidation(c, errMsg, nil)
	}

	hashedPassword, err := passwordhash.Make(registerDTO.Password)
	if err != nil {
		h.AppLogger.Printf("RegisterUser, Password: %s", err)
		return response.InternalError(c, "Failed to hash password.", "internal server error")
	}

	registerInterface := map[string]interface{}{
		"id":        uuid.New(),
		"name":      registerDTO.Name,
		"email":     registerDTO.Email,
		"age":       registerDTO.Age,
		"passwords": hashedPassword,
	}

	if err = query.NewQuery(ctx).Table("users").Insert(registerInterface); err != nil {
		h.AppLogger.Printf("RegisterUser, Insert: %s", err)
		return response.InternalError(c, "Failed to save record.", "internal server error")
	}

	var contentResponse = make([]string, 0)
	return response.NewRecord(c, "Success inserting new record", contentResponse)
}
