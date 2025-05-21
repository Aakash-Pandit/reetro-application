package middlewares

import (
	"net/http"

	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/google/uuid"
)

type MiddlewareInterface interface {
	GetUserById(uuid.UUID) (*models.CreateUserResponse, error)
}

func ChainOfMiddleware(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
