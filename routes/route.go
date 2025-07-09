package routes

import (
	"log"
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/lib/logger"

	"github.com/labstack/echo/v4"
)

var AppLogger *log.Logger = logger.SetLogger("./storage/log/app.log")

func RegisterRoutes(e *echo.Echo) {

	registerUserRoutes(e)

	// Route for not found
	e.Any("/:any", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status":  false,
			"code":    http.StatusNotFound,
			"message": "route not found",
		})
	})
}
