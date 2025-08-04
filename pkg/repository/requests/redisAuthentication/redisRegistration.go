package redisAuthentication

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"IAM/pkg/repository/requests/redisInternal"
	"errors"
	"github.com/go-redis/redis/v8"
)

func (repo *RedisAuthManagementRepository) StartUserRegistration(userID, email, verificationCode string) error {
	// check email exist
	emailExist, err := redisInternal.CheckEmailExist(repo.RDB, email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if emailExist {
		return models.ErrEmailAlreadyRegistered
	}

	userKey := "user:" + userID
	err = repo.RDB.Watch(ctx, func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			if err = pipe.HSet(ctx, userKey, "id", userID, "email", email, "verificationCode", verificationCode).Err(); err != nil {
				logs.ErrorLogger.Error("failed set fields in userKey", err)
				logs.Error.Println("failed set fields in userKey", err)
			}
			if err = pipe.SAdd(ctx, "users", userID).Err(); err != nil {
				logs.ErrorLogger.Error("failed add userID in users", err)
				logs.Error.Println("failed add userID in users", err)
				return err
			}
			if err = pipe.Set(ctx, "email:"+email, userID, 0).Err(); err != nil {
				logs.ErrorLogger.Error("failed add email key", err)
				logs.Error.Println("failed add email key", err)
				return err
			}
			if err != nil {
				logs.ErrorLogger.Error("Redis TX failed", err)
				logs.Error.Println("Redis TX failed", err)
				return err
			}
			return nil
		})
		if err != nil {
			logs.ErrorLogger.Error("failed watch user", err)
			logs.Error.Println("failed watch user", err)
			return err
		}
		return err
	}, userKey)
	return err
}

func (repo *RedisAuthManagementRepository) SaveUser(email, password, name, role, jwt, userVersion string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return err
	}
	if err = redisInternal.SaveUserInRedis(repo.RDB, userID, email, password, name, role, jwt, userVersion); err != nil {
		logs.ErrorLogger.Error("failed save user in redis", err)
		logs.Error.Println("failed save user in redis", err)
		return err
	}
	return nil
}
func (repo *RedisAuthManagementRepository) GetVerificationCode(email string) (string, error) {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return "", err
		}
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return "", err
	}

	userKey := "user:" + userID
	code, err := repo.RDB.HGet(ctx, userKey, "verificationCode").Result()
	if err != nil {
		logs.Error.Println("failed get verification code", err)
		logs.ErrorLogger.Error("failed get verification code", err)
		return "", err
	}
	//err = repo.RDB.HDel(ctx, userKey, "verificationCode").Err()
	//if err != nil {
	//	logs.Error.Println("failed delete verification code field", err)
	//	logs.ErrorLogger.Error("failed delete verification code field", err)
	//	return "", err
	//}
	return code, nil
}

func (repo *RedisAuthManagementRepository) DeleteVerificationCode(email string) error {
	userID, err := redisInternal.GetUserIDByEmail(repo.RDB, email)
	if err != nil {
		if errors.Is(err, models.ErrUserDoesNotExist) {
			return err
		}
		logs.ErrorLogger.Error(err)
		logs.Error.Println(err)
		return err
	}

	userKey := "user:" + userID
	if err := repo.RDB.HDel(ctx, userKey, "verificationCode").Err(); err != nil {
		logs.Error.Println("failed delete verification code", err)
		logs.ErrorLogger.Error("failed delete verification code", err)
		return err
	}
	return nil
}
