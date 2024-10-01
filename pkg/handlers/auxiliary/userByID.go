package auxiliary

import (
	"IAM/pkg/logs"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func GetUserIDByEmail(ctx context.Context, email string, rdb *redis.Client) (string, error) {
	userID, err := rdb.Get(ctx, "email:"+email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("email %s not found", email)
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}
	if userID != "" {
		return userID, nil
	} else {
		return "", fmt.Errorf("email %s not found", email)
	}
}
