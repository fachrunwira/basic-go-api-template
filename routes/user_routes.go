package routes

import (
	"github.com/fachrunwira/basic-go-api-template/handlers/user"
	"github.com/labstack/echo/v4"
)

func registerUserRoutes(e *echo.Echo) {
	userHandler := user.NewUserHandler(AppLogger)

	userGroup := e.Group("/user")
	userGroup.GET("", userHandler.Home)
	userGroup.GET("/list", userHandler.ListUser)
	userGroup.POST("/register", userHandler.Register)
}
