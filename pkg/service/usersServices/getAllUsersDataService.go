package usersServices

import (
	"IAM/pkg/logs"
)

func GetAllUsersDataService() ([]map[string]string, error) {
	allUsersData, err := UserManageRepo.GetAllUsersDataFromDB()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return allUsersData, err
	}

	return allUsersData, nil
}
