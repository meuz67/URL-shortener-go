package main

import (
	"context"
	"log"

	"url-shortener/internal/config"
	appcrypto "url-shortener/internal/crypto"
	"url-shortener/internal/database"
	"url-shortener/internal/httpserver"
	"url-shortener/internal/storage"
)

func main() {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	ctx := context.Background()

	pool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}
	defer pool.Close()

	if err := database.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	encryptor, err := appcrypto.NewEncryptor(cfg.Encryption.Secret)
	if err != nil {
		log.Fatalf("create encryptor: %v", err)
	}

	store := storage.NewPostgresURLStore(pool, encryptor)
	server := httpserver.New(cfg, store)

	if err := server.Start(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
