package storages

import (
	"fmt"
	"log"
	"strings"

	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func ConvertToPQArray(columnTypes []models.ColumnType) string {
	var words []string
	for _, v := range columnTypes {
		words = append(words, fmt.Sprintf("'%s'", v))
	}

	return fmt.Sprintf("{%s}", strings.Join(words, ","))
}

func (p *PostgresStore) GetAllBoards(limit, offset int) ([]*models.Board, error) {
	rows, err := p.DB.Query("SELECT * FROM boards LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*models.Board, 0)
	for rows.Next() {
		board := new(models.Board)

		var columnsString []string
		err := rows.Scan(&board.Id, &board.Name, &board.Template, pq.Array(&columnsString), &board.CreatedById, &board.ModifiedById, &board.CreatedAt, &board.ModifiedAt)
		if err != nil {
			fmt.Printf("Error in scanning the board %v\n", err)
			return nil, err
		}

		stringSlice := make([]string, len(columnsString))
		for i, v := range columnsString {
			stringSlice[i] = string(v)
		}

		for _, v := range stringSlice {
			board.Columns = append(board.Columns, models.ColumnType(strings.ReplaceAll(v, "'", "")))
		}

		board.CreatedBy, _ = p.GetUserById(board.CreatedById)
		board.ModifiedBy, _ = p.GetUserById(board.ModifiedById)

		boards = append(boards, board)
	}

	return boards, nil
}

func (p *PostgresStore) GetBoardById(id uuid.UUID) (*models.Board, error) {
	board := new(models.Board)

	var columnsString []string
	err := p.DB.QueryRow("SELECT * FROM boards WHERE id = $1", id).Scan(&board.Id, &board.Name, &board.Template, pq.Array(&columnsString), &board.CreatedById, &board.ModifiedById, &board.CreatedAt, &board.ModifiedAt)
	if err != nil {
		return nil, err
	}

	stringSlice := make([]string, len(columnsString))
	for i, v := range columnsString {
		stringSlice[i] = string(v)
	}

	for _, v := range stringSlice {
		board.Columns = append(board.Columns, models.ColumnType(strings.ReplaceAll(v, "'", "")))
	}

	board.CreatedBy, _ = p.GetUserById(board.CreatedById)
	board.ModifiedBy, _ = p.GetUserById(board.ModifiedById)

	return board, nil
}

func (p *PostgresStore) CreateBoard(board models.Board) (models.Board, error) {
	arrayInStringFormat := ConvertToPQArray(board.Columns)

	_, err := p.DB.Exec(
		"INSERT INTO boards (id, name, template, columns, created_by_id, modified_by_id, created_at, modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		board.Id, board.Name, board.Template, arrayInStringFormat, board.CreatedById, board.ModifiedById, board.CreatedAt, board.ModifiedAt,
	)
	if err != nil {
		log.Println("Error in creating the user", err)
		return models.Board{}, err
	}

	return board, nil
}

func (p *PostgresStore) UpdateBoard(board models.Board) (models.Board, error) {
	arrayInStringFormat := ConvertToPQArray(board.Columns)

	_, err := p.DB.Exec(
		"UPDATE boards SET name = $1, template = $2, columns = $3, modified_by_id = $4, modified_at = $5 WHERE id = $6",
		board.Name, board.Template, arrayInStringFormat, board.ModifiedById, board.ModifiedAt, board.Id,
	)
	if err != nil {
		log.Println("Error in updating the board", err)
		return models.Board{}, err
	}

	return board, nil
}

func (p *PostgresStore) DeleteBoard(id uuid.UUID) error {
	board := new(models.Board)

	var columnsString []string
	err := p.DB.QueryRow("SELECT * FROM boards WHERE id = $1", id).Scan(&board.Id, &board.Name, &board.Template, pq.Array(&columnsString), &board.CreatedById, &board.ModifiedById, &board.CreatedAt, &board.ModifiedAt)
	if err != nil {
		log.Println("Error while fetching the board", err)
		return err
	}

	_, err = p.DB.Exec("DELETE FROM boards WHERE id = $1", id)
	return err
}
