package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	Id          uuid.UUID           `json:"id"`
	Message     string              `json:"message" validate:"required"`
	BoardId     uuid.UUID           `json:"board_id"`
	Board       *Board              `json:"board"`
	CreatedById uuid.UUID           `json:"created_by_id"`
	CreatedBy   *CreateUserResponse `json:"created_by"`
	CreatedAt   time.Time           `json:"created_at"`
	ModifiedAt  time.Time           `json:"modified_at"`
}

type CreateFeedbackRequest struct {
	Message   string              `json:"message" validate:"required"`
	BoardId   string              `json:"board_id" validate:"required"`
	Board     *Board              `json:"board"`
	CreatedBy *CreateUserResponse `json:"created_by"`
}

type UpdateFeedbackRequest struct {
	Message string `json:"message" validate:"required"`
}

func NewFeedback(feedbackRequest *CreateFeedbackRequest) *Feedback {
	return &Feedback{
		Id:          uuid.New(),
		Message:     feedbackRequest.Message,
		BoardId:     feedbackRequest.Board.Id,
		Board:       feedbackRequest.Board,
		CreatedById: feedbackRequest.CreatedBy.Id,
		CreatedBy:   feedbackRequest.CreatedBy,
		CreatedAt:   time.Now().UTC(),
		ModifiedAt:  time.Now().UTC(),
	}
}

func UpdateFeedback(feedback *Feedback, feedbackRequest *UpdateFeedbackRequest) *Feedback {
	feedback.Message = feedbackRequest.Message
	feedback.ModifiedAt = time.Now().UTC()

	return feedback
}
