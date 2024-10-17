package middlewareJWT

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JwtRepository interface {
	GetDataForJWT(email string) (userID string, userVersion string, error error)
}

var JwtRepo JwtRepository

func CheckTokenValid(tokenString string) (isValid bool, userData []string, err error) {
	// check signature
	claims, err := ValidateTokenSignature(tokenString)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return false, userData, err
	}
	email := claims.Email
	userID := claims.UserID
	// get current userVersion from DB
	_, currentUserVersion, err := JwtRepo.GetDataForJWT(email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return false, userData, err
	}

	// compare
	if claims.UserVersion != currentUserVersion {
		return false, userData, nil
	}

	userData = append(userData, email)
	userData = append(userData, userID)
	userData = append(userData, tokenString)

	return true, userData, nil
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
		logs.Error.Println("failed parse jwt", err)
		logs.ErrorLogger.Error("failed parse jwt", err)
		return nil, err
	}

	// check claims
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
