package entity

import "time"

//Message message struct
type Message struct {

	// ID contains the AggregateID
	ID string

	// Type contains message type
	Type string

	// Version version count
	Version Version

	// Payload contains message data
	Payload map[string]interface{}

	// Timestamp
	Timestamp time.Time
}
