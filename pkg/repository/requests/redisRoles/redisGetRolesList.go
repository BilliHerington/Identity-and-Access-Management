package redisRoles

import (
	"IAM/pkg/models"
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisRolesManagementRepository) GetRolesListFromDB() ([]string, error) {
	roles, err := repo.RDB.SMembers(ctx, "roles").Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, models.ErrRolesListEmpty
	} else if err != nil {
		return []string{}, err
	}
	return roles, nil
}
