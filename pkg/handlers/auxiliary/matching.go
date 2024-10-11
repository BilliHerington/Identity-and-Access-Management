package auxiliary

//// AuxRepository interface for working with auxiliary handlers
//type AuxRepository interface {
//	GetRole(role string) (string, error)
//	GetEmail(email string) (string, error)
//	GetUserIDByEmail(email string) (string, error)
//	GetUserVersion(userID string) (string, error)
//}

//// RoleMatch check role exist in redis, return true if exist
//func RoleMatch(repo AuxRepository, role string) (bool, error) {
//	role, err := repo.GetRole(role)
//	if err != nil {
//		return false, err
//	}
//	if len(role) == 0 {
//		return false, nil
//	}
//	return true, nil
//}
//
//// EmailMatch check email exist in redis, return true if exist
//func EmailMatch(repo AuxRepository, email string) (bool, error) {
//	email, err := repo.GetEmail(email)
//	if err != nil {
//		return false, err
//	}
//	if len(email) == 0 {
//		return false, nil
//	}
//	return true, nil
//}
//
//func UserIDByEmail(repo AuxRepository, email string) (string, error) {
//	userID, err := repo.GetUserIDByEmail(email)
//	return userID, err
//}
//func UserVersion(repo AuxRepository, userID string) (string, error) {
//	userVersion, err := repo.GetUserVersion(userID)
//	if err != nil {
//		return "", err
//	}
//	return userVersion, nil
//}
