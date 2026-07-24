package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	databaseURL string
	redisAddr   string
	jwtSecret   string
	httpPort    string
}

func loadConfig() (config, error) {
	// Ignored deliberately: in Docker the env comes from compose, not a file.
	_ = godotenv.Load()

	cfg := config{
		databaseURL: os.Getenv("DATABASE_URL"),
		redisAddr:   os.Getenv("REDIS_ADDR"),
		jwtSecret:   os.Getenv("JWT_SECRET"),
		httpPort:    os.Getenv("HTTP_PORT"),
	}

	for name, value := range map[string]string{
		"DATABASE_URL": cfg.databaseURL,
		"REDIS_ADDR":   cfg.redisAddr,
		"JWT_SECRET":   cfg.jwtSecret,
	} {
		if value == "" {
			return config{}, fmt.Errorf("%s is required", name)
		}
	}

	if cfg.httpPort == "" {
		cfg.httpPort = "8080"
	}

	return cfg, nil
}
