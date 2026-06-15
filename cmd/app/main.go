package main

import (
	"log"

	"url-shortener/internal/config"
	"url-shortener/internal/httpserver"
	"url-shortener/internal/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewMemoryURLStore()
	server := httpserver.New(cfg, store)

	if err := server.Start(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
