package redisRolesHandlers

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisAssignRoleRepo struct {
	RDB *redis.Client
}

func (repo *RedisAssignRoleRepo) AssignRoleToUser(userID, role string) error {
	err := repo.RDB.HSet(ctx, "user:"+userID, "role", role).Err()
	if err != nil {
		return err
	}
	return nil
}
