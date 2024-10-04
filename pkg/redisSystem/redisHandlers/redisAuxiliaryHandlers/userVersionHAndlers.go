package redisAuxiliaryHandlers

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisUserVersionRepo struct {
	RDB *redis.Client
}

func (repo *RedisUserVersionRepo) GetUserVersion(userID string) (string, error) {
	userVersion, err := repo.RDB.HGet(context.Background(), "user:"+userID, "userVersion").Result()
	if err != nil {
		return "", err
	}
	return userVersion, nil
}
