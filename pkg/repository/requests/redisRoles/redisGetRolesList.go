package redisRoles

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisRolesManagementRepository) GetRolesListFromDB() ([]string, error) {
	roles, err := repo.RDB.SMembers(ctx, "redisRoles").Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, errors.New("redisRoles not found")
	} else if err != nil {
		return []string{}, err
	}
	return roles, nil
}
