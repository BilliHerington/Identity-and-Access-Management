package redisAuxiliaryHandlers

import "github.com/go-redis/redis/v8"

type RedisRegistrationRepo struct {
	RDB *redis.Client
}

func (repo RedisRegistrationRepo) RegisterUser(userID, email, password, name, role, jwt, userVersion string) error {
	err := repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
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
	}, "user:"+userID)
	return err
}
