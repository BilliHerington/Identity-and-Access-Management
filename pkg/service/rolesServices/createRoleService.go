package rolesServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"fmt"
)

func CreateRoleService(inputUserData models.RolesData) (string, error) {
	// save role
	if err := RoleManageRepo.CreateRole(inputUserData.RoleName, inputUserData.Privileges); err != nil {
		if errors.Is(err, models.ErrRoleAlreadyExist) {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}
	msg := fmt.Sprintf("Role %s created successfully", inputUserData.RoleName)
	logs.AuditLogger.Println(msg)
	return msg, nil
}
