package config

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
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		Name:     getEnv("DB_NAME", "go_database"),
		User:     getEnv("DB_USER", "gobase"),
		Password: getEnv("DB_PASSWORD", ""),
		Driver:   getEnv("DB_CONNECTION", "mysql"),
	}
}
