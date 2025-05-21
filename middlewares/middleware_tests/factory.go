package middleware_tests

import (
	"time"

	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
)

func TestMockUserResponse() models.CreateUserResponse {
	return models.CreateUserResponse{
		Id:         uuid.New(),
		FirstName:  "test_first_name",
		LastName:   "test_last_name",
		Username:   "test_username",
		Email:      "test@email.com",
		UserType:   "guest_user",
		CreatedAt:  time.Now().UTC(),
		ModifiedAt: time.Now().UTC(),
	}
}
