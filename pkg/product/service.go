package product

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator"
	"github.com/markus-azer/products-service/pkg/brand"
	"github.com/markus-azer/products-service/pkg/entity"
)

// Validation in Application layer
// In this layer, as validation, we must ensure that domain objects can receive the input. We should reject the input which the domain object can't be received.

// For example, when some mandatory parameters are missing, it should be rejected because the domain object has no way to receive like that parameter.

//DTOhttps://softwareengineering.stackexchange.com/questions/373284/what-is-the-use-of-dto-instead-of-entity
//https://stackoverflow.com/questions/21554977/should-services-always-return-dtos-or-can-they-also-return-domain-models

//clean https://www.entropywins.wtf/blog/2016/11/24/implementing-the-clean-architecture/
//https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/domain-model-layer-validations
//http://www.plainionist.net/Implementing-Clean-Architecture-Controller-Presenter/

//TODO struct to map funchttps://stackoverflow.com/questions/23589564/function-for-converting-a-struct-to-map-in-golang

//Service service interface
type Service struct {
	msgRepo   MessagesRepository
	storeRepo StoreRepository
	brandRepo brand.StoreRepository
}

//NewService create new service
func NewService(msgR MessagesRepository, storeR StoreRepository, brandR brand.StoreRepository) *Service {
	return &Service{
		msgRepo:   msgR,
		storeRepo: storeR,
		brandRepo: brandR,
	}
}

//CreateProductDTO new product DTO
type CreateProductDTO struct {
	Name        string `json:"name" validate:"required,min=3" structs:"name,omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty,min=20" structs:"description,omitempty"`
	Slug        string `json:"slug,omitempty" validate:"omitempty" structs:"slug,omitempty"` //TODO: find slug validation
	Location    string `json:"location,omitempty" validate:"omitempty" structs:"location,omitempty"`
	Image       string `json:"image,omitempty" validate:"omitempty,uri" structs:"image,omitempty"`
	Brand       string `json:"brand,omitempty" validate:"omitempty" structs:"brand,omitempty"`
	Category    string `json:"category,omitempty" validate:"omitempty" structs:"category,omitempty"`
	Price       int8   `json:"price,omitempty" validate:"omitempty,min=1" structs:"price,omitempty"`
	Seller      string `json:"seller,omitempty" validate:"required" structs:"seller,omitempty"`
}

//Create new product
func (s *Service) Create(createProductDTO CreateProductDTO) (entity.ID, entity.Version, []string) {

	//Validate DTOs, Terminate the Create process if the input is not valid
	if err := validator.New().Struct(createProductDTO); err != nil {
		var errs []string

		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, e.Field()+" : "+fmt.Sprint(e))
		}

		return "", 1, errs
	}

	ID := entity.NewID()
	Timestamp := time.Now()

	//Loop through the struct to generate events and validate
	var errs []string
	var messages []*entity.Message
	var version entity.Version = 1

	payload := make(map[string]interface{})
	payload["seller"] = createProductDTO.Seller
	messages = append(messages, &entity.Message{
		ID:        string(ID),
		Type:      "PRODUCT_DRAFT_CREATED",
		Version:   version,
		Payload:   payload,
		Timestamp: Timestamp},
	)

	fields := reflect.TypeOf(createProductDTO)
	values := reflect.ValueOf(createProductDTO)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		switch field.Name {
		case "Name":
			version++

			payload := make(map[string]interface{})
			payload["name"] = value.String()

			messages = append(messages, &entity.Message{
				ID:        string(ID),
				Type:      "PRODUCT_NAME_UPDATED",
				Version:   version,
				Payload:   payload,
				Timestamp: Timestamp})

		case "Description":
			if value.String() != "" {
				version++

				payload := make(map[string]interface{})
				payload["description"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_DESCRIPTION_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Slug":
			//TODO: check uniqueness
			if value.String() != "" {
				version++

				payload := make(map[string]interface{})
				payload["slug"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_SLUG_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Location":
			if value.String() != "" {
				version++

				payload := make(map[string]interface{})
				payload["location"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_Location_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Image":
			//TODO: check if image exist and add event to delete other image
			if value.String() != "" {
				version++

				payload := make(map[string]interface{})
				payload["image"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_IMAGE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Brand":
			if value.String() != "" {
				version++

				_, err := s.brandRepo.FindOneByName(value.String())
				switch err {
				case entity.ErrNotFound:
					errs = append(errs, "Brand "+value.String()+" Not found")
				default:
					if err != nil {
						return "", 1, []string{"Internal Server Error"}
					}
				}

				payload := make(map[string]interface{})
				payload["brand"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_BRAND_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Category":
			if value.String() != "" {
				version++

				payload := make(map[string]interface{})
				payload["category"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_CATEGORY_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Price":
			if value.Int() != 0 {
				version++

				payload := make(map[string]interface{})
				payload["price"] = value.Int()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_PRICE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		}
	}

	if len(errs) >= 1 {
		return "", 1, errs
	}

	p := &entity.Product{
		ID:          ID,
		Version:     version,
		Name:        createProductDTO.Name,
		Description: createProductDTO.Description,
		Slug:        createProductDTO.Slug,
		Image:       createProductDTO.Image,
		Brand:       createProductDTO.Brand,
		Category:    createProductDTO.Category,
		Price:       createProductDTO.Price,
		Status:      "unpublish", // Init the product as unpublish
		CreatedAt:   Timestamp,
	}

	// data, err := json.Marshal(p)

	// if err != nil {
	// 	return "", 1, []string{"Internal Server Error"}
	// }

	// var newMap map[string]interface{}
	// err = json.Unmarshal(data, &newMap) // Convert to a map

	c := &entity.Command{AggregateID: string(ID), Type: "CreateProduct", Payload: structs.Map(createProductDTO), Timestamp: Timestamp}
	_, err := s.storeRepo.StoreCommand(c)
	if err != nil {
		return "", 1, []string{"Internal Server Error"}
	}

	_, err = s.storeRepo.Create(p)
	if err != nil {
		return "", 1, []string{"Internal Server Error"}
	}

	s.msgRepo.SendMessages(messages)

	return ID, p.Version, nil
}

//UpdateProductDTO new product DTO
type UpdateProductDTO struct {
	// Version		entity.Version `json:"_V,omitempty" validate:"omitempty,required,min=3"`
	Name        string `json:"name,omitempty" validate:"omitempty,min=3" structs:"name,omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty,min=20" structs:"description,omitempty"`
	Slug        string `json:"slug,omitempty" validate:"omitempty" structs:"slug,omitempty"` //TODO: find slug validation
	Location    string `json:"location,omitempty" bson:"location,omitempty"`
	Image       string `json:"image,omitempty" validate:"omitempty,uri" structs:"image,omitempty"`
	Brand       string `json:"brand,omitempty" validate:"omitempty" structs:"brand,omitempty"`
	Category    string `json:"category,omitempty" validate:"omitempty" structs:"category,omitempty"`
	Status      string `json:"status,omitempty" validate:"omitempty,oneof=publish unpublish" structs:"status,omitempty"`
	Price       int8   `json:"price,omitempty" validate:"omitempty" structs:"price,omitempty"`
}

//UpdateOne product
func (s *Service) UpdateOne(ID entity.ID, v int32, updateProductDTO UpdateProductDTO) (int32, []string) {
	//Validate DTOs, Terminate the Create process if the input is not valid
	if err := validator.New().Struct(updateProductDTO); err != nil {
		var errs []string

		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, e.Field()+" : "+fmt.Sprint(e))
		}

		return 0, errs
	}

	Timestamp := time.Now()
	version := entity.Version(v)

	p, err := s.storeRepo.FindOneByID(ID)
	switch err {
	case entity.ErrNotFound:
		return 0, []string{"Product with id " + string(ID) + " Not found"}
	default:
		if err != nil {
			return 0, []string{"Internal Server Error"}
		}
	}

	if version != p.Version {
		return 0, []string{"Version conflict"}
	}

	//Loop through the struct to generate events and validate
	var errs []string
	var messages []*entity.Message

	fields := reflect.TypeOf(updateProductDTO)
	values := reflect.ValueOf(updateProductDTO)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		switch field.Name {
		case "Name":
			if value.String() != "" {
				if p.Name == updateProductDTO.Name {
					errs = append(errs, "Name already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["name"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_NAME_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Description":
			if value.String() != "" {
				if p.Description == updateProductDTO.Description {
					errs = append(errs, "Description already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["description"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_DESCRIPTION_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Slug":
			//TODO: check uniqueness
			if value.String() != "" {
				if p.Slug == updateProductDTO.Slug {
					errs = append(errs, "Slug already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["slug"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_SLUG_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Image":
			//TODO: check if image exist and add event to delete other image
			if value.String() != "" {
				if p.Image == updateProductDTO.Image {
					errs = append(errs, "Image already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["image"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_IMAGE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Brand":
			if value.String() != "" {
				if p.Brand == updateProductDTO.Brand {
					errs = append(errs, "Brand already updated")
				}
				version++

				_, err := s.brandRepo.FindOneByName(value.String())
				switch err {
				case entity.ErrNotFound:
					errs = append(errs, "Brand "+value.String()+" Not found")
				default:
					if err != nil {
						return 0, []string{"Internal Server Error"}
					}
				}

				payload := make(map[string]interface{})
				payload["brand"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_BRAND_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Category":
			if value.String() != "" {
				if p.Category == updateProductDTO.Category {
					errs = append(errs, "Category already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["category"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_CATEGORY_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Status":
			if value.String() != "" {
				if p.Status == updateProductDTO.Status {
					errs = append(errs, "Status already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["status"] = value.String()

				if updateProductDTO.Status == "publish" {
					messages = append(messages, &entity.Message{
						ID:        string(ID),
						Type:      "PRODUCT_PUBLISHED",
						Version:   version,
						Payload:   payload,
						Timestamp: Timestamp})

				} else {
					messages = append(messages, &entity.Message{
						ID:        string(ID),
						Type:      "PRODUCT_UNPUBLISHED",
						Version:   version,
						Payload:   payload,
						Timestamp: Timestamp})
				}
			}
		case "Price":
			if value.Int() != 0 {
				if p.Price == updateProductDTO.Price {
					errs = append(errs, "Price already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["price"] = value.Int()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_PRICE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		}
	}

	if len(errs) > 0 {
		return 0, errs
	}

	if version == p.Version {
		return 0, []string{"No Updates Found"}
	}

	c := &entity.Command{AggregateID: string(ID), Type: "UpdateProduct", Payload: structs.Map(updateProductDTO), Timestamp: Timestamp}
	_, err = s.storeRepo.StoreCommand(c)
	if err != nil {
		return 0, []string{"Internal Server Error"}
	}

	up := &entity.UpdateProduct{
		Version:     version,
		Name:        updateProductDTO.Name,
		Description: updateProductDTO.Description,
		Slug:        updateProductDTO.Slug,
		Image:       updateProductDTO.Image,
		Brand:       updateProductDTO.Brand,
		Category:    updateProductDTO.Category,
		Status:      updateProductDTO.Status,
		Price:       updateProductDTO.Price,
	}

	updatedNum, err := s.storeRepo.UpdateOneP(ID, up, entity.Version(v))
	if err != nil {
		return 0, []string{"Internal Server Error"}
	}

	if updatedNum != 1 {
		return 0, []string{"Version conflict"}
	}

	//TODO:handle failure cases
	s.msgRepo.SendMessages(messages)

	return int32(version), nil
}

//Publish publish product //** DEPRECATED **//
func (s *Service) Publish(ID entity.ID, v int32) (int32, []string) {
	Timestamp := time.Now()
	version := entity.Version(v)

	p, err := s.storeRepo.FindOneByID(ID)
	switch err {
	case entity.ErrNotFound:
		return 0, []string{"Product with id " + string(ID) + " Not found"}
	default:
		if err != nil {
			return 0, []string{"Internal Server Error"}
		}
	}

	if version != p.Version {
		return 0, []string{"Version conflict"}
	}

	if p.Status == "Publish" {
		return 0, []string{"Product is already published"}
	}

	newMap := make(map[string]interface{})
	newMap["Status"] = "Publish"

	c := &entity.Command{AggregateID: string(ID), Type: "PublishProduct", Payload: newMap, Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	version++

	//TODO Patch the update
	p.Status = "Publish"
	p.Version = version
	s.storeRepo.UpdateOne(ID, p, version)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_PUBLISHED", Version: version, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return int32(version), nil
}

//Unpublish unpublish product //** DEPRECATED **//
func (s *Service) Unpublish(ID entity.ID, v int32) (int32, error) {
	Timestamp := time.Now()
	version := entity.Version(v)

	//TODO check Transactions
	p, err := s.storeRepo.FindOneByID(ID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if version != p.Version {
		fmt.Println("miss matching version")
		return 0, errors.New("miss matching version")
	}

	//TODO: create status map
	if p.Status == "Unpublish" {
		fmt.Println("Product is already Unpublish")
		return 0, errors.New("Product is already Unpublish")
	}

	newMap := make(map[string]interface{})
	newMap["Status"] = "Unpublish"

	c := &entity.Command{AggregateID: string(ID), Type: "UnpublishProduct", Payload: newMap, Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	version++

	//TODO Patch the update
	p.Status = "Unpublish"
	p.Version = version
	s.storeRepo.UpdateOne(ID, p, version)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_UNPUBLISHED", Version: version, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return int32(version), nil
}

//UpdatePrice product price //** DEPRECATED **//
func (s *Service) UpdatePrice(ID entity.ID, v int32, price int) (int32, error) {
	Timestamp := time.Now()
	version := entity.Version(v)

	//TODO check Transactions
	p, err := s.storeRepo.FindOneByID(ID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if version != p.Version {
		fmt.Println("miss matching version")
		return 0, errors.New("miss matching version")
	}

	newMap := make(map[string]interface{})
	newMap["Price"] = price

	c := &entity.Command{AggregateID: string(ID), Type: "UpdateProductPrice", Payload: newMap, Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	version++

	//TODO Patch the update
	p.Price = int8(price)
	p.Version = version
	s.storeRepo.UpdateOne(ID, p, version)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_PRICE_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return int32(version), nil
}

//Delete product
func (s *Service) Delete(ID entity.ID, v int32) error {
	Timestamp := time.Now()
	version := entity.Version(v)

	//TODO check Transactions
	p, err := s.storeRepo.FindOneByID(ID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if version != p.Version {
		fmt.Println("miss matching version")
		return errors.New("miss matching version")
	}

	c := &entity.Command{AggregateID: string(ID), Type: "DeleteProduct", Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	s.storeRepo.DeleteOne(ID, version)

	version++

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_DELETED", Version: version, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return nil
}
