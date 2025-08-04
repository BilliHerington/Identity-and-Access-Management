package rolesServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"fmt"
)

type RoleManagementRepository interface {
	AssignRoleToUser(email, role string) error
	CreateRole(roleName string, privileges []string) error
	DeleteRole(roleName string) error
	GetAllRolesDataFromDB() ([]map[string]string, error)
	GetRolesListFromDB() ([]string, error)
	RedactRoleDB(role string, privileges []string) error
	GetRolePrivileges(role string) ([]string, error)
}

var RoleManageRepo RoleManagementRepository

func AssignRoleService(inputUserData models.UserRoleData) (string, error) {
	// assign role
	if err := RoleManageRepo.AssignRoleToUser(inputUserData.Email, inputUserData.Role); err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) || errors.Is(err, models.ErrRoleDoesNotExist) {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}

	msg := fmt.Sprintf("role: %s, assign successfully for user: %s", inputUserData.Role, inputUserData.Email)
	logs.AuditLogger.Printf(msg)
	return msg, nil
}
