package response

import (
	"encoding/json"
	"net/http"
)

type ValidResponse[T any] struct {
	Message string `json:"message"`
	Data  	T	   `json:"data"`
}

type ErrorResponse struct {
	Message	  string `json:"message"`
	ErrorCode int `json:"error_code"`
}

func SendValidResponse[T any](w http.ResponseWriter, message string, data T) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ValidResponse[T]{
		Message: message,
		Data:    data,
	})
}

func SendErrorResponse(w http.ResponseWriter, message string, errorCode int) {
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message:   message,
		ErrorCode: errorCode,
	})
}
