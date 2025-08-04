package rolesServices

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"errors"
)

func GetAllRolesDataService() ([]map[string]string, error) {
	roleData, err := RoleManageRepo.GetAllRolesDataFromDB()
	if err != nil {
		if errors.Is(err, models.ErrRolesListEmpty) {
			return roleData, err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return roleData, err
	}
	//msg := fmt.Sprintf("roles data: %s", roleData)
	return roleData, nil
}
