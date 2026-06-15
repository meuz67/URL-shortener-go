package config

import "os"

type Config struct {
	Address string
	BaseURL string
}

func Load() Config {
	address := getEnv("APP_ADDRESS", ":8080")

	return Config{
		Address: address,
		BaseURL: getEnv("BASE_URL", "http://localhost"+address),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
