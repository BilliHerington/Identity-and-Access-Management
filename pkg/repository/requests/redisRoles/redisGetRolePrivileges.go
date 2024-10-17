package redisRoles

import (
	"IAM/pkg/logs"
	"encoding/json"
)

func (repo *RedisRolesManagementRepository) GetRolePrivileges(role string) ([]string, error) {
	var privileges []string
	rawPrivileges, err := repo.RDB.HGet(ctx, "role:"+role, "privileges").Result()
	if err != nil {
		logs.Error.Println("failed to get role privileges:", err)
		logs.ErrorLogger.Fatalln("failed to get role privileges:", err)
		return privileges, err
	}
	if err = json.Unmarshal([]byte(rawPrivileges), &privileges); err != nil {
		logs.Error.Println("failed to unmarshal role privileges:", err)
		logs.ErrorLogger.Fatalln("failed to unmarshal role privileges:", err)
		return privileges, err
	}
	return privileges, nil
}
