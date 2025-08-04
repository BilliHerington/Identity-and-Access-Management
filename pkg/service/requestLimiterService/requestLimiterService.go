package requestLimiterService

import (
	"IAM/pkg/logs"
	"IAM/pkg/models"
	"fmt"
)

type RequestLimitRepository interface {
	GetRequestLimit(redisKey string, limit int, window int64) (bool, error)
}

var RequestLimitRepo RequestLimitRepository

func RequestLimiterService(userID, clientIP string, limit int, window int64) error {
	var redisKey string

	if userID != "" {
		redisKey = fmt.Sprintf("rate_limit_%s", userID)
	} else {
		// if user not authorized, use IP
		redisKey = fmt.Sprintf("rate_limit_ip_%s", clientIP)
	}

	// check requests exceed
	exceeded, err := RequestLimitRepo.GetRequestLimit(redisKey, limit, window)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return err
	}
	if exceeded {
		return models.ErrRequestLimitExceeded
	}
	return nil
}
