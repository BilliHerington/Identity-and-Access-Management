package handlers

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func EmailMatch(c *gin.Context, email string) (bool, error) {
	ctx := context.Background()

	// Проверка наличия email в Redis
	emailKey := "email:" + email
	_, err := initializers.Rdb.Get(ctx, emailKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
func GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	userID, err := initializers.Rdb.Get(ctx, "email:"+email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("email %s not found", email)
		}
		logs.Error.Println(err)
		return "", err
	}
	return userID, nil
}
