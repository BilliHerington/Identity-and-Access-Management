package redisDB

import (
	"IAM/initializers"
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"IAM/pkg/repository/requests/redisRoles"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"reflect"
	"strconv"
)

var ctx = context.Background()

func InitRedis() (*redis.Client, error) {
	initializers.LoadEnvVariables()
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
		logs.Error.Fatalf("redis ping failed: %v", err)
	}
	initFirstTime := os.Getenv("INIT_REDIS_IN_FIRST_TIME")
	if initFirstTime == "true" {
		InitializeRoles(rdb)
		InitializeAdmin(rdb)
	}
	return rdb, nil
}

func InitializeRoles(rdb *redis.Client) {

	var adminPrivileges []string
	val := reflect.ValueOf(models.AdminPrivileges)
	for i := 0; i < val.NumField(); i++ {
		adminPrivileges = append(adminPrivileges, val.Field(i).String())
	}

	var userPrivileges []string
	val = reflect.ValueOf(models.UserPrivileges)
	for i := 0; i < val.NumField(); i++ {
		userPrivileges = append(userPrivileges, val.Field(i).String())
	}

	// create roles
	repository := redisRoles.RedisRolesManagementRepository{RDB: rdb}
	if err := repository.CreateRole("admin", adminPrivileges); err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Fatal(err)
	}
	if err := repository.CreateRole("user", userPrivileges); err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Fatal(err)
	}

	logs.AuditLogger.Println("roles 'user' and 'admin' created successfully")
	logs.Info.Println("roles 'user' and 'admin' created successfully")

}
func InitializeAdmin(rdb *redis.Client) {

	userID := "ROOT"
	adminMail := os.Getenv("ROOT_EMAIL")
	adminPassword := os.Getenv("ROOT_PASSWORD")
	name := "ROOT"
	userVersion := uuid.New().String()

	// hashing pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Fatal(err)
	}
	adminPassword = string(hashedPassword)

	if err = redisInternal.SaveUserInRedis(rdb, userID, adminMail, adminPassword, name, "admin", "", userVersion); err != nil {
		logs.ErrorLogger.Error(err.Error())
		logs.Error.Fatal(err)
	}

	logs.AuditLogger.Printf("ROOT created successfully: %s", adminMail)
	logs.Info.Printf("ROOT created successfully: %s", adminMail)

}
