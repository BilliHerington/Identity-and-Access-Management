package redisRolesHandlers

import (
	"github.com/go-redis/redis/v8"
)

type GetRolesListRepo struct {
	RDB *redis.Client
}

func (repo *GetRolesListRepo) GetRolesListFromDB() ([]string, error) {
	roles, err := repo.RDB.SMembers(ctx, "roles").Result()
	if err != nil {
		return roles, err
	}
	return roles, nil
}
