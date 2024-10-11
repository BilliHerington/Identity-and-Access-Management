package redisUsersHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"context"
)

var ctx = context.Background()

func (repo *RedisUsersRepository) GetPassword(email string) (string, error) {

	// get userID
	userID, err := redisAuxiliaryHandlers.GetUserIDByEmail(repo.RDB, email)
	if err.Error() == "user does not exist" {
		return "", err
	} else if err != nil {
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
