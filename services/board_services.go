package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Aakash-Pandit/reetro-golang/common"
	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/models"
	"github.com/Aakash-Pandit/reetro-golang/storages"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BoardService struct {
	Store       storages.Storage
	User        RequestUser
	RedisClient storages.RedisStoreInterface
}

func NewBoardService(store storages.Storage, user RequestUser, redisClient storages.RedisStoreInterface) *BoardService {
	return &BoardService{Store: store, User: user, RedisClient: redisClient}
}

func (b *BoardService) GetAllBoardsHandler(w http.ResponseWriter, r *http.Request) error {
	limit, offset := core.Pagination(r)
	boards, err := b.Store.GetAllBoards(limit, offset)
	if err != nil {
		log.Println("Error in fetching the boards", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	return core.ListAPIResponse(w, &core.ListAPI{
		Status: http.StatusOK,
		Result: &core.ListAPIResponseBody{
			Count:  len(boards),
			Result: boards,
		},
	})
}

func (b *BoardService) GetBoardByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	var cacheBoard models.Board
	board, err := b.RedisClient.Get(id.String(), cacheBoard)
	if board != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusOK,
			Data:   board,
		})
	}

	if err != nil {
		log.Println("Error in fetching the board from cache", err)
	}

	board, err = b.Store.GetBoardById(id)
	if err != nil {
		log.Println("Error in fetching the Board", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Board not found"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   board,
	})
}

func (b *BoardService) CreateBoardHandler(w http.ResponseWriter, r *http.Request) error {
	var boardRequest models.CreateBoardRequest

	userResponse := b.User.GetRequestUser(r)
	if userResponse == nil {
		log.Println("Error while fetching Request User")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Unable to fetch user from token"},
		})
	}

	if userResponse.UserType != models.SuperAdmin {
		log.Println("Only super admin can create Board")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusUnauthorized,
			Data:   &core.APIError{Detail: "Unauthorized to create board"},
		})
	}

	json.NewDecoder(r.Body).Decode(&boardRequest)
	err := models.ValidateStruct(&boardRequest)
	if err != nil {
		log.Println("Error in validating the Board struct", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	defer r.Body.Close()

	boardRequest.CreatedBy = userResponse
	boardRequest.ModifiedBy = userResponse
	board := models.NewBoard(&boardRequest)

	newBoard, store_error := b.Store.CreateBoard(*board)
	if store_error != nil {
		msg := common.AnyToAnyStructField(store_error, &core.DatabaseError{})

		return core.APIResponse(w, &core.Response{
			Status: http.StatusInternalServerError,
			Data:   msg,
		})
	}

	redisErr := b.RedisClient.Set(newBoard.Id.String(), newBoard)
	if redisErr != nil {
		log.Println("Error in setting the board in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusCreated,
		Data:   newBoard,
	})
}

func (b *BoardService) UpdateBoardHandler(w http.ResponseWriter, r *http.Request) error {
	userResponse := b.User.GetRequestUser(r)
	if userResponse == nil {
		log.Println("Error while fetching Request User")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Unable to fetch user from token"},
		})
	}

	if userResponse.UserType != models.SuperAdmin {
		log.Println("Only super admin can update Board")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusUnauthorized,
			Data:   &core.APIError{Detail: "Unauthorized to create board"},
		})
	}

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: err.Error()},
		})
	}

	board, err := b.Store.GetBoardById(id)
	if err != nil {
		log.Println("Error in fetching the Board", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Board not found"},
		})
	}

	var boardRequest models.UpdateBoardRequest

	json.NewDecoder(r.Body).Decode(&boardRequest)
	structErr := models.ValidateStruct(&boardRequest)
	if structErr != nil {
		log.Println("Error in validating the board struct", structErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   structErr,
		})
	}

	defer r.Body.Close()

	boardRequest.ModifiedBy = userResponse
	board = models.UpdateBoard(board, &boardRequest)

	newBoard, store_error := b.Store.UpdateBoard(*board)
	if store_error != nil {
		msg := common.AnyToAnyStructField(store_error, &core.DatabaseError{})

		return core.APIResponse(w, &core.Response{
			Status: http.StatusInternalServerError,
			Data:   msg,
		})
	}

	redisErr := b.RedisClient.Set(newBoard.Id.String(), newBoard)
	if redisErr != nil {
		log.Println("Error in setting the board in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   newBoard,
	})
}

func (b *BoardService) DeleteBoardHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: err.Error()},
		})
	}

	err = b.Store.DeleteBoard(id)
	if err != nil {
		log.Println("Error in fetching the Board:", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Board not found"},
		})
	}

	redisErr := b.RedisClient.Del(id.String())
	if redisErr != nil {
		log.Println("Error in setting the board in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   map[string]string{"detail": "Board deleted successfully"},
	})
}
