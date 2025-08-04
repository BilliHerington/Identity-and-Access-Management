package redisRoles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"encoding/json"
)

func (repo *RedisRolesManagementRepository) RedactRoleDB(roleName string, privileges []string) error {
	roleExist, err := redisInternal.CheckRoleExist(repo.RDB, roleName)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if !roleExist {
		return models.ErrRoleDoesNotExist
	}
	marshalPrivileges, err := json.Marshal(privileges)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	err = repo.RDB.HSet(ctx, "role:"+roleName, "privileges", marshalPrivileges).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	return nil
}
