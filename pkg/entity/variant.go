package entity

import "time"

//Variant Variant
type Variant struct {
	ID         ID                `json:"id" bson:"_id"`
	Product    ID                `json:"product" bson:"product"`
	Version    Version           `json:"version" bson:"_V"`
	SKU        string            `json:"sku,omitempty" bson:"sku,omitempty"`
	Quantity   int               `json:"quantity" bson:"quantity"`
	Price      int               `json:"price" bson:"price"`
	Image      string            `json:"image,omitempty" bson:"image,omitempty"`
	Attributes map[string]string `json:"attributes" bson:"attributes"`
	CreatedAt  time.Time         `json:"createdAt" bson:"createdAt"`
}

//UpdateVariant data
type UpdateVariant struct {
	Version  Version `bson:"_V,omitempty" structs:",omitempty"`
	SKU      string  `bson:"sku,omitempty" structs:",omitempty"`
	Quantity int     `bson:"quantity,omitempty" structs:",omitempty"`
	Price    int     `bson:"price,omitempty" structs:",omitempty"`
	Image    string  `bson:"image,omitempty" structs:",omitempty"`
}
