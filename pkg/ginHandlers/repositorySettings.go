package ginHandlers

import (
	"IAM/pkg/jwt/authJWT"
	"IAM/pkg/jwt/middlewareJWT"
	//"IAM/pkg/repository/requests/redisAuthentication"
	"IAM/pkg/repository/requests/redisJWT"
	"github.com/go-redis/redis/v8"
)

func InitRepositories(rdb *redis.Client) {
	//logs.Info.Printf("repositories init. redis exam: %s", rdb)
	//rolesHandlers.RoleManageRepo = &redisRoles.RedisRolesManagementRepository{RDB: rdb}
	//users.UserManageRepo = &redisUsers.RedisUserManagementRepository{RDB: rdb}
	//authenticationHandlers.AuthManageRepo = &redisAuthentication.RedisAuthManagementRepository{RDB: rdb}
	//googleAuth.GoogleLoginRepo = &redisAuthentication.RedisAuthManagementRepository{RDB: rdb}
	//requestLimiterHandler.RequestLimitRepo = &requests.RedisRequestRepository{RDB: rdb}
	authJWT.JwtRepo = &redisJWT.JwtManagementRepository{RDB: rdb}
	middlewareJWT.JwtRepo = &redisJWT.JwtManagementRepository{RDB: rdb}
}
