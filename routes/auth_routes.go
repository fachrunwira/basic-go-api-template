package routes

import (
	"github.com/fachrunwira/basic-go-api-template/handlers/auth"
	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(e *echo.Echo) {
	authhandler := auth.NewAuthHandler(AppLogger)
	authGroup := e.Group("/auth")
	authGroup.POST("/login", authhandler.Login)
}
