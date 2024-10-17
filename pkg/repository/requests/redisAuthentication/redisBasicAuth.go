package redisAuthentication

import (
	"IAM/pkg/logs"
	"IAM/pkg/repository/requests/redisInternal"
	"context"
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
		if err.Error() == "user does not exist" {
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
