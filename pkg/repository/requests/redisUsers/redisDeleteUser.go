package redisUsers

import (
	"IAM/pkg/logs"
	"IAM/pkg/repository/requests/redisInternal"
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisUserManagementRepository struct {
	RDB *redis.Client
}

var ctx = context.Background()

func (repo *RedisUserManagementRepository) DeleteUserFromDB(email string) error {

	// get userID
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
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
			if err = pipe.Del(ctx, "user:"+userID).Err(); err != nil {
				logs.Error.Println("Failed to delete user:", err)
				logs.ErrorLogger.Errorf("failed to delete user: %s", err)
				return err
			}
			if err = pipe.Del(ctx, "email:"+email).Err(); err != nil {
				logs.ErrorLogger.Errorf("failed to delete user: %s", err)
				logs.Error.Println("Failed to delete email:", err)
				return err
			}
			if err = pipe.SRem(ctx, "redisUsers", userID).Err(); err != nil {
				logs.Error.Println("Failed to remove user from set:", err)
				logs.ErrorLogger.Errorf("failed to delete user: %s", err)
				return err
			}
			return nil
		})
		if err != nil {
			logs.Error.Println("Redis Tx failed:", err)
			logs.ErrorLogger.Errorf("Redis Tx failed: %s", err)
			return err
		}
		return nil
	})

	if err != nil {
		logs.Error.Println("Redis Watch failed: ", err)
		logs.ErrorLogger.Errorf("Redis Watch failed: %s", err)
		return err
	}
	return nil
}
