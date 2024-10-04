package auxiliary

// RoleRepository interface for work with roles
type RoleRepository interface {
	GetRole(role string) (string, error)
}

// RoleMatch check role exist in redis, return true if exist
func RoleMatch(repo RoleRepository, role string) (bool, error) {
	role, err := repo.GetRole(role)
	if err != nil {
		return false, err
	}
	if role == "" {
		return false, nil
	}
	return true, nil
}

// EmailRepository interface for work with emails
type EmailRepository interface {
	GetEmail(email string) (string, error)
}

// EmailMatch check email exist in redis, return true if exist
func EmailMatch(repo EmailRepository, email string) (bool, error) {
	email, err := repo.GetEmail(email)
	if err != nil {
		return false, err
	}
	if email == "" {
		return false, nil
	}
	return true, nil
}
