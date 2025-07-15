package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func FailedValidation(c echo.Context, message string, errors interface{}) (err error) {
	return c.JSON(http.StatusUnprocessableEntity, echo.Map{
		"message": message,
		"errors":  errors,
	})
}

func InternalError(c echo.Context, message, errors interface{}) (err error) {
	return c.JSON(http.StatusInternalServerError, echo.Map{
		"message": message,
		"errors":  errors,
	})
}

func Success(c echo.Context, message string, content any) (err error) {
	return c.JSON(http.StatusOK, echo.Map{
		"message": message,
		"content": content,
	})
}

func FailedUnknownUser(c echo.Context, message string, errors interface{}) (err error) {
	return c.JSON(http.StatusUnauthorized, echo.Map{
		"message": message,
		"errors":  errors,
	})
}

func NewRecord(c echo.Context, message string, content any) (err error) {
	return c.JSON(http.StatusCreated, echo.Map{
		"message": message,
		"content": content,
	})
}

func BadRequest(c echo.Context, message string, errors interface{}) (err error) {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"message": message,
		"errors":  errors,
	})
}
