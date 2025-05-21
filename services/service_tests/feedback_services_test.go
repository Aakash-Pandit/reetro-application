package service_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/Aakash-Pandit/reetro-golang/services"
	storages_tests "github.com/Aakash-Pandit/reetro-golang/storages/storage_tests"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFeedbacksHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/feedbacks/", nil)
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
	requestUser := services.NewRequestUser(mockRequestUser)
	mockRepo := new(MockStorage)
	feedbackService := services.NewFeedbackService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/feedbacks/", core.HTTPHandleFunc(feedbackService.GetAllFeedbacksHandler)).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var listAPIResponse core.ListAPIResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &listAPIResponse)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, 2, listAPIResponse.Count)
	assert.Equal(t, "this is feedback A", listAPIResponse.Result.([]interface{})[0].(map[string]interface{})["message"])
	assert.Equal(t, "this is feedback B", listAPIResponse.Result.([]interface{})[1].(map[string]interface{})["message"])
}

func TestGetFeedbackByIdHandler(t *testing.T) {
	testFeedback := TestMockFeedback()
	id := testFeedback.Id.String()
	url := fmt.Sprintf("/feedbacks/%s/", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	mockRepo := new(MockStorage)
	feedbackService := services.NewFeedbackService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/feedbacks/{id}/", core.HTTPHandleFunc(feedbackService.GetFeedbackByIdHandler)).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, testFeedback.Message, "this is feedback")
}

func TestCreateFeedbackHandler(t *testing.T) {
	testFeedback := TestMockFeedback()

	payload, _ := json.Marshal(testFeedback)
	req, err := http.NewRequest(http.MethodPost, "/feedbacks/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	user := TestMockUser()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	token = fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("CreateFeedback", testFeedback).Return(testFeedback, nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	userService := services.NewFeedbackService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/feedbacks/", core.HTTPHandleFunc(userService.CreateFeedbackHandler)).Methods(http.MethodPost)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, testFeedback.Message, "this is feedback")
}

func TestUpdateFeedbackHandler(t *testing.T) {
	testFeedback := TestMockFeedback()
	id := testFeedback.Id.String()
	url := fmt.Sprintf("/feedbacks/%s/", id)

	payload, _ := json.Marshal(testFeedback)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	user := TestMockUser()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	token = fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("UpdateFeedback", testFeedback).Return(testFeedback, nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	feedbackService := services.NewFeedbackService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/feedbacks/{id}/", core.HTTPHandleFunc(feedbackService.UpdateFeedbackHandler)).Methods(http.MethodPatch)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteFeedbackHandler(t *testing.T) {
	testFeedback := TestMockFeedback()
	id := testFeedback.Id.String()
	url := fmt.Sprintf("/feedbacks/%s/", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	mockRepo := new(MockStorage)
	feedbackService := services.NewFeedbackService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/feedbacks/{id}/", core.HTTPHandleFunc(feedbackService.DeleteFeedbackHandler)).Methods(http.MethodDelete)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}
