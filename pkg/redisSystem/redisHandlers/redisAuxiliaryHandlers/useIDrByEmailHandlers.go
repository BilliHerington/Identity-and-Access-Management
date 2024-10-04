package redisAuxiliaryHandlers

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

type RedisUserIDByEmailRepo struct {
	RDB *redis.Client
}

func (repo *RedisUserIDByEmailRepo) GetUserIDByEmail(email string) (string, error) {
	userID, err := repo.RDB.Get(ctx, "email:"+email).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return userID, nil
}
