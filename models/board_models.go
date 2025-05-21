package models

import (
	"time"

	"github.com/google/uuid"
)

type ColumnType string
type TemplateType string

const (
	GoodThing ColumnType = "good_thing"
	Learned   ColumnType = "learned"
	ShoutOut  ColumnType = "shout_out"
	WentWell  ColumnType = "went_well"
	ToImprove ColumnType = "to_improve"
	Action    ColumnType = "action"
)

const (
	Agile  TemplateType = "agile"
	Kanban TemplateType = "kanban"
	Pacman TemplateType = "pacman"
)

var ValidColumn = []ColumnType{GoodThing, Learned, ShoutOut, WentWell, ToImprove, Action}
var ValidTemplate = []TemplateType{Agile, Kanban, Pacman}

type Board struct {
	Id           uuid.UUID           `json:"id"`
	Name         string              `json:"name"`
	Template     TemplateType        `json:"template"`
	Columns      []ColumnType        `json:"columns"`
	CreatedById  uuid.UUID           `json:"created_by_id"`
	CreatedBy    *CreateUserResponse `json:"created_by"`
	ModifiedById uuid.UUID           `json:"modified_by_id"`
	ModifiedBy   *CreateUserResponse `json:"modified_by"`
	CreatedAt    time.Time           `json:"created_at"`
	ModifiedAt   time.Time           `json:"modified_at"`
}

type CreateBoardRequest struct {
	Name       string              `json:"name" validate:"required"`
	Template   TemplateType        `json:"template" validate:"oneof=agile kanban pacman"`
	Columns    []ColumnType        `json:"columns" validate:"required,min=1,dive,oneof=good_thing learned shout_out went_well to_improve action"`
	CreatedBy  *CreateUserResponse `json:"created_by"`
	ModifiedBy *CreateUserResponse `json:"modified_by"`
}

type UpdateBoardRequest struct {
	Name       string              `json:"name"`
	Template   TemplateType        `json:"template" validate:"oneof=agile kanban pacman"`
	Columns    []ColumnType        `json:"columns" validate:"required,min=1,dive,oneof=good_thing learned shout_out went_well to_improve action"`
	ModifiedBy *CreateUserResponse `json:"modified_by"`
}

func NewBoard(boardRequest *CreateBoardRequest) *Board {
	return &Board{
		Id:           uuid.New(),
		Name:         boardRequest.Name,
		Template:     boardRequest.Template,
		Columns:      boardRequest.Columns,
		CreatedById:  boardRequest.CreatedBy.Id,
		CreatedBy:    boardRequest.CreatedBy,
		ModifiedById: boardRequest.ModifiedBy.Id,
		ModifiedBy:   boardRequest.ModifiedBy,
		CreatedAt:    time.Now().UTC(),
		ModifiedAt:   time.Now().UTC(),
	}
}

func UpdateBoard(board *Board, boardRequest *UpdateBoardRequest) *Board {
	board.Name = boardRequest.Name
	board.Template = boardRequest.Template
	board.Columns = boardRequest.Columns
	board.ModifiedById = boardRequest.ModifiedBy.Id
	board.ModifiedBy = boardRequest.ModifiedBy
	board.ModifiedAt = time.Now().UTC()

	return board
}
