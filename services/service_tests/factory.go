package service_tests

import (
	"time"

	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/Aakash-Pandit/reetro-golang/services"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestMockUsers() []*models.CreateUserResponse {
	userA := models.CreateUserResponse{
		Id:         uuid.New(),
		FirstName:  "test_first_name",
		LastName:   "test_last_name",
		Username:   "username B",
		Email:      "testb@email.com",
		UserType:   "super_admin",
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
	userB := models.CreateUserResponse{
		Id:         uuid.New(),
		FirstName:  "test_first_name",
		LastName:   "test_last_name",
		Username:   "username A",
		Email:      "testa@email.com",
		UserType:   "super_admin",
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}

	users := []*models.CreateUserResponse{
		&userA, &userB,
	}

	return users
}

func TestMockUser() models.User {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	return models.User{
		Id:         uuid.New(),
		FirstName:  "test_first_name",
		LastName:   "test_last_name",
		Username:   "test_username",
		Password:   string(password),
		Email:      "test@email.com",
		UserType:   "super_admin",
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
}

func TestMockUserResponse() models.CreateUserResponse {
	return models.CreateUserResponse{
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

func TestLoginPayload() models.LoginRequest {
	return models.LoginRequest{
		Username: "test_username",
		Password: "password",
	}
}

func TestForgotPasswordPayload() services.ForgotPasswordPayload {
	return services.ForgotPasswordPayload{
		Email: "test_email@email.com",
	}
}

func TestResetPasswordPayload() services.ResetPasswordPayload {
	return services.ResetPasswordPayload{
		Username:    "test_username",
		OldPassword: "old_password",
		NewPassword: "new_password",
	}
}

func TestMockBoards() []*models.Board {
	user := TestMockUserResponse()
	boardA := models.Board{
		Id:           uuid.New(),
		Name:         "test board A",
		Template:     models.Agile,
		Columns:      []models.ColumnType{models.WentWell, models.ToImprove, models.Action},
		CreatedById:  user.Id,
		ModifiedById: user.Id,
		CreatedBy:    &user,
		ModifiedBy:   &user,
		CreatedAt:    time.Now().UTC(),
		ModifiedAt:   time.Now().UTC(),
	}

	boardB := models.Board{
		Id:           uuid.New(),
		Name:         "test board B",
		Template:     models.Agile,
		Columns:      []models.ColumnType{models.WentWell, models.ToImprove, models.Action},
		CreatedById:  user.Id,
		ModifiedById: user.Id,
		CreatedBy:    &user,
		ModifiedBy:   &user,
		CreatedAt:    time.Now().UTC(),
		ModifiedAt:   time.Now().UTC(),
	}

	boards := []*models.Board{
		&boardA, &boardB,
	}

	return boards
}

func TestMockBoard() models.Board {
	return models.Board{
		Id:         uuid.New(),
		Name:       "test_board",
		Template:   models.Agile,
		Columns:    []models.ColumnType{models.WentWell, models.ToImprove, models.Action},
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
}

func TestMockFeedbacks() []*models.Feedback {
	user := TestMockUserResponse()
	feedbackA := models.Feedback{
		Id:          uuid.New(),
		Message:     "this is feedback A",
		CreatedById: user.Id,
		CreatedBy:   &user,
		CreatedAt:   time.Now().UTC(),
		ModifiedAt:  time.Now().UTC(),
	}

	feedbackB := models.Feedback{
		Id:          uuid.New(),
		Message:     "this is feedback B",
		CreatedById: user.Id,
		CreatedBy:   &user,
		CreatedAt:   time.Now().UTC(),
		ModifiedAt:  time.Now().UTC(),
	}

	return []*models.Feedback{
		&feedbackA, &feedbackB,
	}
}

func TestMockFeedback() models.Feedback {
	return models.Feedback{
		Id:         uuid.New(),
		Message:    "this is feedback",
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
}
