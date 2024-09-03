package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnvVariables() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
