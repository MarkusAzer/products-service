package variant

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/go-playground/validator"
	"github.com/markus-azer/products-service/pkg/entity"
	"github.com/markus-azer/products-service/pkg/product"
)

//Service service interface
type Service struct {
	msgRepo     MessagesRepository
	storeRepo   StoreRepository
	productRepo product.StoreRepository
}

//NewService create new service
func NewService(msgR MessagesRepository, storeR StoreRepository, productR product.StoreRepository) *Service {
	return &Service{
		msgRepo:     msgR,
		storeRepo:   storeR,
		productRepo: productR,
	}
}

//CreateVariantDTO new variant DTO
type CreateVariantDTO struct {
	Product    entity.ID         `json:"product" validate:"required" structs:"product"`
	SKU        string            `json:"sku,omitempty" validate:"omitempty" structs:"sku,omitempty"`
	Quantity   int               `json:"quantity,omitempty" validate:"omitempty" structs:"quantity,omitempty"`
	Price      int               `json:"price,omitempty" validate:"omitempty,min=1" structs:"price,omitempty"`
	Image      string            `json:"image,omitempty" validate:"omitempty,uri" structs:"image,omitempty"`
	Attributes map[string]string `json:"attributes" validate:"required" structs:"attributes"`
}

//Create new variant
func (s *Service) Create(createVariantDTO CreateVariantDTO) (entity.ID, int32, []entity.ClientError) {

	//Validate DTOs, Terminate the Create process if the input is not valid
	if err := validator.New().Struct(createVariantDTO); err != nil {
		var errs []entity.ClientError

		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, entity.ClientError(e.Field()+" : "+fmt.Sprint(e)))
		}

		return "", 1, errs
	}

	ID := entity.NewID()
	Timestamp := time.Now()

	//Loop through the struct to generate events and validate
	var errs []entity.ClientError
	var messages []*entity.Message
	var version entity.Version = 1

	//Lower Case Attributes
	for k, v := range createVariantDTO.Attributes {
		delete(createVariantDTO.Attributes, k)
		createVariantDTO.Attributes[strings.ToLower(k)] = strings.ToLower(v)
	}

	payload := make(map[string]interface{})
	payload["product"] = createVariantDTO.Product
	payload["attributes"] = createVariantDTO.Attributes

	messages = append(messages, &entity.Message{
		ID:        string(ID),
		Type:      "PRODUCT_VARIANT_DRAFT_CREATED",
		Version:   version,
		Payload:   payload,
		Timestamp: Timestamp},
	)

	fields := reflect.TypeOf(createVariantDTO)
	values := reflect.ValueOf(createVariantDTO)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		switch field.Name {
		case "Product":
			_, err := s.productRepo.FindOneByID(entity.ID(value.String()))
			switch err {
			case entity.ErrNotFound:
				errs = append(errs, entity.ClientError("Product with ID "+value.String()+" doesnt Exist"))
			default:
				if err != nil {
					return "", 0, []entity.ClientError{"Internal Server Error"}
				}
			}

		case "Attributes":
			duplicatedVariant, err := s.storeRepo.FindOneByAttribute(createVariantDTO.Product, createVariantDTO.Attributes)
			if err != entity.ErrNotFound && err != nil {
				return "", 0, []entity.ClientError{"Internal Server Error"}
			}

			if duplicatedVariant != nil {
				errs = append(errs, entity.ClientError("Variant Attributes Duplication with ID "+string(duplicatedVariant.ID)))
			}

			if len(createVariantDTO.Attributes) > 3 {
				errs = append(errs, entity.ClientError("Max Attributes is 3"))
			}
		case "SKU":
			if value.String() != "" {
				version++

				payload := make(map[string]interface{})
				payload["sku"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_SKU_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Quantity":
			if value.Int() != 0 {
				version++

				payload := make(map[string]interface{})
				payload["quantity"] = value.Int()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_QUANTITY_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Price":
			if value.Int() != 0 {
				version++

				payload := make(map[string]interface{})
				payload["price"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_PRICE_UPDATED",
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
					Type:      "PRODUCT_VARIANT_IMAGE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		}
	}

	if len(errs) >= 1 {
		return "", 1, errs
	}

	v := &entity.Variant{
		ID:         ID,
		Version:    version,
		Product:    createVariantDTO.Product,
		SKU:        createVariantDTO.SKU,
		Quantity:   createVariantDTO.Quantity,
		Price:      createVariantDTO.Price,
		Image:      createVariantDTO.Image,
		Attributes: createVariantDTO.Attributes,
		CreatedAt:  Timestamp,
	}

	c := &entity.Command{AggregateID: string(ID), Type: "CreateVariant", Payload: structs.Map(createVariantDTO), Timestamp: Timestamp}
	_, err := s.storeRepo.StoreCommand(c)
	if err != nil {
		return "", 1, []entity.ClientError{"Internal Server Error"}
	}

	_, err = s.storeRepo.Create(v)
	if err != nil {
		return "", 1, []entity.ClientError{"Internal Server Error"}
	}

	s.msgRepo.SendMessages(messages)

	return ID, int32(v.Version), nil
}

//UpdateVariantDTO update variant DTO
type UpdateVariantDTO struct {
	SKU      string `json:"sku,omitempty" validate:"omitempty" structs:"sku,omitempty"`
	Quantity int    `json:"quantity,omitempty" validate:"omitempty" structs:"quantity,omitempty"`
	Price    int    `json:"price,omitempty" validate:"omitempty,min=1" structs:"price,omitempty"`
	Image    string `json:"image,omitempty" validate:"omitempty,uri" structs:"image,omitempty"`
}

//UpdateOne product
func (s *Service) UpdateOne(ID entity.ID, v int32, updateVariantDTO UpdateVariantDTO) (int32, []entity.ClientError) {
	//Validate DTOs, Terminate the Create process if the input is not valid
	if err := validator.New().Struct(updateVariantDTO); err != nil {
		var errs []entity.ClientError

		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, entity.ClientError(e.Field()+" : "+fmt.Sprint(e)))
		}

		return 0, errs
	}

	Timestamp := time.Now()
	version := entity.Version(v)

	variant, err := s.storeRepo.FindOneByID(ID)
	switch err {
	case entity.ErrNotFound:
		e := entity.ClientError("Variant with id " + string(ID) + " Not found")
		return 0, []entity.ClientError{e}
	default:
		if err != nil {
			return 0, []entity.ClientError{"Internal Server Error"}
		}
	}

	if version != variant.Version {
		return 0, []entity.ClientError{"Version conflict"}
	}

	//Loop through the struct to generate events and validate
	var errs []entity.ClientError
	var messages []*entity.Message

	fields := reflect.TypeOf(updateVariantDTO)
	values := reflect.ValueOf(updateVariantDTO)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		switch field.Name {
		case "SKU":
			if value.String() != "" {
				if variant.SKU == updateVariantDTO.SKU {
					errs = append(errs, "Sku already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["sku"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_SKU_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Quantity":
			if value.String() != "" {
				if variant.Quantity == updateVariantDTO.Quantity {
					errs = append(errs, "Quantity already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["quantity"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_QUANTITY_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Image":
			//TODO: check if image exist and add event to delete other image
			if value.String() != "" {
				if variant.Image == updateVariantDTO.Image {
					errs = append(errs, "Image already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["image"] = value.String()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_IMAGE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		case "Price":
			if value.Int() != 0 {
				if variant.Price == updateVariantDTO.Price {
					errs = append(errs, "Price already updated")
				}
				version++

				payload := make(map[string]interface{})
				payload["price"] = value.Int()

				messages = append(messages, &entity.Message{
					ID:        string(ID),
					Type:      "PRODUCT_VARIANT_PRICE_UPDATED",
					Version:   version,
					Payload:   payload,
					Timestamp: Timestamp})
			}
		}
	}

	if len(errs) > 0 {
		return 0, errs
	}

	if version == variant.Version {
		return 0, []entity.ClientError{"No Updates Found"}
	}

	c := &entity.Command{AggregateID: string(ID), Type: "UpdateProduct", Payload: structs.Map(updateVariantDTO), Timestamp: Timestamp}
	_, err = s.storeRepo.StoreCommand(c)
	if err != nil {
		return 0, []entity.ClientError{"Internal Server Error"}
	}

	up := &entity.UpdateVariant{
		Version:  version,
		SKU:      updateVariantDTO.SKU,
		Quantity: updateVariantDTO.Quantity,
		Price:    updateVariantDTO.Price,
		Image:    updateVariantDTO.Image,
	}

	updatedNum, err := s.storeRepo.UpdateOne(ID, up, entity.Version(v))
	if err != nil {
		return 0, []entity.ClientError{"Internal Server Error"}
	}

	if updatedNum != 1 {
		return 0, []entity.ClientError{"Version conflict"}
	}

	//TODO:handle failure cases
	s.msgRepo.SendMessages(messages)

	return int32(version), nil
}

//Delete product
func (s *Service) Delete(id entity.ID, v int32) *entity.ClientError {
	Timestamp := time.Now()
	version := entity.Version(v)

	//TODO check Transactions
	p, err := s.storeRepo.FindOneByID(id)
	switch err {
	case entity.ErrNotFound:
		e := entity.ClientError("Variant with id " + string(id) + " Not found")
		return &e
	default:
		if err != nil {
			e := entity.ClientError("Internal Server Error")
			return &e
		}
	}

	if version != p.Version {
		e := entity.ClientError("miss matching version")
		return &e
	}

	c := &entity.Command{AggregateID: string(id), Type: "DeleteVariant", Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	updatedNum, err := s.storeRepo.DeleteOne(id, version)
	if err != nil {
		e := entity.ClientError("Internal Server Error")
		return &e
	}

	if updatedNum != 1 {
		e := entity.ClientError("Version conflict")
		return &e
	}

	version++

	//TODO:handle failure cases
	m := &entity.Message{ID: string(id), Type: "PRODUCT_VARIANT_DELETED", Version: version, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return nil
}
