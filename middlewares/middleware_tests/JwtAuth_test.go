package middleware_tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJSONWebToken(t *testing.T) {
	user := TestMockUserResponse()
	token, err := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestValidateJWT(t *testing.T) {
	user := TestMockUserResponse()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	claims, err := middlewares.ValidateJWT(token)

	assert.Nil(t, err)
	assert.NotNil(t, claims)
}

func TestJWTAuthentication(t *testing.T) {
	user := TestMockUserResponse()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	mockStorage := new(MockMiddlewareStorage)
	middleware := middlewares.NewMiddleware(mockStorage)
	handler := middleware.JWTAuthentication(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	token = fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJWTAuthenticationWithEmptyToken(t *testing.T) {
	mockStorage := new(MockMiddlewareStorage)
	middleware := middlewares.NewMiddleware(mockStorage)
	handler := middleware.JWTAuthentication(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTAuthenticationWithIncorrectLength(t *testing.T) {

	mockStorage := new(MockMiddlewareStorage)
	middleware := middlewares.NewMiddleware(mockStorage)
	handler := middleware.JWTAuthentication(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "token")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJWTAuthenticationWithoutBearer(t *testing.T) {
	user := TestMockUserResponse()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	mockStorage := new(MockMiddlewareStorage)
	middleware := middlewares.NewMiddleware(mockStorage)
	handler := middleware.JWTAuthentication(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
