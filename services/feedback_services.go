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

type FeedbackService struct {
	Store       storages.Storage
	User        RequestUser
	RedisClient storages.RedisStoreInterface
}

func NewFeedbackService(store storages.Storage, user RequestUser, redisClient storages.RedisStoreInterface) *FeedbackService {
	return &FeedbackService{Store: store, User: user, RedisClient: redisClient}
}

func (f *FeedbackService) GetAllFeedbacksHandler(w http.ResponseWriter, r *http.Request) error {
	limit, offset := core.Pagination(r)
	feedbacks, err := f.Store.GetAllFeedbacks(limit, offset)
	if err != nil {
		log.Println("Error in fetching the feedbacks", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	return core.ListAPIResponse(w, &core.ListAPI{
		Status: http.StatusOK,
		Result: &core.ListAPIResponseBody{
			Count:  len(feedbacks),
			Result: feedbacks,
		},
	})
}

func (f *FeedbackService) GetFeedbackByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	var cacheFeedback models.Feedback
	feedback, err := f.RedisClient.Get(id.String(), cacheFeedback)
	if feedback != nil {
		return core.APIResponse(w, &core.Response{
			Status: http.StatusOK,
			Data:   feedback,
		})
	}

	if err != nil {
		log.Println("Error in fetching the feedback from cache", err)
	}

	feedback, err = f.Store.GetFeedbackById(id)
	if err != nil {
		log.Println("Error in fetching the Feedback", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Feedback not found"},
		})
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   feedback,
	})
}

func (f *FeedbackService) CreateFeedbackHandler(w http.ResponseWriter, r *http.Request) error {
	var feedbackRequest models.CreateFeedbackRequest

	userResponse := f.User.GetRequestUser(r)
	if userResponse == nil {
		log.Println("Error while fetching Request User")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Unable to fetch user from token"},
		})
	}

	json.NewDecoder(r.Body).Decode(&feedbackRequest)
	err := models.ValidateStruct(&feedbackRequest)
	if err != nil {
		log.Println("Error in validating the Feedback struct", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   err,
		})
	}

	defer r.Body.Close()

	boardId, idErr := uuid.Parse(feedbackRequest.BoardId)
	if idErr != nil {
		log.Println("Error in parsing the board id", idErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: idErr.Error()},
		})
	}

	board, dbErr := f.Store.GetBoardById(boardId)
	if dbErr != nil {
		log.Println("Error while fetching Request Board")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Unable to fetch board from database"},
		})
	}

	feedbackRequest.Board = board
	feedbackRequest.CreatedBy = userResponse
	feedback := models.NewFeedback(&feedbackRequest)

	newFeedback, store_error := f.Store.CreateFeedback(*feedback)
	if store_error != nil {
		msg := common.AnyToAnyStructField(store_error, &core.DatabaseError{})
		return core.APIResponse(w, &core.Response{
			Status: http.StatusInternalServerError,
			Data:   msg,
		})
	}

	redisErr := f.RedisClient.Set(newFeedback.Id.String(), newFeedback)
	if redisErr != nil {
		log.Println("Error in setting the Feedback in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusCreated,
		Data:   newFeedback,
	})
}

func (f *FeedbackService) UpdateFeedbackHandler(w http.ResponseWriter, r *http.Request) error {
	userResponse := f.User.GetRequestUser(r)
	if userResponse == nil {
		log.Println("Error while fetching Request User")
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Unable to fetch user from token"},
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

	feedback, err := f.Store.GetFeedbackById(id)
	if err != nil {
		log.Println("Error in fetching the feedback", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "feedback not found"},
		})
	}

	var feedbackRequest models.UpdateFeedbackRequest

	json.NewDecoder(r.Body).Decode(&feedbackRequest)
	structErr := models.ValidateStruct(&feedbackRequest)
	if structErr != nil {
		log.Println("Error in validating the feedback struct", structErr)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   structErr,
		})
	}

	defer r.Body.Close()

	feedback = models.UpdateFeedback(feedback, &feedbackRequest)

	newFeedback, store_error := f.Store.UpdateFeedback(*feedback)
	if store_error != nil {
		msg := common.AnyToAnyStructField(store_error, &core.DatabaseError{})

		return core.APIResponse(w, &core.Response{
			Status: http.StatusInternalServerError,
			Data:   msg,
		})
	}

	redisErr := f.RedisClient.Set(newFeedback.Id.String(), newFeedback)
	if redisErr != nil {
		log.Println("Error in setting the feedback in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   newFeedback,
	})
}

func (f *FeedbackService) DeleteFeedbackHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		log.Println("Error in parsing the id", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: err.Error()},
		})
	}

	err = f.Store.DeleteFeedback(id)
	if err != nil {
		log.Println("Error in fetching the Feedback:", err)
		return core.APIResponse(w, &core.Response{
			Status: http.StatusBadRequest,
			Data:   &core.APIError{Detail: "Feedback not found"},
		})
	}

	redisErr := f.RedisClient.Del(id.String())
	if redisErr != nil {
		log.Println("Error in setting the Feedback in redis", redisErr)
	}

	return core.APIResponse(w, &core.Response{
		Status: http.StatusOK,
		Data:   map[string]string{"detail": "Feedback deleted successfully"},
	})
}
