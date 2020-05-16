package entity

//Category data
type Category struct {
	ID          ID     `json:"id" bson:"_id,omitempty"`
	Version     int32  `json:"version" bson:"_V,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description,omitempty"`
	Slug        string `json:"slug" bson:"slug,omitempty"`
	CreatedAt   Time   `json:"createdAt" bson:"createdAt"`
}
