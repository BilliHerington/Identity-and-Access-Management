package redisAuxiliaryHandlers

import (
	"github.com/go-redis/redis/v8"
)

func (repo *RedisAuxiliaryRepository) StartUserRegistration(userID, email, verificationCode string) error {
	userKey := "user:" + userID
	err := repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, userKey, map[string]interface{}{
				"id":               userID,
				"email":            email,
				"verificationCode": verificationCode,
			})
			pipe.SAdd(ctx, "users", userID)
			pipe.Set(ctx, "email:"+email, userID, 0)
			return nil
		})
		return err
	}, userKey)
	return err
}

func (repo *RedisAuxiliaryRepository) SaveUser(userID, email, password, name, role, jwt, userVersion string) error {
	userKey := "user:" + userID
	err := repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, userKey, map[string]interface{}{
				"id":          userID,
				"email":       email,
				"name":        name,
				"password":    password,
				"role":        role,
				"jwt":         jwt,
				"userVersion": userVersion,
			})
			pipe.SAdd(ctx, "users", userID)
			pipe.Set(ctx, "email:"+email, userID, 0)
			return nil
		})
		return err
	}, userKey)
	return err
}
func (repo *RedisAuxiliaryRepository) GetVerificationCode(userID string) (string, error) {
	userKey := "user:" + userID
	code, err := repo.RDB.HGet(ctx, userKey, "verificationCode").Result()
	if err != nil {
		return "", err
	}
	err = repo.RDB.HDel(ctx, userKey, "verificationCode").Err()
	return code, err
}
