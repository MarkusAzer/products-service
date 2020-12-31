package entity

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
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

//ClientError ClientError
type ClientError string

// Op A unique string operation pointing to a function
// Multiple operations can construct a friendly stack trace.
type Op string

// Kind category of the error
type Kind int

const (
	// NotFound NotFound
	NotFound Kind = iota + 1
	// ValidationFailed ValidationFailed
	ValidationFailed
	// ConcurrentModification ConcurrentModification
	ConcurrentModification
	// Unexpected Unexpected
	Unexpected
	//NoUpdates no updates found
	NoUpdates
)

//ErrorMessage ErrorMessage
type ErrorMessage string

// ErrorField ErrorField
type ErrorField struct {
	Field string
	Error string
}

// Error error
type Error struct {
	Op
	Kind
	ErrorMessage
	Severity logrus.Level
	Errors   []ErrorField
	Err      error
	// ...application specific fields
}

//Error err
func (e *Error) Error() string {
	return e.Err.Error()
}

//E e
func E(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case error:
			e.Err = arg
		case logrus.Level:
			e.Severity = arg
		case []ErrorField:
			e.Errors = arg
		case ErrorMessage:
			e.ErrorMessage = arg
		default:
			panic("bad call to E")
		}
	}

	return e
}

//Ops ops
func Ops(e *Error) []Op {
	res := []Op{e.Op}

	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)
	return res
}

//KindE kind
func KindE(err error) Kind {
	e, ok := err.(*Error)
	if !ok {
		return Unexpected
	}

	if e.Kind != 0 {
		return e.Kind
	}

	return KindE(e.Err)
}

// func SystemErr(err error) {
// 	sysErr, ok := err.(*Error)
// 	if !ok{
// 		logrus.Error(err)
// 		return
// 	}

// 	entry := logrus.WithFields(
// 		"operations", Ops(sysErr),
// 		"kind", Kind(sysErr),
// 		// application specific data
// 	)

// 	switch Level(err){
// 	case Warning:
// 		entry.Warnf("%s: %v", sysErr.Op, err)
// 	case Info:
// 		entry.Infof("%s: %v", sysErr.Op, err)
// 	case Debug:
// 		entry.Debugf("%s: %v", sysErr.Op, err)
// 	default:
// 		entry.Errorf("%s: %v", sysErr.Op, err)
// 	}
// }
