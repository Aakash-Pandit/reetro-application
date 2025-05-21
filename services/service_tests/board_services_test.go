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

func TestGetAllBoardsHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/boards/", nil)
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
	boardService := services.NewBoardService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/boards/", core.HTTPHandleFunc(boardService.GetAllBoardsHandler)).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)

	var listAPIResponse core.ListAPIResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &listAPIResponse)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, 2, listAPIResponse.Count)
	assert.Equal(t, "test board A", listAPIResponse.Result.([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, "test board B", listAPIResponse.Result.([]interface{})[1].(map[string]interface{})["name"])
}

func TestGetBoardByIdHandler(t *testing.T) {
	testBoard := TestMockBoard()
	id := testBoard.Id.String()
	url := fmt.Sprintf("/boards/%s/", id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	mockRepo := new(MockStorage)
	boardService := services.NewBoardService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/boards/{id}/", core.HTTPHandleFunc(boardService.GetBoardByIdHandler)).Methods(http.MethodGet)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, testBoard.Name, "test_board")
}

func TestCreateBoardHandler(t *testing.T) {
	board := TestMockBoard()

	payload, _ := json.Marshal(board)
	req, err := http.NewRequest(http.MethodPost, "/boards/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	user := TestMockUser()
	token, _ := middlewares.GenerateJSONWebToken(user.Id.String(), user.Email)

	token = fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()

	mockRepo := new(MockStorage)
	mockRepo.On("CreateBoard", board).Return(board, nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	boardService := services.NewBoardService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/boards/", core.HTTPHandleFunc(boardService.CreateBoardHandler)).Methods(http.MethodPost)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, board.Name, "test_board")
}

func TestUpdateBoardHandler(t *testing.T) {
	testBoard := TestMockBoard()
	id := testBoard.Id.String()
	url := fmt.Sprintf("/boards/%s/", id)

	payload, _ := json.Marshal(testBoard)
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
	mockRepo.On("UpdateBoard", testBoard).Return(testBoard, nil)

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	boardService := services.NewBoardService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/boards/{id}/", core.HTTPHandleFunc(boardService.UpdateBoardHandler)).Methods(http.MethodPatch)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteBoardHandler(t *testing.T) {
	testBoard := TestMockBoard()
	id := testBoard.Id.String()
	url := fmt.Sprintf("/boards/%s/", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mockRequestUser := new(MockRequestUserStorage)
	mockRedisClient := new(storages_tests.MockRedisClient)
	requestUser := services.NewRequestUser(mockRequestUser)
	mockRepo := new(MockStorage)
	boardService := services.NewBoardService(mockRepo, *requestUser, mockRedisClient)

	r := mux.NewRouter()
	r.HandleFunc("/boards/{id}/", core.HTTPHandleFunc(boardService.DeleteBoardHandler)).Methods(http.MethodDelete)
	r.ServeHTTP(rr, req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rr.Code)
}
