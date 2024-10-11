package redisRolesHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisRolesManagementRepository) CreateRole(roleName string, privileges []string) error {

	// marshal Privileges for writing in redis
	privilegesMarshaled, err := json.Marshal(privileges)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}

	// check role exist
	roleExist, err := redisAuxiliaryHandlers.CheckRoleExist(repo.RDB, roleName)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if roleExist {
		return errors.New("role already exist")
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
