package common_tests

import "github.com/stretchr/testify/mock"

type MockEmail struct {
	mock.Mock
}

func (m *MockEmail) SendEmailForPasswordReset(emailId, subject, password string) error {
	return nil
}
