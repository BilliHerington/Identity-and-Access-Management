package redisRolesHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisRolesManagementRepository) DeleteRole(roleName string) error {
	roleExist, err := redisAuxiliaryHandlers.CheckRoleExist(repo.RDB, roleName)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if !roleExist {
		return errors.New("role does not exist")
	}
	err = repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.SRem(ctx, "roles", roleName) // deleting role from role-list in redis
			pipe.Del(ctx, "role:"+roleName)   // deleting role data from redis
			return nil
		})
		return err
	})
	return err
}
