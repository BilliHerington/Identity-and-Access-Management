package rolesServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"fmt"
)

func DeleteRoleService(inputUserData models.RolesData) (string, error) {

	// delete role from db
	if err := RoleManageRepo.DeleteRole(inputUserData.RoleName); err != nil {
		if errors.Is(err, models.ErrRoleDoesNotExist) {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}
	msg := fmt.Sprintf("Role %s has been deleted", inputUserData.RoleName)
	logs.AuditLogger.Println(msg)
	return msg, nil
}
