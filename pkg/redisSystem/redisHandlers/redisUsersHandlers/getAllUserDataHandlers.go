package redisUsersHandlers

import (
	"IAM/pkg/logs"
	"github.com/go-redis/redis/v8"
)

type RedisGetAllUsersDataRepo struct {
	RDB *redis.Client
}

func (repo *RedisGetAllUsersDataRepo) GetAllUsersDataFromDB() ([]map[string]string, error) {
	var users []map[string]string

	// get all ID from user-list in redis
	listID, err := repo.RDB.SMembers(ctx, "users").Result()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err.Error())
		return users, err
	}

	// get all data for every ID
	for _, id := range listID {
		userData, err := repo.RDB.HGetAll(ctx, "user:"+id).Result()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
			return users, err
		}
		if len(userData) > 0 {
			users = append(users, userData)
		} else {
			users = append(users, map[string]string{
				"error": "No data found",
			})
		}
	}
	return users, nil
}
