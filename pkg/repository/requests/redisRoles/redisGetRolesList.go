package redisRoles

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisRolesManagementRepository) GetRolesListFromDB() ([]string, error) {
	roles, err := repo.RDB.SMembers(ctx, "roles").Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, errors.New("roles not found")
	} else if err != nil {
		return []string{}, err
	}
	return roles, nil
}
