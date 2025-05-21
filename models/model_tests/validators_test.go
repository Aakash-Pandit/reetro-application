package model_tests

import (
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/models"
)

func TestValidateUserStruct(t *testing.T) {
	createUserRequest := &models.CreateUserRequest{
		FirstName: "test_first_name",
		LastName:  "test_last_name",
		Username:  "test_username",
		Password:  "password",
		Email:     "test@email.com",
		UserType:  "guest_user",
	}

	err := models.ValidateStruct(createUserRequest)
	if err != nil {
		t.Errorf("returned unexpected output: got %v want %v", err, nil)
	}
}

func TestValidateUserStructWithEmptyFirstName(t *testing.T) {
	createUserRequest := &models.CreateUserRequest{
		FirstName: "",
		LastName:  "test_last_name",
	}

	err := models.ValidateStruct(createUserRequest)
	if err == nil {
		t.Errorf("returned unexpected output: got %v want %v", err, "error")
	}
}

func TestValidateFeedbackStruct(t *testing.T) {
	createFeedbackRequest := &models.CreateFeedbackRequest{
		Message: "This is a feedback",
		BoardId: "992c4b2b-f83f-45bc-a258-8b1074fb7a8e",
	}

	err := models.ValidateStruct(createFeedbackRequest)
	if err != nil {
		t.Errorf("returned unexpected output: got %v want %v", err, nil)
	}
}

func TestFeedbackWithEmptyMessage(t *testing.T) {
	createFeedbackRequest := &models.CreateFeedbackRequest{
		Message: "",
	}

	createFeedbackRequest.CreatedBy = TestMockCreateUserResponse()
	createFeedbackRequest.Board = TestMockBoard()
	feedback := models.NewFeedback(createFeedbackRequest)

	if feedback.Message != createFeedbackRequest.Message {
		t.Errorf("returned unexpected output: got %v want %v", feedback.Message, createFeedbackRequest.Message)
	}
}
