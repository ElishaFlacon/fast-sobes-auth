package config

import (
	"os"
	"strconv"
)

type GrpcConfig struct {
	Port string
}

type PgConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type LoggerConfig struct {
	BufSize int
}

type Config struct {
	Env string

	Grpc *GrpcConfig
	Pg   *PgConfig
	Log  *LoggerConfig
}

func NewConfig() *Config {
	env := getEnv("APP_ENV", "dev")

	grpcCfg := &GrpcConfig{
		Port: getEnv("GRPC_PORT", "5000"),
	}

	pgCfg := &PgConfig{
		Host:     getEnv("POSTGRES_HOST", "fast-sobes-auth-db"),
		Port:     getIntEnv("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnv("POSTGRES_DB", "auth"),
		SSLMode:  getEnv("POSTGRES_SSL", "disable"),
	}

	logCfg := &LoggerConfig{
		BufSize: 10,
	}

	return &Config{
		Env: env,

		Grpc: grpcCfg,
		Pg:   pgCfg,
		Log:  logCfg,
	}
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getIntEnv(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}

	return intVal
}
