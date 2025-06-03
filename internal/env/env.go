package env

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/michaeldebetaz/unilscrap/internal/assert"
)

func Load() {
	filepath := ".env"
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", filepath, err)
	}

	slog.Debug("Environment variables loaded")
}

func BASE_PATH() string {
	return assert.NotEmpty(os.Getenv("BASE_PATH"))
}

func ORIGIN() string {
	return assert.NotEmpty(os.Getenv("ORIGIN"))
}

func PORT() string {
	return assert.NotEmpty(os.Getenv("PORT"))
}
