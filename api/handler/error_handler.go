package handler

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/markus-azer/products-service/pkg/entity"
)

func getStatusCode(kind entity.Kind) int {
	ErrHTTPStatusMap := map[entity.Kind]int{
		entity.NotFound:               http.StatusBadRequest,
		entity.ValidationFailed:       http.StatusBadRequest,
		entity.ConcurrentModification: http.StatusConflict,
		entity.Unexpected:             http.StatusInternalServerError,
		entity.NoUpdates:              http.StatusBadRequest,
	}
	code := ErrHTTPStatusMap[kind]

	// If error code is not found then apply a default case
	if code == 0 {
		code = http.StatusInternalServerError
	}

	return code
}
func errorHandler(err error) *response {

	e, ok := err.(*entity.Error)

	if !ok {
		return &response{StatusCode: 500, Message: "Internal Service Error", Successful: false}
	}

	StatusCode := getStatusCode(e.Kind)

	return &response{StatusCode: StatusCode, Message: string(e.ErrorMessage), Errors: e.Errors, Successful: false}
}

func serializationErrorHandler(err error) *response {
	switch {
	case err == io.EOF:
		return &response{StatusCode: 400, Message: "Provide valid Body", Successful: false}
	case err != nil && strings.Contains(err.Error(), "json: unknown field"):
		m := regexp.MustCompile(`\"(.*)\"`)
		field := m.FindString(err.Error())
		errors := []entity.ErrorField{
			{
				Field: field,
				Error: "Not allowed",
			},
		}
		return &response{StatusCode: http.StatusBadRequest, Message: "Provide valid Payload", Errors: errors, Successful: false}
	default:
		return &response{StatusCode: 400, Message: "Provide valid Body", Successful: false}
	}
}
