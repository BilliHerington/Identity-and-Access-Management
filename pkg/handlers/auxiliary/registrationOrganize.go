package auxiliary

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func RegistrationOrganizeHandler(rdb *redis.Client, ctx context.Context, userID, email, password, name, role, jwt, userVersion string) error {
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
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
			return nil
		})
		return err
	}, "user:"+userID)
	return err
}
