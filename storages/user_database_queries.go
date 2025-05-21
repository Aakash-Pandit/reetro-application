package storages

import (
	"log"

	"github.com/Aakash-Pandit/reetro-golang/common"
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (p *PostgresStore) GetAllUsers(limit, offset int) ([]*models.CreateUserResponse, error) {
	rows, err := p.DB.Query("SELECT * FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.CreateUserResponse, 0)
	for rows.Next() {
		user := new(models.User)
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.Email, &user.UserType, &user.CreatedAt, &user.ModifiedAt)
		if err != nil {
			return nil, err
		}

		msg := common.AnyToAnyStructField(user, &models.CreateUserResponse{})

		users = append(users, msg.(*models.CreateUserResponse))
	}

	return users, nil
}

func (p *PostgresStore) GetUserById(id uuid.UUID) (*models.CreateUserResponse, error) {
	user := new(models.User)

	err := p.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.Email, &user.UserType, &user.CreatedAt, &user.ModifiedAt)
	if err != nil {
		return nil, err
	}

	msg := common.AnyToAnyStructField(user, &models.CreateUserResponse{})

	return msg.(*models.CreateUserResponse), nil
}

func (p *PostgresStore) CreateUser(user models.User) (models.User, error) {
	_, err := p.DB.Exec(
		"INSERT INTO users (id, first_name, last_name, username, password, email, user_type, created_at, modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		user.Id, user.FirstName, user.LastName, user.Username, user.Password, user.Email, user.UserType, user.CreatedAt, user.ModifiedAt,
	)
	if err != nil {
		log.Println("Error in creating the user", err)
		return models.User{}, err
	}

	return user, nil
}

func (p *PostgresStore) DeleteUser(id uuid.UUID) error {
	user := new(models.User)

	err := p.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.Email, &user.UserType, &user.CreatedAt, &user.ModifiedAt)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (p *PostgresStore) VerifyUserByUsername(username string) (*models.User, error) {
	user := new(models.User)
	err := p.DB.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.Email, &user.UserType, &user.CreatedAt, &user.ModifiedAt)

	if err != nil {
		log.Println("Error in fetching the user", err)
		return nil, err
	}

	return user, nil
}

func (p *PostgresStore) VerifyUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := p.DB.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.Email, &user.UserType, &user.CreatedAt, &user.ModifiedAt)

	if err != nil {
		log.Println("Error in fetching the user", err)
		return nil, err
	}

	return user, nil
}

func (p *PostgresStore) SavePassword(user models.User) error {
	_, err := p.DB.Exec("UPDATE users SET password = $1 WHERE id = $2", user.Password, user.Id)
	if err != nil {
		log.Println("Error in updating the password", err)
		return err
	}

	return nil
}

func (p *PostgresStore) VerifyUserByUsernamePassword(username, password string) (*models.User, error) {
	user := new(models.User)
	err := p.DB.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Password, &user.Email, &user.UserType, &user.CreatedAt, &user.ModifiedAt)

	if err != nil {
		log.Println("Invalid Username", err)
		return nil, err
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passwordError != nil {
		log.Println("Invalid Password", passwordError)
		return nil, passwordError
	}

	return user, nil
}
