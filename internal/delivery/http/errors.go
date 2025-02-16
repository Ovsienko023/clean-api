package http

import (
	"api/internal/domain/errdomain"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorResponseDetails struct {
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
	Position    string `json:"position,omitempty"`
}

type ErrorResponseError struct {
	Code        int                    `json:"code,omitempty"`
	Description string                 `json:"description,omitempty"`
	Details     []ErrorResponseDetails `json:"details"`
}

type ErrorResponse struct {
	Error ErrorResponseError `json:"error,omitempty"`
}

func (r *ErrorResponse) Add(reason, description, position string) {
	details := ErrorResponseDetails{
		Reason:      reason,
		Description: description,
		Position:    position,
	}
	r.Error.Details = append(r.Error.Details, details)
}

func (r *ErrorResponse) Done(w http.ResponseWriter, err error) {
	var validateErr validator.ValidationErrors

	ok := errors.As(err, &validateErr)
	if ok {
		r.Error.Code = http.StatusBadRequest
	} else {

		if code, ok := errToCode[err]; ok {
			r.Error.Code = code
		} else {
			r.Error.Code = http.StatusInternalServerError
		}

	}

	r.Error.Description = err.Error()

	if len(r.Error.Details) == 0 {
		r.Error.Details = []ErrorResponseDetails{}
	}

	response, _ := json.Marshal(r)

	w.WriteHeader(r.Error.Code)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(response)
}

var errToCode = map[error]int{
	errdomain.ErrObjectNotFound: http.StatusNotFound,
}

func getStatusCode(err error) int {
	code, ok := errToCode[err]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}
