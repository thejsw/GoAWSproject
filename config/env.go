package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env load failed")
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Fatal("OPENAI_API_KEY missing")
	}
}
