package redisDB

import (
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

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	//password := ""
	intDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logs.Error.Printf("invalid REDIS_DB value %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       intDb,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		logs.ErrorLogger.Errorln(err)
		logs.Info.Printf("Trying connect to redis on address: %s, with password:%s", addr, password)
		logs.Error.Fatalf("redis ping failed: %v", err)
	}
	initFirstTime := os.Getenv("INIT_REDIS_IN_FIRST_TIME")
	//logrus.Warn("test")
	if initFirstTime == "true" {
		logs.Info.Printf("Starting project in first time\nCheck .env if WANT change this\nCreating base DB structure...\n")
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

	defRoleNameAdmin := "admin"
	defRoleNameUser := "user"
	// create roles
	repository := redisRoles.RedisRolesManagementRepository{RDB: rdb}
	if err := repository.CreateRole(defRoleNameAdmin, adminPrivileges); err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Fatal(err)
	}
	if err := repository.CreateRole(defRoleNameUser, userPrivileges); err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Fatal(err)
	}

	logs.AuditLogger.Printf("roles '%s' and '%s' created successfully", defRoleNameAdmin, defRoleNameUser)
	logs.Info.Printf("roles '%s' and '%s' created successfully", defRoleNameAdmin, defRoleNameUser)

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
