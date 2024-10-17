package handlers

import (
	"IAM/pkg/handlers/authentication"
	"IAM/pkg/handlers/common"
	"IAM/pkg/handlers/roles"
	"IAM/pkg/handlers/users"
	"IAM/pkg/jwt/authJWT"
	"IAM/pkg/jwt/middlewareJWT"
	"IAM/pkg/repository/requests"
	"IAM/pkg/repository/requests/redisAuthentication"
	"IAM/pkg/repository/requests/redisJWT"
	"IAM/pkg/repository/requests/redisRoles"
	"IAM/pkg/repository/requests/redisUsers"
	"github.com/go-redis/redis/v8"
)

func InitRepositories(rdb *redis.Client) {
	roles.RoleManageRepo = &redisRoles.RedisRolesManagementRepository{RDB: rdb}
	users.UserManageRepo = &redisUsers.RedisUserManagementRepository{RDB: rdb}
	authentication.AuthManageRepo = &redisAuthentication.RedisAuthManagementRepository{RDB: rdb}
	common.RequestLimitRepo = &requests.RedisRequestRepository{RDB: rdb}
	authJWT.JwtRepo = &redisJWT.JwtManagementRepository{RDB: rdb}
	middlewareJWT.JwtRepo = &redisJWT.JwtManagementRepository{RDB: rdb}
}
