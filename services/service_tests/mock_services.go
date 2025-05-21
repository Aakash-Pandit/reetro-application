package service_tests

import (
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockStorage struct {
	mock.Mock
}

type MockRequestUserStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAllUsers(int, int) ([]*models.CreateUserResponse, error) {
	users := TestMockUsers()
	return users, nil
}

func (m *MockStorage) GetUserById(id uuid.UUID) (*models.CreateUserResponse, error) {
	user := TestMockUserResponse()
	user.Id = id
	return &user, nil
}

func (m *MockStorage) CreateUser(user models.User) (models.User, error) {
	// args := m.Called(user)
	// return args.Get(0).(models.User), args.Error(1)
	return TestMockUser(), nil
}

func (m *MockStorage) DeleteUser(id uuid.UUID) error {
	return nil
}

func (m *MockStorage) VerifyUserByUsername(username string) (*models.User, error) {
	user := TestMockUser()
	user.Username = username
	return &user, nil
}

func (m *MockStorage) VerifyUserByEmail(email string) (*models.User, error) {
	user := TestMockUser()
	user.Email = email
	return &user, nil
}

func (m *MockStorage) SavePassword(user models.User) error {
	return nil
}

func (m *MockStorage) VerifyUserByUsernamePassword(username, password string) (*models.User, error) {
	user := TestMockUser()
	user.Username = username
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user.Password = string(bcryptPassword)
	return &user, nil
}

func (m *MockStorage) GetAllBoards(int, int) ([]*models.Board, error) {
	boards := TestMockBoards()
	return boards, nil
}

func (m *MockStorage) GetBoardById(id uuid.UUID) (*models.Board, error) {
	board := TestMockBoard()
	board.Id = id
	return &board, nil
}

func (m *MockStorage) CreateBoard(board models.Board) (models.Board, error) {
	// args := m.Called()
	// return args.Get(0).(models.Feedback), args.Error(1)
	return board, nil
}

func (m *MockStorage) UpdateBoard(board models.Board) (models.Board, error) {
	return board, nil
}

func (m *MockStorage) DeleteBoard(id uuid.UUID) error {
	return nil
}

func (m *MockStorage) GetAllFeedbacks(int, int) ([]*models.Feedback, error) {
	return TestMockFeedbacks(), nil
}

func (m *MockStorage) GetFeedbackById(id uuid.UUID) (*models.Feedback, error) {
	feedback := TestMockFeedback()
	feedback.Id = id
	return &feedback, nil
}

func (m *MockStorage) CreateFeedback(feedback models.Feedback) (models.Feedback, error) {
	// args := m.Called()
	// return args.Get(0).(models.Feedback), args.Error(1)
	return TestMockFeedback(), nil
}

func (m *MockStorage) UpdateFeedback(feedback models.Feedback) (models.Feedback, error) {
	return feedback, nil
}

func (m *MockStorage) DeleteFeedback(id uuid.UUID) error {
	return nil
}

func (m *MockRequestUserStorage) GetUserById(id uuid.UUID) (*models.CreateUserResponse, error) {
	user := TestMockUserResponse()
	user.Id = id
	return &user, nil
}
