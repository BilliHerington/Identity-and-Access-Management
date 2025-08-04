package rolesServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
	"fmt"
)

func RedactRoleService(inputUserData models.RolesData) (string, error) {

	// redact role
	if err := RoleManageRepo.RedactRoleDB(inputUserData.RoleName, inputUserData.Privileges); err != nil {
		if errors.Is(err, models.ErrRoleDoesNotExist) {
			return "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return "", err
	}
	msg := fmt.Sprintf("%s updated successfully. New privileges: %s", inputUserData.RoleName, inputUserData.Privileges)
	logs.AuditLogger.Printf(msg)
	return msg, nil
}
