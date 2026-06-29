package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
}


func Load() *Config {

	_ = godotenv.Load()

	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-me"),
	}

	if cfg.DatabaseURL == "" {
		
		host := getEnv("DB_HOST")
		port := getEnv("DB_PORT")
		user := getEnv("DB_USER")
		pass := getEnv("DB_PASSWORD")
		name := getEnv("DB_NAME")
		ssl := getEnv("DB_SSLMODE", "disable")
		cfg.DatabaseURL = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Bangkok",
			host, port, user, pass, name, ssl,
		)
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
