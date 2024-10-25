package rolesTest

import (
	"IAM/pkg/handlers/users"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockUserManagementRepository struct {
	mock.Mock
}

func (m *MockUserManagementRepository) DeleteUserFromDB(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserManagementRepository) GetAllUsersDataFromDB() ([]map[string]string, error) {
	return nil, nil
}

func (m *MockUserManagementRepository) GetUsersListFromDB() ([]string, error) {
	return nil, nil
}

func (m *MockUserManagementRepository) GetUserRole(email string) (string, error) {
	return "", nil
}

func TestDeleteUser_Success(t *testing.T) {
	// mocks setting
	mockRepo := new(MockUserManagementRepository)
	mockRepo.On("DeleteUserFromDB", "test@example.com").Return(nil)

	// set new Global repo var
	users.UserManageRepo = mockRepo

	// create HTTP-request
	
}
