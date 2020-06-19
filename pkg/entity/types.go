package entity

import (
	"time"

	"github.com/google/uuid"
)

//ID Type
type ID string

//Time Type
type Time time.Time

//Version Type
type Version int32

//NewID create a new id
func NewID() ID {
	id := uuid.New()
	return ID(id.String())
}

//TimeNow create Time Now
func TimeNow() Time {
	return Time(time.Now().UTC())
}
