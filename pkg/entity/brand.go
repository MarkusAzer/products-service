package entity

import "time"

//Brand data
type Brand struct {
	ID          ID        `json:"id" bson:"_id,omitempty"`
	Version     Version   `json:"version" bson:"_V,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description,omitempty"`
	Slug        string    `json:"slug" bson:"slug,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}
