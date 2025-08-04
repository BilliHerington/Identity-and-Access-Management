package redisAuthentication

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"errors"
)

func (repo *RedisAuthManagementRepository) SavePassCode(email, passCode string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.Error.Println(err)
		logs.Error.Println(email)
		return err
	}
	if err = repo.RDB.HSet(ctx, "user:"+userID, "verificationCode", passCode).Err(); err != nil {
		logs.Error.Println("failed save pass code", err)
		logs.ErrorLogger.Error("failed save pass code", err)
		return err
	}
	return err
}
func (repo *RedisAuthManagementRepository) SaveNewUserData(email, password, userVersion string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	//	logs.Info.Printf("saving user:%s\nuservers:%s", email, userVersion)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.Error.Println(err)
		logs.Error.Println(email)
		return err
	}
	//logs.Info.Printf("userdata:\nemail:%s\npass:%s\nuserID:%s\nuserVersion:%s\n", email, password, userID, userVersion)
	err = repo.RDB.HSet(ctx, "user:"+userID, "password", password, "userVersion", userVersion).Err()
	if err != nil {
		logs.Error.Println("failed save user data", err)
		logs.ErrorLogger.Error("failed save user data", err)
		return err
	}
	return err
}
