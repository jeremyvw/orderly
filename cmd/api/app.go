package main

import (
	"log"

	"orderly/internal/pkg/token"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := token.Init([]byte(cfg.jwtSecret)); err != nil {
		log.Fatalf("init token: %v", err)
	}

	db, err := initPostgres(cfg.databaseURL)
	if err != nil {
		log.Fatalf("init postgres: %v", err)
	}
	defer db.Close()

	// To Do: Initialize Redis, repositories, services, and start the HTTP server.
}
