package entity

import "time"

//Command command struct
type Command struct {

	// AggregateID contains the AggregateID
	AggregateID string `json:"aggregateId" bson:"aggregateId"`

	// Type contains message type
	Type string `json:"type" bson:"type"`

	// Payload contains message data
	Payload map[string]interface{} `json:"payload" bson:"payload"`

	// Timestamp
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
