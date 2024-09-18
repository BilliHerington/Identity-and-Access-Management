package initializers

import (
	"github.com/joho/godotenv"
	"log"

	"os"
)

var JwtSecretKey []byte

func LoadEnvVariables() {
	err := godotenv.Load("config/.env")
	JwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
