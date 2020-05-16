package entity

import (
	"time"

	"github.com/google/uuid"
)

//ID Type
type ID string

//Time Type
type Time time.Time

//NewID create a new id
func NewID() ID {
	id := uuid.New()
	return ID(id.String())
}

//TimeNow create Time Now
func TimeNow() Time {
	return Time(time.Now().UTC())
}
