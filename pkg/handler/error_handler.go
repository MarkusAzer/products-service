package handler

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/markus-azer/products-service/pkg/entity"
)

func getStatusCode(kind entity.Kind) int {
	switch kind {
	case entity.NotFound:
		return http.StatusBadRequest
	case entity.ValidationFailed:
		return http.StatusBadRequest
	case entity.ConcurrentModification:
		return http.StatusConflict
	case entity.Unexpected:
		return http.StatusInternalServerError
	case entity.NoUpdates:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
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
			entity.ErrorField{
				Field: field,
				Error: "Not allowed",
			},
		}
		return &response{StatusCode: http.StatusBadRequest, Message: "Provide valid Payload", Errors: errors, Successful: false}
	default:
		return &response{StatusCode: 400, Message: "Provide valid Body", Successful: false}
	}
}
