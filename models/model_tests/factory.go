package model_tests

import (
	"time"

	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
)

func TestMockCreateUserResponse() *models.CreateUserResponse {
	return &models.CreateUserResponse{
		Id:         uuid.New(),
		FirstName:  "test_first_name",
		LastName:   "test_last_name",
		Username:   "test_username",
		Email:      "test@email.com",
		UserType:   "super_admin",
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
}

func TestMockBoard() *models.Board {
	return &models.Board{
		Id:       uuid.New(),
		Name:     "test_board_name",
		Template: models.Agile,
		Columns:  []models.ColumnType{models.WentWell, models.ToImprove, models.Action},
	}
}
