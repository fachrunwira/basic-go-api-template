package config

import (
	"fmt"
	"os"
)

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var intVal int
		fmt.Scanf(val, "%d", intVal)
		return intVal
	}

	return defaultVal
}
