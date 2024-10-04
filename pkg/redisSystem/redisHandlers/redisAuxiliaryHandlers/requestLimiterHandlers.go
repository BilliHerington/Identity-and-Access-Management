package redisAuxiliaryHandlers

import (
	"IAM/pkg/logs"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRequestRepo struct {
	RDB *redis.Client
}

func (repo *RedisRequestRepo) GetRequestLimit(key string, limit int, window int64) (bool, error) {

	// using INCR for increasing value by key
	count, err := repo.RDB.Incr(ctx, key).Result()
	if err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Printf("redisSystem Incr failed: %v", err)
		return false, err
	}

	// if 1st request, setting TTL (key lifetime)
	if count == 1 {
		err = repo.RDB.Expire(ctx, key, time.Duration(window)*time.Second).Err()
		if err != nil {
			logs.ErrorLogger.Error(err.Error())
			logs.Error.Printf("redisSystem Expire failed: %v", err)
			return false, err
		}
	}

	// if user exceed request limit
	if count > int64(limit) {
		return true, nil
	}

	return false, nil
}
