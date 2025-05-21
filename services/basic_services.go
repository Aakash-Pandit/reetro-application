package services

import (
	"net/http"
	"strings"

	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/Aakash-Pandit/reetro-golang/storages"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AboutPageResponse struct {
	Detail      string `json:"detail"`
	Information string `json:"information"`
}

type ForgotPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordPayload struct {
	Username    string `json:"username" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type APISuccessResponse struct {
	Detail string `json:"detail"`
}

type ServiceStorage interface {
	GetUserById(uuid.UUID) (*models.CreateUserResponse, error)
}

type RequestUser struct {
	Store ServiceStorage
}

func NewRequestUser(store ServiceStorage) *RequestUser {
	return &RequestUser{Store: store}
}

func (u *RequestUser) GetRequestUser(r *http.Request) *models.CreateUserResponse {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		return nil
	}

	token_info := strings.Split(tokenString, " ")
	if len(token_info) != 2 {
		return nil
	}

	if token_info[0] != "Bearer" {
		return nil
	}

	token, err := middlewares.ValidateJWT(token_info[1])
	if err != nil {
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)
	id, _ := uuid.Parse(claims["id"].(string))
	user, err := u.Store.GetUserById(id)

	if err != nil {
		return nil
	}

	return user
}

type BasicService struct {
	RedisClient storages.RedisStoreInterface
}

func HomeHandler(w http.ResponseWriter, r *http.Request) error {
	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   APISuccessResponse{Detail: "Reetro Application Home Page"},
	})
}

func AboutHandler(w http.ResponseWriter, r *http.Request) error {
	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data: AboutPageResponse{
			Detail:      "Reetro Application About Page",
			Information: "This application is used for retrospective.",
		},
	})
}

func (b *BasicService) ClearRedisCache(w http.ResponseWriter, r *http.Request) error {
	err := b.RedisClient.FlushAll()
	if err != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusInternalServerError,
			Data:   &core.APIError{Detail: "Error in clearing the Redis Cache"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   APISuccessResponse{Detail: "Redis Cache Cleared Successfully"},
	})
}
