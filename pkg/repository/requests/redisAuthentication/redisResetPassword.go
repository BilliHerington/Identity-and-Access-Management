package redisAuthentication

import (
	"IAM/pkg/logs"
	"IAM/pkg/repository/requests/redisInternal"
)

func (repo *RedisAuthManagementRepository) SavePassCode(email, passCode string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return err
		}
		logs.Error.Println(err)
		logs.Error.Println(email)
		return err
	}
	if err = repo.RDB.HSet(ctx, "user:"+userID, passCode, 0).Err(); err != nil {
		logs.Error.Println("failed save pass code", err)
		logs.ErrorLogger.Error("failed save pass code", err)
		return err
	}
	return err
}
func (repo *RedisAuthManagementRepository) SaveNewUserData(email, password, userVersion string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return err
		}
		logs.Error.Println(err)
		logs.Error.Println(email)
		return err
	}
	err = repo.RDB.HMSet(ctx, "user:"+userID, map[string]interface{}{
		"password":    password,
		"userVersion": userVersion,
	}, 0).Err()
	if err != nil {
		logs.Error.Println("failed save user data", err)
		logs.ErrorLogger.Error("failed save user data", err)
		return err
	}
	return err
}
