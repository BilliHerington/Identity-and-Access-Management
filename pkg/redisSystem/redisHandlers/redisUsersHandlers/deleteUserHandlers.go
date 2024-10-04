package redisUsersHandlers

import (
	"IAM/pkg/logs"
	"github.com/go-redis/redis/v8"
)

type RedisDeleteUserRepo struct {
	RDB *redis.Client
}

func (repo *RedisDeleteUserRepo) DeleteUserDB(userID, email string) error {

	// delete user data from redis
	err := repo.RDB.Del(ctx, "user:"+userID).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return err
	}

	// delete email from redis
	err = repo.RDB.Del(ctx, "email:"+email).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return err
	}

	// delete userID from list in redis
	err = repo.RDB.SRem(ctx, "users", userID).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return err
	}
	return nil
}
