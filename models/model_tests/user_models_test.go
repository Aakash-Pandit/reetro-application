package model_tests

import (
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/models"
)

func TestNewUser(t *testing.T) {
	createUserRequest := &models.CreateUserRequest{
		FirstName: "test_first_name",
		LastName:  "test_last_name",
		Username:  "test_username",
		Password:  "password",
		Email:     "test@email.com",
		UserType:  "guest_user",
	}

	user, err := models.NewUser(createUserRequest)

	if err != nil {
		t.Errorf("returned unexpected error: got %v", err)
	}

	if user.FirstName != createUserRequest.FirstName {
		t.Errorf("returned unexpected output: got %v want %v", user.FirstName, createUserRequest.FirstName)
	}

	if user.Email != createUserRequest.Email {
		t.Errorf("returned unexpected output: got %v want %v", user.Email, createUserRequest.Email)
	}
}
