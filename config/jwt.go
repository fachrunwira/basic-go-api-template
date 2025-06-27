package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	Key []byte
	Iat int64
	Exp int64
}

func LoadJWTConfig() JWTConfig {
	return JWTConfig{
		Key: []byte(os.Getenv("JWT_SECRET")),
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(time.Minute * 15).Unix(),
	}
}
