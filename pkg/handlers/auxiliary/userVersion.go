package auxiliary

type UserVersionRepository interface {
	GetUserVersion(userID string) (string, error)
}

func UserVersion(repo UserVersionRepository, userID string) (string, error) {
	userVersion, err := repo.GetUserVersion(userID)
	if err != nil {
		return "", err
	}
	return userVersion, nil
}
