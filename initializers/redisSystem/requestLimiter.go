package redisDB

import (
	"IAM/pkg/logs"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func RateLimitExceeded(key string, limit int, window int64, rdb *redis.Client) (bool, error) {
	ctx := context.Background()
	// using INCR for increasing value by key
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Printf("redisSystem Incr failed: %v", err)
		return false, err
	}
	// if 1st request, setting TTL (key lifetime)
	if count == 1 {
		err = rdb.Expire(ctx, key, time.Duration(window)*time.Second).Err()
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
