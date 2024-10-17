package redisInternal

import (
	"IAM/pkg/logs"
	"github.com/go-redis/redis/v8"
)

func SaveUserInRedis(rdb *redis.Client, userID, email, password, name, role, jwt, userVersion string) error {
	userKey := "user:" + userID
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := pipe.HMSet(ctx, userKey, map[string]interface{}{
				"id":          userID,
				"email":       email,
				"name":        name,
				"password":    password,
				"role":        role,
				"jwt":         jwt,
				"userVersion": userVersion,
			}).Err(); err != nil {
				logs.ErrorLogger.Error("failed set fields in userKey", err)
				logs.Error.Println("failed set fields in userKey", err)
				return err
			}
			if err := pipe.SAdd(ctx, "redisUsers", userID).Err(); err != nil {
				logs.ErrorLogger.Error("failed add userID in redisUsers", err)
				logs.Error.Println("failed add userID in redisUsers", err)
				return err
			}
			if err := pipe.Set(ctx, "email:"+email, userID, 0).Err(); err != nil {
				logs.ErrorLogger.Error("failed add email key", err)
				logs.Error.Println("failed add email key", err)
				return err
			}
			return nil
		})
		if err != nil {
			logs.ErrorLogger.Error("failed watch user", err)
			logs.Error.Println("failed watch user", err)
			return err
		}
		return nil
	}, userKey)
	return err
}
