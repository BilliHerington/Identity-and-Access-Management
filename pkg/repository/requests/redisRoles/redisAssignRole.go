package redisRoles

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisRolesManagementRepository struct {
	RDB *redis.Client
}

func (repo *RedisRolesManagementRepository) AssignRoleToUser(email, role string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	roleExist, err := redisInternal.CheckRoleExist(repo.RDB, role)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if !roleExist {
		return models.ErrRoleDoesNotExist
	}
	err = repo.RDB.HSet(ctx, "user:"+userID, "role", role).Err()
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return fmt.Errorf("failed to add role to user: %w", err)
	}
	return err
}
