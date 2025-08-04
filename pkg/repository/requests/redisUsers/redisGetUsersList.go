package redisUsers

import (
	"IAM/pkg/models"
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisUserManagementRepository) GetUsersListFromDB() ([]string, error) {
	user, err := repo.RDB.SMembers(ctx, "users").Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, models.ErrUserDoesNotExist
	} else if err != nil {
		return []string{}, err
	}
	return user, nil
}
