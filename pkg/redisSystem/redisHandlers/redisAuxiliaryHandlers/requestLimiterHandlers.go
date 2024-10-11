package redisAuxiliaryHandlers

import (
	"IAM/pkg/logs"
	"fmt"
	"time"
)

func (repo *RedisAuxiliaryRepository) GetRequestLimit(key string, limit int, window int64) (bool, error) {

	// using INCR for increasing value by key
	count, err := repo.RDB.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis Incr failed: %v", err)
	}

	// if 1st request, setting TTL (key lifetime)
	if count == 1 {
		err = repo.RDB.Expire(ctx, key, time.Duration(window)*time.Second).Err()
		if err != nil {
			logs.ErrorLogger.Error(err.Error())
			logs.Error.Printf("redis Expire failed: %v", err)
			return false, fmt.Errorf("redis Expire failed: %v", err)
		}
	}

	// if user exceed request limit
	if count > int64(limit) {
		return true, nil
	}

	return false, nil
}
