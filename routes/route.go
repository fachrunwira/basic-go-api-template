package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/controllers/auth"
	"github.com/fachrunwira/basic-go-api-template/controllers/user"
	"github.com/fachrunwira/basic-go-api-template/lib/logger"

	"github.com/labstack/echo/v4"
)

var AppLogger *log.Logger

func RegisterRoutes(e *echo.Echo, db *sql.DB) {
	AppLogger = logger.SetLogger("./storage/log/app.log")

	authHandler := auth.NewLoginHandler(AppLogger, db)
	userHandler := user.NewUserHandler(AppLogger, db)

	e.GET("/", userHandler.Home)
	e.POST("/auth/login/v2", authHandler.LoginUser)

	// Route for not found
	e.Any("/:any", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]any{
			"status":  false,
			"code":    http.StatusNotFound,
			"message": "route not found",
		})
	})
}
