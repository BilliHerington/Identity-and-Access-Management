package handlers

import (
	"IAM/pkg/handlers/roles"
	"IAM/pkg/handlers/users"
	"IAM/pkg/redisSystem/redisHandlers/redisRolesHandlers"
	"IAM/pkg/redisSystem/redisHandlers/redisUsersHandlers"
	"github.com/go-redis/redis/v8"
)

func InitRepositories(rdb *redis.Client) {
	roles.RoleManageRepo = &redisRolesHandlers.RedisRolesManagementRepository{RDB: rdb}
	users.UserManageRepo = &redisUsersHandlers.RedisUsersRepository{RDB: rdb}
}
