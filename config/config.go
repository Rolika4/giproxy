package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv("GIT_TOKEN")
}