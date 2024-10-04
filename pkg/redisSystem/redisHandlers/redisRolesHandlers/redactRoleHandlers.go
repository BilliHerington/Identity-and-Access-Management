package redisRolesHandlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

type RedisRedactRoleRepo struct {
	RDB *redis.Client
}

func (repo RedisRedactRoleRepo) RedactRoleDB(roleName string, privileges []string) error {
	marshalPrivileges, err := json.Marshal(privileges)
	if err != nil {
		return err
	}
	err = repo.RDB.HSet(ctx, "role:"+roleName, "privileges", marshalPrivileges).Err()
	if err != nil {
		return err
	}
	return nil
}
