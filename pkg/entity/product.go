package entity

import "time"

//Product data
type Product struct {
	ID          ID        `json:"id" bson:"_id"`
	Version     Version   `json:"version" bson:"_V"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	Slug        string    `json:"slug,omitempty" bson:"slug,omitempty"`
	Location    string    `json:"location,omitempty" bson:"location,omitempty"`
	Image       string    `json:"image,omitempty" bson:"image,omitempty"`
	Brand       string    `json:"brand,omitempty" bson:"brand,omitempty"`
	Category    string    `json:"category,omitempty" bson:"category,omitempty"`
	Price       int8      `json:"price,omitempty" bson:"price,omitempty"`
	Status      string    `json:"status,omitempty" bson:"status,omitempty"`
	Seller      string    `json:"seller,omitempty" bson:"seller,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

//UpdateProduct data
type UpdateProduct struct {
	Version     Version `bson:"_V,omitempty"`
	Name        string  `bson:"name,omitempty" structs:",omitempty"`
	Description string  `bson:"description,omitempty" structs:",omitempty"`
	Slug        string  `bson:"slug,omitempty" structs:",omitempty"`
	Location    string  `json:"location,omitempty" bson:"location,omitempty"`
	Image       string  `bson:"image,omitempty" structs:",omitempty"`
	Brand       string  `bson:"brand,omitempty" structs:",omitempty"`
	Category    string  `bson:"category,omitempty" structs:",omitempty"`
	Price       int8    `bson:"price,omitempty" structs:",omitempty"`
	Status      string  `bson:"status,omitempty" structs:",omitempty"`
}

//Validate Validate Product Struct
// TODO: better validation https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
func (p *Product) Validate() []string {

	var errs []string

	if p.Name == "" {
		errs = append(errs, "Name : Name is required")
	}

	if (p.Price == 0) || (p.Price < 1) {
		errs = append(errs, "Price :- Provide Valid Price")
	}

	return errs
}
