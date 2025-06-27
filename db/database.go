package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fachrunwira/basic-go-api-template/config"
	"github.com/fachrunwira/basic-go-api-template/lib/logger"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var DBLogger *log.Logger = logger.SetLogger("./storage/log/db.log")

func InitDB() (*sql.DB, error) {
	db_cfg := config.LoadDBConfig()

	db, err := connectTo(&db_cfg)

	if err != nil {
		DBLogger.Printf("%v", err)
		return nil, err
	}

	db.Close()
	err = db.Ping()
	if err != nil {
		DBLogger.Printf("Failed to ping DB: %v", err)
		return nil, err
	}

	return db, nil
}

func connectTo(cfg *config.DBConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error

	switch cfg.Driver {
	case "pgsql":
		if db, err = sql.Open("postgres", dsn(cfg)); err != nil {
			return nil, fmt.Errorf("failed to open connection to database: %v", err)
		}
	case "mysql":
		if db, err = sql.Open("mysql", dsn(cfg)); err != nil {
			return nil, fmt.Errorf("failed to open connection to database: %v", err)
		}
	default:
		db = nil
		err = fmt.Errorf("the driver '%s' has not been implemented", cfg.Driver)
	}

	return db, err
}

func dsn(cfg *config.DBConfig) string {
	var dsn string

	switch cfg.Driver {
	case "pgsql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	}

	return dsn
}
