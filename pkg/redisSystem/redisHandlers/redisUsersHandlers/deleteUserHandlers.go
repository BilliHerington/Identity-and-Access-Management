package redisUsersHandlers

import (
	"IAM/pkg/logs"
	"IAM/pkg/redisSystem/redisHandlers/redisAuxiliaryHandlers"
	"github.com/go-redis/redis/v8"
)

type RedisUsersRepository struct {
	RDB *redis.Client
}

func (repo *RedisUsersRepository) DeleteUserFromDB(email string) error {

	// get userID
	userID, err := redisAuxiliaryHandlers.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}

	err = repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := pipe.Del(ctx, "user:"+userID).Err(); err != nil {
				logs.Error.Println("Failed to delete user:", err)
				logs.ErrorLogger.Errorf("failed to delete user: %s", err)
				return err
			}
			if err := pipe.Del(ctx, "email:"+email).Err(); err != nil {
				logs.ErrorLogger.Errorf("failed to delete user: %s", err)
				logs.Error.Println("Failed to delete email:", err)
				return err
			}
			if err := pipe.SRem(ctx, "users", userID).Err(); err != nil {
				logs.Error.Println("Failed to remove user from set:", err)
				
				return err
			}
			return nil
		})
		return err
	})
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	return nil
}
