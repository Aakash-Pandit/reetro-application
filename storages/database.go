package storages

import (
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetAllUsers(int, int) ([]*models.CreateUserResponse, error)
	GetUserById(uuid.UUID) (*models.CreateUserResponse, error)
	CreateUser(models.User) (models.User, error)
	DeleteUser(uuid.UUID) error
	VerifyUserByUsername(string) (*models.User, error)
	VerifyUserByEmail(string) (*models.User, error)
	SavePassword(models.User) error
	VerifyUserByUsernamePassword(string, string) (*models.User, error)

	GetAllBoards(int, int) ([]*models.Board, error)
	GetBoardById(uuid.UUID) (*models.Board, error)
	CreateBoard(models.Board) (models.Board, error)
	UpdateBoard(models.Board) (models.Board, error)
	DeleteBoard(uuid.UUID) error

	GetAllFeedbacks(int, int) ([]*models.Feedback, error)
	GetFeedbackById(uuid.UUID) (*models.Feedback, error)
	CreateFeedback(models.Feedback) (models.Feedback, error)
	UpdateFeedback(models.Feedback) (models.Feedback, error)
	DeleteFeedback(uuid.UUID) error
}

type Database struct {
	Postgres *PostgresStore
	Redis    *RedisStore
}

func DBInit() (*Database, error) {
	PG, err := NewPostgresDB()
	if err != nil {
		return nil, err
	}

	RD, err := NewRedisDB()
	if err != nil {
		return nil, err
	}

	return &Database{Postgres: PG, Redis: RD}, nil
}

func (d *Database) Close() {
	d.Postgres.DB.Close()
	d.Redis.Client.Close()
}
