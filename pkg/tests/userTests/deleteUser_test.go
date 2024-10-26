package userTests

import (
	"IAM/pkg/handlers/users"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
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
	body := map[string]string{"email": "test@example.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/deleteUser", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// use httptest for recording answer
	w := httptest.NewRecorder()

	// create Gin-context and call handler
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// call func DeleteUser
	handler := users.DeleteUser()
	handler(c)

	// check result
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"user deleted":"test@example.com"}`, w.Body.String())
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_UserDoesNotExist(t *testing.T) {
	// mocks setting
	mockRepo := new(MockUserManagementRepository)
	mockRepo.On("DeleteUserFromDB", "nonexistent@example.com").Return(errors.New("user does not exist"))

	// set new Global repo var
	users.UserManageRepo = mockRepo

	// create HTTP-request
	body := map[string]string{"email": "nonexistent@example.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/delete-user", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// use httptest for recording answer
	w := httptest.NewRecorder()

	// create Gin-context and call handler
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// call func DeleteUser
	handler := users.DeleteUser()
	handler(c)

	// check result
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "user does not exist"}`, w.Body.String())
	mockRepo.AssertExpectations(t)
}
func TestDeleteUser_InvalidEmail(t *testing.T) {
	//mockRepo := new(MockUserManagementRepository)
	//
	//// Здесь не надо настраивать mockRepo, потому что оно не должно быть вызвано при ошибке валидации
	//users.UserManageRepo = mockRepo

	// Некорректный email
	body := map[string]string{"email": "emailWithoutMailTag.com"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/delete-user", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Вызываем обработчик
	handler := users.DeleteUser()
	handler(c)

	// Ожидаем, что код ответа будет 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверяем тело ответа
	assert.JSONEq(t, `{"error": "incorrect data format, please check your input data"}`, w.Body.String())

	// Проверяем, что метод DeleteUserFromDB не был вызван
	//mockRepo.AssertNotCalled(t, "DeleteUserFromDB", "emailWithoutMailTag.com")
}
