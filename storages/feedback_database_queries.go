package storages

import (
	"fmt"
	"log"

	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
)

func (p *PostgresStore) GetAllFeedbacks(limit, offset int) ([]*models.Feedback, error) {
	rows, err := p.DB.Query("SELECT * FROM feedbacks LIMIT $1 OFFSET $2", limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := make([]*models.Feedback, 0)
	for rows.Next() {
		feedback := new(models.Feedback)
		err := rows.Scan(&feedback.Id, &feedback.Message, &feedback.BoardId, &feedback.CreatedById, &feedback.CreatedAt, &feedback.ModifiedAt)
		if err != nil {
			fmt.Printf("Error in scanning the feedback %v\n", err)
			return nil, err
		}

		feedback.CreatedBy, _ = p.GetUserById(feedback.CreatedById)

		feedbacks = append(feedbacks, feedback)
	}

	return feedbacks, nil
}

func (p *PostgresStore) GetFeedbackById(id uuid.UUID) (*models.Feedback, error) {
	feedback := new(models.Feedback)

	err := p.DB.QueryRow("SELECT * FROM feedbacks WHERE id = $1", id).Scan(&feedback.Id, &feedback.Message, &feedback.BoardId, &feedback.CreatedById, &feedback.CreatedAt, &feedback.ModifiedAt)
	if err != nil {
		return nil, err
	}

	feedback.Board, _ = p.GetBoardById(feedback.BoardId)
	feedback.CreatedBy, _ = p.GetUserById(feedback.CreatedById)
	return feedback, nil
}

func (p *PostgresStore) CreateFeedback(feedback models.Feedback) (models.Feedback, error) {
	_, err := p.DB.Exec(
		"INSERT INTO feedbacks (id, message, board_id, created_by_id, created_at, modified_at) VALUES ($1, $2, $3, $4, $5, $6)",
		feedback.Id, feedback.Message, feedback.BoardId, feedback.CreatedById, feedback.CreatedAt, feedback.ModifiedAt,
	)
	if err != nil {
		log.Println("Error in creating the user", err)
		return models.Feedback{}, err
	}

	return feedback, nil
}

func (p *PostgresStore) UpdateFeedback(feedback models.Feedback) (models.Feedback, error) {
	_, err := p.DB.Exec(
		"UPDATE feedbacks SET message = $1, modified_at = $2 WHERE id = $3",
		feedback.Message, feedback.ModifiedAt, feedback.Id,
	)
	if err != nil {
		log.Println("Error in updating the board", err)
		return models.Feedback{}, err
	}

	return feedback, nil
}

func (p *PostgresStore) DeleteFeedback(id uuid.UUID) error {
	feedback := new(models.Feedback)

	err := p.DB.QueryRow("SELECT * FROM feedbacks WHERE id = $1", id).Scan(&feedback.Id, &feedback.Message, &feedback.BoardId, &feedback.CreatedById, &feedback.CreatedAt, &feedback.ModifiedAt)
	if err != nil {
		log.Println("Error while fetching the feedback", err)
		return err
	}

	_, err = p.DB.Exec("DELETE FROM feedbacks WHERE id = $1", id)
	return err
}
