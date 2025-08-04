package redisAuthentication

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisAuthManagementRepository struct {
	RDB *redis.Client
}

func (repo *RedisAuthManagementRepository) GetPassword(email string) (string, error) {

	// get userID
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return "", err
	}

	// get password from redis
	redisPassword, err := repo.RDB.HGet(ctx, "user:"+userID, "password").Result()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return "", err
	}
	return redisPassword, nil
}
