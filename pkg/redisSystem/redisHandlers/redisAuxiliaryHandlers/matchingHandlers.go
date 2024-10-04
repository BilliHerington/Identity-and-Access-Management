package redisAuxiliaryHandlers

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisRoleRepo realization RoleRepository interface for redis
type RedisRoleRepo struct {
	RDB *redis.Client
}

// GetRole return role by key from redis
func (repo *RedisRoleRepo) GetRole(role string) (string, error) {
	resultRole, err := repo.RDB.HGet(ctx, "role:"+role, "name").Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return resultRole, nil
}

type RedisEmailRepo struct {
	RDB *redis.Client
}

func (repo *RedisEmailRepo) GetEmail(email string) (string, error) {
	resultEmail, err := repo.RDB.HGet(ctx, "email:"+email, "name").Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return resultEmail, nil
}
