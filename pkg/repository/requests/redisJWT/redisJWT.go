package redisJWT

import (
	"IAM/pkg/logs"
	"IAM/pkg/repository/requests/redisInternal"
	"context"
	"github.com/go-redis/redis/v8"
)

type JwtManagementRepository struct {
	RDB *redis.Client
}

var ctx = context.Background()

func (repo *JwtManagementRepository) GetDataForJWT(email string) (userID string, userVersion string, error error) {

	// get userID
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return "", "", err
		}
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return "", "", err
	}

	// get userVersion
	userVersion, err = redisInternal.GetUserVersion(repo.RDB, userID)
	if err != nil {
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return "", "", err
	}

	return userID, userVersion, nil
}
func (repo *JwtManagementRepository) SetJWT(email, jwt string) error {

	// get userID
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if err.Error() == "user does not exist" {
			return err
		}
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return err
	}

	// set jwt
	if err = repo.RDB.HSet(ctx, "user:"+userID, "jwt", jwt).Err(); err != nil {
		logs.Error.Println("failed set jwt", err)
		logs.ErrorLogger.Error("failed set jwt", err)
		return err
	}
	return nil
}
