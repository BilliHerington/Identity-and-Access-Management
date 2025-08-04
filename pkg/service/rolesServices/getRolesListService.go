package rolesServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
)

func GetRolesListService() ([]string, error) {
	rolesList, err := RoleManageRepo.GetRolesListFromDB()
	if err != nil {
		if errors.Is(err, models.ErrRolesListEmpty) {
			return rolesList, err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return rolesList, err
	}
	//msg := fmt.Sprintf("roles list: %v", roles)
	return rolesList, nil
}
