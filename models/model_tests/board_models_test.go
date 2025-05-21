package model_tests

import (
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/models"
)

func TestNewBoard(t *testing.T) {
	createBoardRequest := &models.CreateBoardRequest{
		Name:     "test_board_name",
		Template: models.Agile,
		Columns:  []models.ColumnType{models.WentWell, models.ToImprove, models.Action},
	}

	user := TestMockCreateUserResponse()
	createBoardRequest.CreatedBy = user
	createBoardRequest.ModifiedBy = user
	board := models.NewBoard(createBoardRequest)

	if board.Name != createBoardRequest.Name {
		t.Errorf("returned unexpected output: got %v want %v", board.Name, createBoardRequest.Name)
	}
}
