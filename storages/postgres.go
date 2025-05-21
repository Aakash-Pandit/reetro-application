package storages

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type PostgresStore struct {
	DB *sql.DB
}

func NewPostgresDB() (*PostgresStore, error) {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_SSL_MODE"))

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres")
	return &PostgresStore{DB: db}, nil
}
