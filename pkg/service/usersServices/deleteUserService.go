package usersServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"fmt"
)

type UserManagementRepository interface {
	DeleteUserFromDB(email string) error
	GetAllUsersDataFromDB() ([]map[string]string, error)
	GetUsersListFromDB() ([]string, error)
	GetUserRole(email string) (string, error)
}

var UserManageRepo UserManagementRepository

func DeleteUserService(inputUserData models.EmailData) (string, error) {
	// deleting userdata from DB
	if err := UserManageRepo.DeleteUserFromDB(inputUserData.Email); err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}
	msg := fmt.Sprintf("user %v deleted successfully", inputUserData.Email)
	logs.AuditLogger.Info(msg)
	return msg, nil
}
