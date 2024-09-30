package jwtHandlers

import (
	"IAM/initializers"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

func CreateJWT(c *gin.Context, email string, rdb *redis.Client) (string, error) {
	// get userID from redis
	id, err := auxiliary.GetUserIDByEmail(c, email, rdb)
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

	// sign token with secret key
	signedToken, err := token.SignedString(initializers.JwtSecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
