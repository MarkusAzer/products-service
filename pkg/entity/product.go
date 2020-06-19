package entity

//Product data
type Product struct {
	ID          ID      `json:"id" bson:"_id,omitempty"`
	Version     Version `json:"version" bson:"_V,omitempty"`
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description,omitempty"`
	Slug        string  `json:"slug" bson:"slug,omitempty"`
	Image       string  `json:"image" bson:"image,omitempty"`
	Brand       string  `json:"brand" bson:"brand,omitempty"`
	Category    string  `json:"category" bson:"category,omitempty"`
	Price       int8    `json:"price" bson:"price"`
	Status      string  `json:"status" bson:"status"`
	CreatedAt   Time    `json:"createdAt" bson:"createdAt"`
}

//UpdateProduct data
type UpdateProduct struct {
	Version     Version `json:"version" bson:"_V,omitempty"`
	Name        string  `json:"name" bson:"name,omitempty" structs:",omitempty"`
	Description string  `json:"description" bson:"description,omitempty" structs:",omitempty"`
	Slug        string  `json:"slug" bson:"slug,omitempty" structs:",omitempty"`
	Image       string  `json:"image" bson:"image,omitempty" structs:",omitempty"`
	Brand       string  `json:"brand" bson:"brand,omitempty" structs:",omitempty"`
	Category    string  `json:"category" bson:"category,omitempty" structs:",omitempty"`
	Price       int8    `json:"price" bson:"price,omitempty" structs:",omitempty"`
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
