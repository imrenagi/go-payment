package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/imrenagi/go-payment"
)

type Meta struct {
	TotalItems  int    `json:"total_items"`
	TotalPages  int    `json:"total_pages"`
	CurrentPage int    `json:"cur_page"`
	Cursor      string `json:"last_cursor"`
}

// Error is struct used to return error message to the client
type Error struct {
	StatusCode int    `json:"error_code"`
	Message    string `json:"error_message"`
}

// Empty used to return nothing
type Empty struct{}

// WriteSuccessResponse creates success response for the http handler
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, data interface{}, headMap map[string]string) {
	w.Header().Add("Content-Type", "application/json")
	if headMap != nil && len(headMap) > 0 {
		for key, val := range headMap {
			w.Header().Add(key, val)
		}
	}
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

// WriteFailResponse creates error response for the http handler
func WriteFailResponse(w http.ResponseWriter, statusCode int, error interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(error)
	w.Write(jsonData)
}

// WriteFailResponseFromError creates error response based on the given error
func WriteFailResponseFromError(w http.ResponseWriter, err error) {

	var statusCode int
	if errors.Is(err, payment.ErrNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, payment.ErrInternal) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, payment.ErrDatabase) {
		statusCode = http.StatusInternalServerError
	} else if errors.Is(err, payment.ErrBadRequest) {
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, payment.ErrCantProceed) {
		statusCode = http.StatusUnprocessableEntity
	} else if errors.Is(err, payment.ErrUnauthorized) {
		statusCode = http.StatusUnauthorized
	} else if errors.Is(err, payment.ErrForbidden) {
		statusCode = http.StatusForbidden
	} else {
		statusCode = http.StatusInternalServerError
	}

	errorMsg := Error{
		Message:    err.Error(),
		StatusCode: statusCode,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(errorMsg.StatusCode)
	jsonData, _ := json.Marshal(errorMsg)
	w.Write(jsonData)
}
