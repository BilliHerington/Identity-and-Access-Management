package redisRolesHandlers

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

type RedisGetAllRolesDataRepo struct {
	RDB *redis.Client
}

func (repo RedisGetAllRolesDataRepo) GetAllRolesDataFromDB() ([]map[string]string, error) {
	var roles []map[string]string

	// get roles list from redis
	allRoles, err := repo.RDB.SMembers(ctx, "roles").Result()
	if errors.Is(err, redis.Nil) {
		return roles, errors.New("no roles found")
	}
	if err != nil {
		return roles, err
	}

	// get all data by role
	for _, roleName := range allRoles {
		roleData, err := repo.RDB.HGetAll(ctx, "role:"+roleName).Result()
		if err != nil {
			return roles, err
		}
		if len(roleData) > 0 {
			roles = append(roles, roleData)
		} else {
			roles = append(roles, map[string]string{
				"error": "No data found",
			})
		}
	}
	return roles, nil
}
