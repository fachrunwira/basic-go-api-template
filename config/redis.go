package config

import "github.com/fachrunwira/basic-go-api-template/lib/env"

type RedisConfig struct {
	Address  string
	Password string
	Port     string
	Database int
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Address:  env.Get("REDIS_HOST", "localhost"),
		Port:     env.Get("REDIS_PORT", "6379"),
		Password: env.Get("REDIS_PASS", ""),
		Database: env.GetInt("REDIS_DB", 0),
	}
}
