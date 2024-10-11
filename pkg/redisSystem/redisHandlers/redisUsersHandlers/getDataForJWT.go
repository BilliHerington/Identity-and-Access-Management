package redisUsersHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
)

func (repo *RedisUsersRepository) GetDataForJWT(email string) (userID string, userVersion string, error error) {

	// get userID
	userID, err := redisAuxiliaryHandlers.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return "", "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return "", "", err
	}

	// get userVersion
	userVersion, err = redisAuxiliaryHandlers.GetUserVersion(repo.RDB, userID)
	if err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return "", "", err
	}

	return userID, userVersion, nil
}
