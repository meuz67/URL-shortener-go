package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Address    string
	BaseURL    string
	Database   DatabaseConfig
	Encryption EncryptionConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type EncryptionConfig struct {
	Secret string
}

func Load() Config {
	_ = godotenv.Load()

	address := getEnv("APP_ADDRESS", ":8080")

	return Config{
		Address: address,
		BaseURL: getEnv("BASE_URL", "http://localhost"+address),
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "sar58yeaf"),
			Name:     getEnv("POSTGRES_DB", "password"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		Encryption: EncryptionConfig{
			Secret: os.Getenv("ENCRYPTION_SECRET"),
		},
	}
}

func (c Config) Validate() error {
	if c.Encryption.Secret == "" {
		return errors.New("ENCRYPTION_SECRET is required")
	}

	return nil
}

func (c DatabaseConfig) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + c.Port,
		Path:   c.Name,
	}

	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
