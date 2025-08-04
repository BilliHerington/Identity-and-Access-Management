package usersServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
)

func GetUsersListService() ([]string, error) {
	// get all Users from user-list in redis
	usersList, err := UserManageRepo.GetUsersListFromDB()
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return usersList, err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return usersList, err
	}

	return usersList, nil
}
