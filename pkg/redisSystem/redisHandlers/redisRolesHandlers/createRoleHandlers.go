package redisRolesHandlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

type RedisCreateRoleRepo struct {
	RDB *redis.Client
}

func (repo RedisCreateRoleRepo) CreateRole(roleName string, privileges []string) error {

	// marshal Privileges for writing in redis
	privilegesMarshaled, err := json.Marshal(privileges)
	if err != nil {
		return err
	}

	// writing in redis
	err = repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, "role:"+roleName, map[string]interface{}{
				"name":       roleName,
				"privileges": privilegesMarshaled,
			})
			pipe.SAdd(ctx, "roles", roleName)
			return nil
		})
		return err
	}, "role:"+roleName)
	return err
}
