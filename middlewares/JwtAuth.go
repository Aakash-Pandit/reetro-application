package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Middleware struct {
	Store MiddlewareInterface
}

func NewMiddleware(store MiddlewareInterface) *Middleware {
	return &Middleware{Store: store}
}

func GenerateJSONWebToken(id, email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := jwt.MapClaims{
		"authorized": true,
		"id":         id,
		"email":      email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
}

func (m *Middleware) JWTAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			core.APIResponse(w, &core.Response{
				Status: http.StatusUnauthorized,
				Data:   &core.APIError{Detail: "Empty Token"},
			})
			return
		}

		token_info := strings.Split(tokenString, " ")
		if len(token_info) != 2 {
			core.APIResponse(w, &core.Response{
				Status: http.StatusUnauthorized,
				Data:   &core.APIError{Detail: "Invalid token"},
			})
			return
		}

		if token_info[0] != "Bearer" {
			core.APIResponse(w, &core.Response{
				Status: http.StatusUnauthorized,
				Data:   &core.APIError{Detail: "Bearer token is required"},
			})
			return
		}

		token, err := ValidateJWT(token_info[1])
		if err != nil {
			core.APIResponse(w, &core.Response{
				Status: http.StatusUnauthorized,
				Data:   &core.APIError{Detail: "Invalid token"},
			})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		id, _ := uuid.Parse(claims["id"].(string))
		_, err = m.Store.GetUserById(id)

		if err != nil {
			core.APIResponse(w, &core.Response{
				Status: http.StatusUnauthorized,
				Data:   &core.APIError{Detail: "Unauthorized"},
			})
			return
		}

		next(w, r)
	}
}
