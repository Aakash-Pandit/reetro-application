package model_tests

import (
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/models"
)

func TestNewFeedback(t *testing.T) {
	createFeedbackRequest := &models.CreateFeedbackRequest{
		Message: "This is a feedback",
	}

	createFeedbackRequest.CreatedBy = TestMockCreateUserResponse()
	createFeedbackRequest.Board = TestMockBoard()
	feedback := models.NewFeedback(createFeedbackRequest)

	if feedback.Message != createFeedbackRequest.Message {
		t.Errorf("returned unexpected output: got %v want %v", feedback.Message, createFeedbackRequest.Message)
	}
}
