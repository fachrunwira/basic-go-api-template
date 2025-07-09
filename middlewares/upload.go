package middlewares

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func AllowedFileTypes(maxSize int64, allowed_ext []string) echo.MiddlewareFunc {
	fileTypesMap := make(map[string]bool, len(allowed_ext))
	for _, ext := range allowed_ext {
		fileTypesMap[strings.ToLower(ext)] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, maxSize)

			// Optional: parse form early to detect file
			err := c.Request().ParseMultipartForm(maxSize)
			if err != nil {
				return c.JSON(http.StatusRequestEntityTooLarge, echo.Map{
					"message": "Uploaded file too large.",
					"errors":  "file too large",
				})
			}

			form := c.Request().MultipartForm
			for _, files := range form.File {
				for _, file := range files {
					ext := strings.ToLower(filepath.Ext(file.Filename))
					if !fileTypesMap[ext] {
						return c.JSON(http.StatusUnsupportedMediaType, echo.Map{
							"message": "Selected file types is not allowed.",
							"errors":  "file types: " + ext,
						})
					}
				}
			}

			return next(c)
		}
	}
}
