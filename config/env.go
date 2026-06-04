package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() error {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env load skipped:", err)
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		return fmt.Errorf("OPENAI_API_KEY missing")
	}

	if os.Getenv("DATABASE_URL") == "" {
		return fmt.Errorf("DATABASE_URL missing")
	}

	if strings.HasPrefix(os.Getenv("DATABASE_URL"), "https://") || strings.HasPrefix(os.Getenv("DATABASE_URL"), "http://") {
		return fmt.Errorf("DATABASE_URL must be a PostgreSQL connection string, not a Supabase project URL")
	}

	return nil
}
