package main

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() error {

	// llamada al paquete, lee por defecto el .env
	err := godotenv.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func validateEnvironments() error {
	if strings.TrimSpace(os.Getenv("SERVER_PORT")) == "" {
		return errors.New("the SERVER_PORT env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS")) == "" {
		return errors.New("the ALLOWED_ORIGINS env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("ALLOWED_METHODS")) == "" {
		return errors.New("the ALLOWED_METHODS env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("IMAGES_DIR")) == "" {
		return errors.New("the IMAGES_DIR env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("JWT_SECRET_KEY")) == "" {
		return errors.New("the JWT_SECRET_KEY env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("DB_USER")) == "" {
		return errors.New("the DB_USER env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("DB_PASSWORD")) == "" {
		return errors.New("the DB_PASSWORD env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("DB_HOST")) == "" {
		return errors.New("the DB_HOST env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("DB_PORT")) == "" {
		return errors.New("the DB_PORT env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("DB_NAME")) == "" {
		return errors.New("the DB_NAME env var is mandatory")
	}

	if strings.TrimSpace(os.Getenv("DB_SSL_MODE")) == "" {
		return errors.New("the DB_SSL_MODE env var is mandatory")
	}

	return nil
}

// IsSemanticSearchEnabled returns whether the ENABLE_SEMANTIC_SEARCH env var
// is explicitly set to "true". Defaults to false (semantic search off).
// This is an optional env var — it is NOT validated in validateEnvironments().
func IsSemanticSearchEnabled() bool {
	return strings.ToLower(os.Getenv("ENABLE_SEMANTIC_SEARCH")) == "true"
}
