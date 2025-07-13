package config

import "github.com/fachrunwira/basic-go-api-template/lib/env"

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     env.Get("DB_HOST", "localhost"),
		Port:     env.Get("DB_PORT", "3306"),
		Name:     env.Get("DB_NAME", "go_database"),
		User:     env.Get("DB_USER", "gobase"),
		Password: env.Get("DB_PASSWORD", ""),
		Driver:   env.Get("DB_CONNECTION", "mysql"),
	}
}
