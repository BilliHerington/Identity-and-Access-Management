package jwtHandlers

import (
	"IAM/initializers"
	"IAM/pkg/handlers/auxiliary"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

func IsTokenValid(tokenString string, rdb *redis.Client) (bool, string, error) {

	// check signature
	claims, err := ValidateTokenSignature(tokenString)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return false, "", err
	}
	userID := claims.UserID
	// get current userVersion from redis
	currentUserVersion, err := auxiliary.UserVersion(&redisAuxiliaryHandlers.RedisAuxiliaryRepository{RDB: rdb}, userID)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return false, "", err
	}

	// compare
	if claims.UserVersion != currentUserVersion {
		return false, "", nil
	}
	return true, userID, nil
}
func ValidateTokenSignature(tokenString string) (*models.Claims, error) {

	// parse token
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// check signed method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return initializers.JwtSecretKey, nil
	})
	if err != nil {
		logs.Error.Println(err)
		logs.Error.Println(tokenString, err)
		return nil, err
	}

	// check claims
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
