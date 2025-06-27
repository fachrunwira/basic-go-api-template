package config

type RedisConfig struct {
	Address  string
	Password string
	Port     string
	Database int
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Address:  getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASS", ""),
		Database: getEnvAsInt("REDIS_DB", 0),
	}
}
