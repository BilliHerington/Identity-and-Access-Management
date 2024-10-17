package redisRoles

import (
	"IAM/pkg/logs"
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisRolesManagementRepository) GetAllRolesDataFromDB() ([]map[string]string, error) {
	var roles []map[string]string

	// get redisRoles list from redis
	allRoles, err := repo.RDB.SMembers(ctx, "roles").Result()
	if errors.Is(err, redis.Nil) {
		return roles, errors.New("no redisRoles found")
	}
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return roles, err
	}

	// get all data by role
	for _, roleName := range allRoles {
		roleData, err := repo.RDB.HGetAll(ctx, "role:"+roleName).Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err)
			return roles, err
		}
		if len(roleData) > 0 {
			roles = append(roles, roleData)
		} else {
			roles = append(roles, map[string]string{
				"error": "No data found for role: " + roleName,
			})
		}
	}
	return roles, nil
}
