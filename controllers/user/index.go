package user

import (
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/db"
	"github.com/fachrunwira/basic-go-api-template/lib/response"

	"github.com/labstack/echo/v4"
)

func (h *userHandler) Home(c echo.Context) error {
	rows, err := db.DB.Query("SELECT id, nama_admin FROM admin")
	if err != nil {
		h.AppLogger.Printf("Failed to fetch data %v", err)
		return response.InternalError(c, "Unknown error occured", "Internal server error")
	}

	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			h.AppLogger.Printf("Internal server error %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"status":  false,
				"code":    http.StatusInternalServerError,
				"message": "Terjadi kesalahan.",
			})
		}

		result = append(result, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"code":    http.StatusOK,
		"message": "success",
		"content": result,
	})
}
