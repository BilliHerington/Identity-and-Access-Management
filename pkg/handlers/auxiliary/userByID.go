package auxiliary

import (
	"IAM/pkg/logs"
)

type UserIDByEmailRepository interface {
	GetUserIDByEmail(email string) (string, error)
}

func UserIDByEmail(repo UserIDByEmailRepository, email string) (string, error) {
	userID, err := repo.GetUserIDByEmail(email)
	if err != nil {
		logs.Error.Println(err)
		logs.ErrorLogger.Error(err)
		return "", err
	}
	if userID == "" {
		return "", nil
	}
	return userID, nil
}
