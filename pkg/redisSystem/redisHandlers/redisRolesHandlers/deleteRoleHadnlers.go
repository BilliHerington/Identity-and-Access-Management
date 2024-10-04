package redisRolesHandlers

import (
	"github.com/go-redis/redis/v8"
)

type RedisDeleteRoleRepo struct {
	RDB *redis.Client
}

func (repo *RedisDeleteRoleRepo) DeleteRole(roleName string) error {

	// deleting role from role-list in redis
	err := repo.RDB.SRem(ctx, "roles", roleName).Err()
	if err != nil {
		return err
	}

	// deleting role data from redis
	err = repo.RDB.Del(ctx, "role:"+roleName).Err()
	if err != nil {
		return err
	}
	return nil
}
