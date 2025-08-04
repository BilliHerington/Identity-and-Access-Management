package redisRoles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"encoding/json"
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
	roleExist, err := redisInternal.CheckRoleExist(repo.RDB, roleName)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if roleExist {
		return models.ErrRoleAlreadyExist
	}

	logs.Info.Printf("name:%s and priv:%s", roleName, privilegesMarshaled)
	// writing in redis
	//var someString []string
	//someString = []string{"penis", "penis", "penis"}

	err = repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, "role:"+roleName, "privileges", string(privilegesMarshaled))
			pipe.SAdd(ctx, "roles", roleName)
			return nil
		})
		if err != nil {
			logs.ErrorLogger.Error(err)
			logs.Error.Println(err)
		}
		return err
	}, "role:"+roleName)
	return err
}
