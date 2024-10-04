package jwtHandlers

import (
	"IAM/initializers"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"time"
)

func CreateJWT(email string, userVersion string, rdb *redis.Client) (string, error) {
	// get userID from redis
	repo := &redisAuxiliaryHandlers.RedisUserIDByEmailRepo{RDB: rdb}
	userID, err := auxiliary.UserIDByEmail(repo, email)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := models.Claims{
		UserID:      userID,
		UserVersion: userVersion,
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
