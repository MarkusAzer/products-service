package handler

import "github.com/markus-azer/products-service/pkg/entity"

// response struct which contains an API Response
type response struct {
	StatusCode int                    `json:"statusCode,omitempty"` //todo enum and hide Response
	Message    string                 `json:"message,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Errors     []entity.ErrorField    `json:"errors,omitempty"`
	Successful bool                   `json:"successful"`
}
