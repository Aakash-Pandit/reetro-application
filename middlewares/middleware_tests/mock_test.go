package middleware_tests

import (
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMiddlewareStorage struct {
	mock.Mock
}

func (m *MockMiddlewareStorage) GetUserById(id uuid.UUID) (*models.CreateUserResponse, error) {
	user := TestMockUserResponse()
	return &user, nil
}
