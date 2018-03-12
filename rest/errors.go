package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error defines HTTP API's error type
type Error struct {
	Message   string      `json:"message"`
	ErrorCode int         `json:"errorCode"`
	Status    int         `json:"-"`
	Details   interface{} `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: %s. Code: %d", e.Message, e.ErrorCode)
}

// Send sends the chosen error as JSON response
func (e *Error) Send(w http.ResponseWriter) (int, error) {
	m, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	return w.Write(m)
}

// ErrNotFound is used for 404 responses
var ErrNotFound = Error{
	Message:   "Not found",
	ErrorCode: 1000,
	Status:    http.StatusNotFound,
}

// ErrMethodNotAllowed is used for 405 responses
var ErrMethodNotAllowed = Error{
	Message:   "Method not allowed",
	ErrorCode: 1001,
	Status:    http.StatusMethodNotAllowed,
}

// ErrUnsupportedMediaType is used when incoming request does not match Content-Type header
var ErrUnsupportedMediaType = Error{
	Message:   "Content-Type should be application/json",
	ErrorCode: 1002,
	Status:    http.StatusUnsupportedMediaType,
}

// ErrInvalidCredentials is used when requester does not have access to an API resource
var ErrInvalidCredentials = Error{
	Message:   "Invalid credentials",
	ErrorCode: 1003,
	Status:    http.StatusUnauthorized,
}

// ErrMissingInput is used when a required input field is missing in an API request
var ErrMissingInput = Error{
	Message:   "Missing required input",
	ErrorCode: 1004,
	Status:    http.StatusBadRequest,
}

// ErrInvalidInput is used when request contains proper fields but invalid data
var ErrInvalidInput = Error{
	Message:   "Invalid input",
	ErrorCode: 1005,
	Status:    http.StatusBadRequest,
}

// ErrCreateResource is used when a resource cannot be created. This is used for any resources being created.
var ErrCreateResource = Error{
	Message:   "Unable to create resource",
	ErrorCode: 1006,
	Status:    http.StatusBadRequest,
}

// ErrUpdateResource is used when a resource cannot be updated
var ErrUpdateResource = Error{
	Message:   "Unable to update resource",
	ErrorCode: 1007,
	Status:    http.StatusBadRequest,
}
