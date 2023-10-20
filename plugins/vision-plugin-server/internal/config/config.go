package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host           string
	Port           string
	DatabaseURL    string
	AllowedOrigins []string
}

func New() Config {
	host := getEnvDefault("HOST", "localhost")
	port := getEnvDefault("PORT", "8080")
	return Config{
		Host:        host,
		Port:        port,
		DatabaseURL: getEnvDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		AllowedOrigins: strings.Split(
			getEnvDefault("ALLOWED_ORIGINS", fmt.Sprintf("http://%s:%s,https://%s:%s", host, port, host, port)), ",",
		),
	}
}

func getEnvDefault(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}
