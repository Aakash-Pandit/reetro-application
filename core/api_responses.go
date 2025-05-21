package core

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type Response struct {
	Data   any
	Status int
}

type ListAPIResponseBody struct {
	Result any
	Count  int
}

type ListAPI struct {
	Status int
	Result *ListAPIResponseBody
}

type DatabaseError struct {
	Detail  string `json:"Detail"`
	Message string `json:"Message"`
}

type APIError struct {
	Detail string `json:"detail"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func APIResponse(w http.ResponseWriter, response *Response) error {
	w.WriteHeader(response.Status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response.Data)
	return err
}

func ListAPIResponse(w http.ResponseWriter, response *ListAPI) error {
	w.WriteHeader(response.Status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response.Result)
	return err
}

func HTTPHandleFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		if err := fn(w, r); err != nil {
			log.Println("Error in handling the request", err)
			APIResponse(w, &Response{
				Status: http.StatusBadRequest,
				Data:   &APIError{Detail: err.Error()},
			})
		}
	}
}
