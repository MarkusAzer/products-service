package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ID Type
type ID string

//Time Type
type Time time.Time

//NewID create a new id
func NewID() ID {
	return ID(primitive.NewObjectID().Hex())
}

//TimeNow create Time Now
func TimeNow() Time {
	return Time(time.Now().UTC())
}
