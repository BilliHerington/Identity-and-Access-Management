package redisInternal

import (
	"IAM/pkg/models"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisAuxiliaryRepository struct {
	RDB *redis.Client
}

func CheckRoleExist(rdb *redis.Client, role string) (bool, error) {
	isMember, err := rdb.SIsMember(ctx, "roles", role).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}
	return isMember, nil
}

func CheckEmailExist(rdb *redis.Client, email string) (bool, error) {
	err := rdb.Get(ctx, "email:"+email).Err()
	//logs.Info.Printf("error from redis:%s", err)
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot get email from redis: %w", err)
	}
	return true, nil
}

func GetUserIDByEmail(rdb *redis.Client, email string) (string, error) {
	userID, err := rdb.Get(ctx, "email:"+email).Result()
	if errors.Is(redis.Nil, err) {
		return "", models.ErrUserDoesNotExist
	} else if err != nil {
		return "", fmt.Errorf("cannot get user id from redis: %w", err)
	}
	return userID, nil
}

func GetUserVersion(rdb *redis.Client, userID string) (string, error) {
	userVersion, err := rdb.HGet(ctx, "user:"+userID, "userVersion").Result()
	if err != nil {
		return "", fmt.Errorf("cannot get user version from redis: %w", err)
	}
	return userVersion, nil
}
