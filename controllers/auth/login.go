package auth

import (
	"database/sql"
	"log"
	"net/http"

	passwordhash "github.com/fachrunwira/basic-go-api-template/lib/password_hash"
	"github.com/fachrunwira/basic-go-api-template/lib/response"
	"github.com/fachrunwira/basic-go-api-template/lib/validation"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	loginHandler struct {
		AppLogger *log.Logger
		DB        *sql.DB
	}

	userLoginDTO struct {
		Email    string `json:"email_login" form:"email_login" validate:"required,email,max_string=150"`
		Password string `json:"pass_login" form:"pass_login" validate:"required,max_string=250"`
	}
)

func NewLoginHandler(logger *log.Logger, db *sql.DB) *loginHandler {
	return &loginHandler{
		AppLogger: logger,
		DB:        db,
	}
}

func (lh *loginHandler) LoginUser(c echo.Context) (err error) {
	var uDTO userLoginDTO

	if err := c.Bind(&uDTO); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
			"errors":  nil,
		})
	}

	if err = c.Validate(&uDTO); err != nil {
		errorMessage := validation.Errors(err.(validator.ValidationErrors)).(string)
		return response.FailedValidation(c, errorMessage, nil)
	}

	rows, err_row := lh.DB.Query("SELECT id_login, password FROM data_user WHERE email = ? LIMIT 1;", uDTO.Email)
	if err_row != nil {
		lh.AppLogger.Fatalf("Failed to fetch data: %s", err_row)
		return response.InternalError(c, "Unknown error occured", "Internal server error")
	}
	defer rows.Close()

	result_row := make(map[string]string)
	for rows.Next() {
		var id, pass string
		err := rows.Scan(&id, &pass)
		if err != nil {
			lh.AppLogger.Fatalf("Failed to fetch data: %s", err)
			return response.InternalError(c, "Terjadi kesalahan saat mengambil data.", "ERRCODE:500")
		}
		result_row["id"] = id
		result_row["pass"] = pass
	}

	if !passwordhash.Check(uDTO.Password, result_row["pass"]) {
		return response.FailedUnknownUser(c, "Email atau Password anda salah.", nil)
	}

	return response.Success(c, "Berhasil", result_row)
}
