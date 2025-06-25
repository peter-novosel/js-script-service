package config

import (
	"os"
)

type Config struct {
	Port         string
	Env          string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "admin"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "secret"),
		DBName:     getEnv("POSTGRES_DB", "scriptdb"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
