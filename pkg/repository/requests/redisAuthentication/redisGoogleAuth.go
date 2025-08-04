package redisAuthentication

import (
	"IAM/pkg/logs"
	"IAM/pkg/repository/requests/redisInternal"
)

func (repo *RedisAuthManagementRepository) SaveUserByGoogle(userID, email, password, name, role, jwt, userVersion string) error {
	if err := redisInternal.SaveUserInRedis(repo.RDB, userID, email, password, name, role, jwt, userVersion); err != nil {
		logs.Error.Println("failed save user in redis", err)
		logs.ErrorLogger.Error("failed save user in redis", err)
		return err
	}
	return nil
}
func (repo *RedisAuthManagementRepository) CheckEmailExist(email string) (bool, error) {
	//logs.Info.Printf("trying find: %s", email)
	emailExist, err := redisInternal.CheckEmailExist(repo.RDB, email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return false, err
	}
	if emailExist {
		return true, nil
	} else {
		return false, nil
	}
}
