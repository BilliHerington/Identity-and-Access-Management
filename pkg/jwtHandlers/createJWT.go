package jwtHandlers

import (
	"IAM/initializers"
	"IAM/pkg/handlers"
	"IAM/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateJWT(c *gin.Context, email string) (string, error) {
	id, err := handlers.GetUserIDByEmail(c, email)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := models.Claims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа
	signedToken, err := token.SignedString(initializers.JwtSecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
