package authJWT

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtRepository interface {
	GetDataForJWT(email string) (userID string, userVersion string, error error)
	SetJWT(userID, jwt string) error
}

var JwtRepo JwtRepository

func CreateJWT(email, userID, userVersion string) (string, error) {

	expirationTime := time.Now().Add(time.Hour * 24)
	claims := models.Claims{
		Email:          email,
		UserID:         userID,
		UserVersion:    userVersion,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token with secret key
	signedToken, err := token.SignedString(initializers.JwtSecretKey)
	if err != nil {
		logs.Error.Println("failed sign token", err)
		logs.ErrorLogger.Error("failed sign token", err)
		return "", err
	}
	return signedToken, nil
}
