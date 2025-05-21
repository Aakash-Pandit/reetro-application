package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	GuestUser  Role = "guest_user"
	TeamMember Role = "team_member"
	SuperAdmin Role = "super_admin"
)

var ValidUserType = []Role{GuestUser, TeamMember, SuperAdmin}

type User struct {
	Id         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	UserType   Role      `json:"user_type"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	UserType  Role   `json:"user_type"`
}

type CreateUserResponse struct {
	Id         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	UserType   Role      `json:"user_type"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewUser(userRequest *CreateUserRequest) (*User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error while hashing password: ", err)
		return nil, err
	}

	return &User{
		Id:         uuid.New(),
		FirstName:  userRequest.FirstName,
		LastName:   userRequest.LastName,
		Username:   userRequest.Username,
		Password:   string(password),
		Email:      userRequest.Email,
		UserType:   userRequest.UserType,
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}, nil
}
