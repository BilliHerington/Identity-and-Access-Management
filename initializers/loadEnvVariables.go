package initializers

import (
	"github.com/joho/godotenv"
	"log"

	"os"
)

var JwtSecretKey []byte

func LoadEnvVariables() {
	envPath := "config/.env"
	err := godotenv.Load(envPath)
	JwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		log.Fatal("Error loading .env file. Tried load .env from: \nPlease, make sure that you have configured your .env file\n", envPath)
	}
	if os.Getenv("USE_TEST_MODE_WITHOUT_GOOGLE") == "true" {
		log.Printf("Warning! You using project without Google Account, some functional may not be able. \n You can change this in .env\n")
	}
}
