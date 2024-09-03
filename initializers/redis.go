package initializers

import (
	"IAM/pkg/logs"
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	LoadEnvVariables()

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logs.Error.Fatalf("invalid REDIS_DB value %v", err)
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err = Rdb.Ping(Ctx).Result()
	if err != nil {
		logs.Error.Fatalf("redis ping failed: %v", err)
	}
}
