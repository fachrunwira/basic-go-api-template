package main

import (
	"log"
	"os"

	"github.com/fachrunwira/basic-go-api-template/db"
	"github.com/fachrunwira/basic-go-api-template/lib/cache"
	"github.com/fachrunwira/basic-go-api-template/lib/logger"
	"github.com/fachrunwira/basic-go-api-template/lib/validation"
	"github.com/fachrunwira/basic-go-api-template/middlewares"
	"github.com/fachrunwira/basic-go-api-template/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	HTTPLogger *log.Logger
	AppLogger  *log.Logger
	DBLogger   *log.Logger
)

func main() {
	// Load .env files
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatalf("Failed to open env files: %v", errEnv)
		return
	}

	// Logger using lumberjack
	HTTPLogger = logger.SetLogger("./storage/log/http.log")
	AppLogger = logger.SetLogger("./storage/log/app.log")
	DBLogger = logger.SetLogger("./storage/log/db.log")

	database, err := db.InitDB()
	if err != nil {
		DBLogger.Printf("Failed to connect to DB: %v", err)
	}

	defer database.Close()

	cache.InitRedis()

	appPort := os.Getenv("APP_PORT")

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: HTTPLogger.Writer(),
		Format: `${time_rfc3339} | ${status} | ${method} | ${uri} | ${latency_human}` + "\n",
	}))

	e.Logger.SetOutput(HTTPLogger.Writer())

	// Set Validation
	e.Validator = validation.New()

	// Whitelist IP
	// allowedIPs := []string{}
	// e.Use(ipwhitelisting.IPWhitelist(allowedIPs))

	// Rate Limiter
	limiter := middlewares.NewClientLimiter(5, 10)
	e.Use(limiter.Middleware())

	e.Use(middleware.CORS())

	AppLogger.Printf("Starting server on port %s", appPort)
	routes.RegisterRoutes(e, database)
	e.Start(":" + appPort)
}
