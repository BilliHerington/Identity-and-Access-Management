package service

import (
	"IAM/pkg/repository/requests"
	"IAM/pkg/repository/requests/redisAuthentication"
	"IAM/pkg/repository/requests/redisRoles"
	"IAM/pkg/repository/requests/redisUsers"
	"IAM/pkg/service/authenticationServices"
	"IAM/pkg/service/authenticationServices/googleAuthService"
	"IAM/pkg/service/requestLimiterService"
	"IAM/pkg/service/rolesServices"
	"IAM/pkg/service/usersServices"
	"github.com/go-redis/redis/v8"
)

func InitRepositoriesForServices(rdb *redis.Client) {
	//logs.Info.Printf("repositories init. redis exam: %s", rdb)
	rolesServices.RoleManageRepo = &redisRoles.RedisRolesManagementRepository{RDB: rdb}
	usersServices.UserManageRepo = &redisUsers.RedisUserManagementRepository{RDB: rdb}
	authenticationServices.AuthManageRepo = &redisAuthentication.RedisAuthManagementRepository{RDB: rdb}
	googleAuthService.GoogleLoginRepo = &redisAuthentication.RedisAuthManagementRepository{RDB: rdb}
	requestLimiterService.RequestLimitRepo = &requests.RedisRequestRepository{RDB: rdb}
	//authJWT.JwtRepo = &redisJWT.JwtManagementRepository{RDB: rdb}
	//middlewareJWT.JwtRepo = &redisJWT.JwtManagementRepository{RDB: rdb}
}
