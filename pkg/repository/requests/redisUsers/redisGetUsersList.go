package redisUsers

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisUserManagementRepository) GetUsersListFromDB() ([]string, error) {
	user, err := repo.RDB.SMembers(ctx, "users").Result()
	if errors.Is(err, redis.Nil) {
		return []string{}, errors.New("users not found")
	} else if err != nil {
		return []string{}, err
	}
	return user, nil
}
