package redisDB

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"reflect"
	"strconv"
)

//var (
//	Rdb *redis.Client
//	Ctx = context.Background()
//)

func InitRedis() (*redis.Client, error) {
	initializers.LoadEnvVariables()
	ctx := context.Background()
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logs.Error.Fatalf("invalid REDIS_DB value %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Fatalf("redisSystem ping failed: %v", err)
	}
	InitializeRoles(rdb)
	InitializeAdmin(rdb)
	return rdb, nil
}

func InitializeRoles(rdb *redis.Client) {
	var (
		adminMatch bool
		userMatch  bool
	)

	adminRoleKey := "role:admin"
	userRoleKey := "role:user"

	ctx := context.Background()

	// check role:admin exist
	adminRes, err := rdb.HGetAll(ctx, adminRoleKey).Result()
	if err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Fatalf("redisSystem HGetAll failed: %v", err)
	} else if len(adminRes) == 0 {
		adminMatch = false
	} else {
		adminMatch = true
	}
	// check roel:user exist
	userRes, err := rdb.HGetAll(ctx, userRoleKey).Result()
	if err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Fatalf("redisSystem HGetAll failed: %v", err)
	} else if len(userRes) == 0 {
		userMatch = false
	} else {
		userMatch = true
	}

	// create admin if not exist
	if !adminMatch {
		val := reflect.ValueOf(models.AdminPrivileges)
		var adminPrivileges []string
		for i := 0; i < val.NumField(); i++ {
			adminPrivileges = append(adminPrivileges, val.Field(i).String())
		}
		marshaledAdminPrivileges, err := json.Marshal(adminPrivileges)
		if err != nil {
			logs.ErrorLogger.Errorf("marshal admin privileges failed: %v", err)
			logs.Error.Fatalf("marshal admin privileges failed: %v", err)
		}
		err = rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, adminRoleKey, map[string]interface{}{
					"name":       "admin",
					"privileges": marshaledAdminPrivileges,
				})
				pipe.SAdd(ctx, "roles", "admin")
				return nil
			})
			return err
		}, adminRoleKey)
		logs.AuditLogger.Println("role 'admin' created successfully")
		logs.Info.Println("role 'admin' created successfully")
	}
	// create user if not exist
	if !userMatch {
		val := reflect.ValueOf(models.UserPrivileges)
		var userPrivileges []string
		for i := 0; i < val.NumField(); i++ {
			userPrivileges = append(userPrivileges, val.Field(i).String())
		}

		marshaledUserPrivileges, err := json.Marshal(userPrivileges)
		if err != nil {
			logs.ErrorLogger.Errorf("marshal user privileges failed: %v", err)
			logs.Error.Fatalf("marshal user privileges failed: %v", err)
		}
		err = rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, userRoleKey, map[string]interface{}{
					"name":       "user",
					"privileges": marshaledUserPrivileges,
				})
				pipe.SAdd(ctx, "roles", "user")
				return nil
			})
			return err
		}, userRoleKey)
		logs.AuditLogger.Println("role 'user' created successfully")
		logs.Info.Println("role 'user' created successfully")
	}
}
func InitializeAdmin(rdb *redis.Client) {
	ctx := context.Background()

	userID := "ROOT"
	adminMail := os.Getenv("ROOT_EMAIL")
	adminPassword := os.Getenv("ROOT_PASSWORD")
	name := "ROOT"
	res, err := rdb.HGetAll(ctx, "user:"+userID).Result()
	if err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Fatalf("redisSystem HGetAll failed: %v", err)
	}
	if len(res) == 0 {
		// hashing pass
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			logs.ErrorLogger.Error(err.Error())
			logs.Error.Fatal(err)
		}
		adminPassword = string(hashedPassword)

		// creating userVersion
		userVersion := uuid.New().String()

		err = rdb.Watch(ctx, func(tx *redis.Tx) error {
			_, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.HMSet(ctx, "user:"+userID, map[string]interface{}{
					"id":          userID,
					"email":       adminMail,
					"name":        name,
					"password":    adminPassword,
					"role":        "admin",
					"jwt":         "",
					"userVersion": userVersion,
				})
				pipe.SAdd(ctx, "users", userID)
				return nil
			})
			return err
		}, "user:"+userID)
		if err != nil {
			logs.ErrorLogger.Errorf("admin creataion failed: %v", err)
			logs.Error.Fatalf("admin creataion failed: %v", err)
		}
		// add new EmailKey for User
		err = rdb.Set(ctx, "email:"+adminMail, userID, 0).Err()
		if err != nil {
			logs.Error.Println(err)
			logs.ErrorLogger.Error(err.Error())
		}
		logs.AuditLogger.Printf("ROOT created successfully: %s", adminMail)
		logs.Info.Printf("ROOT created successfully: %s", adminMail)
	}
}
