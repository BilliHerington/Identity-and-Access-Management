package redisAuxiliaryHandlers

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

type RedisUserVersionRepo struct {
	RDB *redis.Client
}

func (repo *RedisUserVersionRepo) GetUserVersion(userID string) (string, error) {
	userVersion, err := repo.RDB.HGet(ctx, "user:"+userID, "userVersion").Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("user version not found")
	} else if err != nil {
		return "", err
	}
	return userVersion, nil
}
