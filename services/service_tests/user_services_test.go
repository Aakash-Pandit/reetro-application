package service_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/common/common_tests"
	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/Aakash-Pandit/reetro-golang/services"
	storages_tests "github.com/Aakash-Pandit/reetro-golang/storages/storage_tests"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/users/", nil)
	if err != nil {
		t.Fatal(err)
	}

	user := TestMockUser()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	token = fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)

	mockRepo := new(MockStorage)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/users/", core.HTTPHandleFunc(userService.GetAllUsersHandler)).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var listAPIResponse core.ListAPIResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &listAPIResponse)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, 2, listAPIResponse.Count)

	assert.Equal(t, "test_first_name", listAPIResponse.Result.([]interface{})[0].(map[string]interface{})["first_name"])
	assert.Equal(t, "testb@email.com", listAPIResponse.Result.([]interface{})[0].(map[string]interface{})["email"])

	assert.Equal(t, "testa@email.com", listAPIResponse.Result.([]interface{})[1].(map[string]interface{})["email"])
}

func TestGetUserByIdHandler(t *testing.T) {
	testUser := TestMockUser()
	id := testUser.Id.String()
	url := fmt.Sprintf("/users/%s/", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)

	mockRepo := new(MockStorage)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/", core.HTTPHandleFunc(userService.GetUserByIdHandler)).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, testUser.Email, "test@email.com")
}

func TestCreateUserHandler(t *testing.T) {
	testUser := TestMockUser()

	payload, _ := json.Marshal(testUser)
	req, err := http.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("CreateUser", testUser).Return(testUser, nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/users/", core.HTTPHandleFunc(userService.CreateUserHandler)).Methods(http.MethodPost)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, testUser.Email, "test@email.com")
}

func TestDeleteUserHandler(t *testing.T) {
	testUser := TestMockUser()
	id := testUser.Id.String()
	url := fmt.Sprintf("/users/%s/", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)

	mockRepo := new(MockStorage)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/users/{id}/", core.HTTPHandleFunc(userService.DeleteUserHandler)).Methods(http.MethodDelete)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLoginHandler(t *testing.T) {
	testLogin := TestLoginPayload()
	payload, _ := json.Marshal(testLogin)
	req, err := http.NewRequest(http.MethodPost, "/login/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("VerifyUserByUsername", testLogin.Username).Return(&testLogin, nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/login/", core.HTTPHandleFunc(userService.LoginHandler)).Methods(http.MethodPost)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tokenResponse core.TokenResponse
	err = json.Unmarshal(rr.Body.Bytes(), &tokenResponse)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	assert.NotEmpty(t, tokenResponse.Token)
}

func TestResetPasswordHandler(t *testing.T) {
	testUser := TestMockUser()
	resetPasswordPayload := TestResetPasswordPayload()

	payload, _ := json.Marshal(resetPasswordPayload)
	req, err := http.NewRequest(http.MethodPost, "/reset_password/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("VerifyUserByUsernamePassword", resetPasswordPayload.Username, resetPasswordPayload.OldPassword).Return(&testUser, nil)
	mockRepo.On("SavePassword", testUser).Return(nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/reset_password/", core.HTTPHandleFunc(userService.ResetPasswordHandler)).Methods(http.MethodPost)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestForgotPasswordHandler(t *testing.T) {
	testUser := TestMockUser()
	ForgotPasswordPayload := TestForgotPasswordPayload()

	payload, _ := json.Marshal(ForgotPasswordPayload)
	req, err := http.NewRequest(http.MethodPost, "/forgot_password/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("VerifyUserByEmail", testUser.Email).Return(&testUser, nil)
	mockRepo.On("SavePassword", testUser).Return(nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	mockEmail := new(common_tests.MockEmail)
	requestUser := services.NewRequestUser(mockRequestUser)
	userService := services.NewUserService(mockRepo, *requestUser, mockRedisClient, mockEmail)

	r := mux.NewRouter()
	r.HandleFunc("/forgot_password/", core.HTTPHandleFunc(userService.ForgotPasswordHandler)).Methods(http.MethodPost)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}
