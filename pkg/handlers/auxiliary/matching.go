package auxiliary

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

// RoleMatch check role exist in redis, return true if exist
func RoleMatch(roleKey string, rdb *redis.Client) (bool, error) {
	ctx := context.Background()
	res, err := rdb.HGetAll(ctx, roleKey).Result()
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, nil
	}
	return true, nil
}

// EmailMatch check email exist in redis, return true if exist
func EmailMatch(email string, rdb *redis.Client) (bool, error) {
	ctx := context.Background()

	// check email exist in redis
	emailKey := "email:" + email
	_, err := rdb.Get(ctx, emailKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}
