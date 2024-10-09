package auxiliary

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

// RoleRepository interface for work with roles
type RoleRepository interface {
	GetRole(role string) (string, error)
}

// RoleMatch check role exist in redis, return true if exist
func RoleMatch(repo RoleRepository, role string) (bool, error) {
	role, err := repo.GetRole(role)
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
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
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
