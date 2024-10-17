package redisInternal

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisAuxiliaryRepository struct {
	RDB *redis.Client
}

//// GetRole return role by key from redis
//func (repo *RedisAuxiliaryRepository) GetRole(role string) (string, error) {
//	resultRole, err := repo.RDB.HGet(ctx, "role:"+role, "name").Result()
//	if errors.Is(err, redis.Nil) {
//		return "", nil
//	}
//	return resultRole, err
//}
//
//func (repo *RedisAuxiliaryRepository) GetEmail(email string) (string, error) {
//	resultEmail, err := repo.RDB.Get(ctx, "email:"+email).Result()
//	if errors.Is(err, redis.Nil) {
//		return "", nil
//	}
//	return resultEmail, err
//}
//func (repo *RedisAuxiliaryRepository) GetUserIDByEmail(email string) (string, error) {
//	userID, err := repo.RDB.Get(ctx, "email:"+email).Result()
//	if errors.Is(err, redis.Nil) {
//		return "", errors.New("email not found")
//	}
//	return userID, err
//}
//func (repo *RedisAuxiliaryRepository) GetUserVersion(userID string) (string, error) {
//	userVersion, err := repo.RDB.HGet(ctx, "user:"+userID, "userVersion").Result()
//	return userVersion, err
//}

func CheckRoleExist(rdb *redis.Client, role string) (bool, error) {
	isMember, err := rdb.SIsMember(ctx, "redisRoles", role).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}
	return isMember, nil
}

func CheckEmailExist(rdb *redis.Client, email string) (bool, error) {
	err := rdb.Get(ctx, "email:"+email).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("cannot get email from redis: %w", err)
	}
	return true, nil
}

func GetUserIDByEmail(rdb *redis.Client, email string) (string, error) {
	userID, err := rdb.HGet(ctx, "email:"+email, "id").Result()
	if errors.Is(redis.Nil, err) {
		return "", errors.New("user does not exist")
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
