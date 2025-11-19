package config

import (
	"os"
)

type Config struct {
	GRPCPort string
	DBURL    string
	Env      string

	LogBufSize int
}

func NewConfig() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "5000"),
		DBURL:    getEnv("DB_URL", "postgres://user:pass@localhost:5432/db"),
		Env:      getEnv("APP_ENV", "dev"),

		LogBufSize: 10,
	}
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
