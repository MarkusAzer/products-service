package entity

import (
	"errors"
	"fmt"
)

//CustomError CustomError
type CustomError struct {
	Status int
	Err    error
}

//NewCustomError create new custom Error
func NewCustomError(err string, status int) *CustomError {
	return &CustomError{
		Status: status,
		Err:    errors.New(err),
	}
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("status: %d , err: %v", c.Status, c.Err)
}

//ErrNotFound not found
var ErrNotFound = errors.New("Not found")
