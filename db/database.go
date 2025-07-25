package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/fachrunwira/basic-go-api-template/config"
	"github.com/fachrunwira/basic-go-api-template/lib/logger"
	"github.com/fachrunwira/basic-go-api-template/lib/storage"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type dbKeyType string

const dbKey dbKeyType = "database"

var (
	DBLogger *log.Logger = logger.SetLogger("./storage/log/db.log")
	DB       *sql.DB
)

func InitDB() {
	db_cfg := config.LoadDBConfig()

	DB, err := connectTo(&db_cfg)

	if err != nil {
		DBLogger.Fatalf("%v", err)
	}

	err = DB.Ping()
	if err != nil {
		DBLogger.Fatalf("Failed to ping DB: %v", err)
	}
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
	case "sqlite":
		st := storage.Init().StoragePath().Directory("databases").Get("database.db")
		if st.Error() != nil {
			return nil, st.Error()
		}
		url, err := st.Url()
		if err != nil {
			return nil, err
		}

		if db, err = sql.Open("sqlite3", *url); err != nil {
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

func InjectDB(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbKey, DB)
}

func FromContext(ctx context.Context) *sql.DB {
	db, ok := ctx.Value(dbKey).(*sql.DB)
	if !ok {
		DBLogger.Fatal("database not found in context")
	}

	return db
}
