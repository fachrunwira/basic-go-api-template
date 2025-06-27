package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fachrunwira/basic-go-api-template/lib/jwt"
	"github.com/fachrunwira/basic-go-api-template/lib/response"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	AppLogger *log.Logger
	DB        *sql.DB
}

func NewUserHandler(logger *log.Logger, db *sql.DB) *userHandler {
	return &userHandler{
		AppLogger: logger,
		DB:        db,
	}
}

func (h *userHandler) Home(c echo.Context) error {
	rows, err := h.DB.Query("SELECT id, nama_admin FROM admin")
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

// type Admin struct {
// 	ID         string
// 	Email      string
// 	Nama_admin string
// }

// func Home(c echo.Context) error {

// 	// example execute query
// 	err := db.Connect()
// 	if err != nil {
// 		log.Fatalf("db connection failed: %v", err)
// 	}
// 	defer db.DB.Close()

// 	result, err := db.DB.Query("SELECT id, nama_admin, email FROM admin")
// 	if err != nil {
// 		log.Fatalf("failed to fetch data: %v", err)
// 	}
// 	defer result.Close()

// 	var admins []Admin
// 	for result.Next() {
// 		var admin Admin
// 		err := result.Scan(&admin.ID, &admin.Email, &admin.Nama_admin)
// 		if err != nil {
// 			log.Fatalf("Row scan failed: %v", err)
// 		}
// 		admins = append(admins, admin)
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"status":  true,
// 		"message": "Hello World!",
// 		"content": admins,
// 	})
// }

func GetToken(c echo.Context) error {
	claims := map[string]interface{}{
		"username": "johny",
	}

	token := jwt.GenerateToken(claims)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"code":    http.StatusOK,
		"content": token,
	})
}

// func ParseToken(c echo.Context) error {
// 	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.VFb0qJ1LRg_4ujbZoRMXnVkUgiuKq5KxWqNdbKq_G9Vvz-S1zZa9LPxtHWKa64zDl2ofkT8F6jBt_K4riU-fPg"
// 	parsedToken, err := jwt.ValidateToken(token)

// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 			"status": false,
// 			"code":   http.StatusUnauthorized,
// 			"errors": err,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"status":  true,
// 		"code":    http.StatusOK,
// 		"content": parsedToken,
// 	})
// }
