package redisUsersHandlers

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisGetPasswordRepo struct {
	RDB *redis.Client
}

func (repo *RedisGetPasswordRepo) GetPassword(userID string) (string, error) {
	pass, err := repo.RDB.HGet(ctx, "user:"+userID, "password").Result()
	if err != nil {
		return "", err
	}
	return pass, nil
}
