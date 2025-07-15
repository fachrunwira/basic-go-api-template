package routes

import (
	"github.com/fachrunwira/basic-go-api-template/handlers/files"
	"github.com/labstack/echo/v4"
)

func registerFilesRouter(c *echo.Echo) {
	fh := files.NewFilesHandler(AppLogger)
	fhGroup := c.Group("/file")
	fhGroup.POST("/upload", fh.UploadFile)
}
