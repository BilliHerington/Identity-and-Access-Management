package redisUsers

import (
	"IAM/pkg/logs"
	"IAM/pkg/repository/requests/redisInternal"
)

func (repo *RedisUserManagementRepository) GetUserRole(email string) (string, error) {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return "", err
	}
	role, err := repo.RDB.HGet(ctx, "user:"+userID, "role").Result()
	if err != nil {
		logs.Error.Println("failed get user role", err)
		logs.ErrorLogger.Error("failed get user role", err)
		return "", err
	}
	return role, nil
}
