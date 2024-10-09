package auxiliary

type UserIDByEmailRepository interface {
	GetUserIDByEmail(email string) (string, error)
}

func UserIDByEmail(repo UserIDByEmailRepository, email string) (string, error) {
	userID, err := repo.GetUserIDByEmail(email)
	if err != nil {
		return "", err
	}
	return userID, nil
}
