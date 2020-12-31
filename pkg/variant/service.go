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
	"github.com/sirupsen/logrus"
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
func (s *Service) Create(createVariantDTO CreateVariantDTO) (*entity.ID, *int32, *entity.Error) {

	//Validate DTOs, Terminate the Create process if the input is not valid
	if err := validator.New().Struct(createVariantDTO); err != nil {
		errs := entity.Error{Op: "Create", Kind: entity.ValidationFailed, ErrorMessage: "Provide valid Payload", Severity: logrus.InfoLevel}

		for _, e := range err.(validator.ValidationErrors) {
			errs.Errors = append(errs.Errors, entity.ErrorField{Field: e.Field(), Error: fmt.Sprint(e)})
		}

		return nil, nil, &errs
	}

	ID := entity.NewID()
	Timestamp := time.Now()

	//Loop through the struct to generate events and validate
	errs := entity.Error{Op: "Create", Kind: entity.ValidationFailed, ErrorMessage: "Provide valid Payload", Severity: logrus.InfoLevel}
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
		fieldName := fields.Field(i).Name
		field := fields.Field(i)
		value := values.Field(i)

		switch field.Name {
		case "Product":
			_, err := s.productRepo.FindOneByID(entity.ID(value.String()))
			switch err {
			case entity.ErrNotFound:
				errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Product with ID " + value.String() + " doesn't Exist"})
			default:
				if err != nil {
					return nil, nil, &entity.Error{Op: "Create", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
				}
			}

		case "Attributes":
			duplicatedVariant, err := s.storeRepo.FindOneByAttribute(createVariantDTO.Product, createVariantDTO.Attributes)
			if err != entity.ErrNotFound && err != nil {
				return nil, nil, &entity.Error{Op: "Create", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
			}

			if duplicatedVariant != nil {
				errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Variant Attributes Duplication with ID " + string(duplicatedVariant.ID)})
			}

			if len(createVariantDTO.Attributes) > 3 {
				errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Max Attributes is 3"})
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

	if len(errs.Errors) >= 1 {
		return nil, nil, &errs
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
		return nil, nil, &entity.Error{Op: "Create", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
	}

	_, err = s.storeRepo.Create(v)
	if err != nil {
		return nil, nil, &entity.Error{Op: "Create", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
	}

	s.msgRepo.SendMessages(messages)

	Version := int32(v.Version)
	return &ID, &Version, nil
}

//UpdateVariantDTO update variant DTO
type UpdateVariantDTO struct {
	SKU      string `json:"sku,omitempty" validate:"omitempty" structs:"sku,omitempty"`
	Quantity int    `json:"quantity,omitempty" validate:"omitempty" structs:"quantity,omitempty"`
	Price    int    `json:"price,omitempty" validate:"omitempty,min=1" structs:"price,omitempty"`
	Image    string `json:"image,omitempty" validate:"omitempty,uri" structs:"image,omitempty"`
}

//UpdateOne product
func (s *Service) UpdateOne(ID entity.ID, v int32, updateVariantDTO UpdateVariantDTO) (*int32, *entity.Error) {
	//Validate DTOs, Terminate the Create process if the input is not valid
	if err := validator.New().Struct(updateVariantDTO); err != nil {
		errs := entity.Error{Op: "UpdateOne", Kind: entity.ValidationFailed, ErrorMessage: "Provide valid Payload", Severity: logrus.InfoLevel}

		for _, e := range err.(validator.ValidationErrors) {
			errs.Errors = append(errs.Errors, entity.ErrorField{Field: e.Field(), Error: fmt.Sprint(e)})
		}

		return nil, &errs
	}

	Timestamp := time.Now()
	version := entity.Version(v)

	variant, err := s.storeRepo.FindOneByID(ID)
	switch err {
	case entity.ErrNotFound:
		return nil, &entity.Error{Op: "Update", Kind: entity.NotFound, ErrorMessage: entity.ErrorMessage("Variant with id " + string(ID) + " Not found"), Severity: logrus.InfoLevel}
	default:
		if err != nil {
			return nil, &entity.Error{Op: "Update", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
		}
	}

	if version != variant.Version {
		return nil, &entity.Error{Op: "Update", Kind: entity.ConcurrentModification, ErrorMessage: entity.ErrorMessage("Version conflict"), Severity: logrus.InfoLevel}
	}

	//Loop through the struct to generate events and validate
	errs := entity.Error{Op: "Create", Kind: entity.ValidationFailed, ErrorMessage: "Provide valid Payload", Severity: logrus.InfoLevel}
	var messages []*entity.Message

	fields := reflect.TypeOf(updateVariantDTO)
	values := reflect.ValueOf(updateVariantDTO)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		fieldName := fields.Field(i).Name
		field := fields.Field(i)
		value := values.Field(i)

		switch field.Name {
		case "SKU":
			if value.String() != "" {
				if variant.SKU == updateVariantDTO.SKU {
					errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Sku already updated"})
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
					errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Quantity already updated"})
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
					errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Image already updated"})
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
					errs.Errors = append(errs.Errors, entity.ErrorField{Field: fieldName, Error: "Price already updated"})
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

	if len(errs.Errors) > 0 {
		return nil, &errs
	}

	if version == variant.Version {
		return nil, &entity.Error{Op: "Update", Kind: entity.NoUpdates, ErrorMessage: entity.ErrorMessage("No updates found"), Severity: logrus.InfoLevel}
	}

	c := &entity.Command{AggregateID: string(ID), Type: "UpdateProduct", Payload: structs.Map(updateVariantDTO), Timestamp: Timestamp}
	_, err = s.storeRepo.StoreCommand(c)
	if err != nil {
		return nil, &entity.Error{Op: "Update", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
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
		return nil, &entity.Error{Op: "Update", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
	}

	if updatedNum != 1 {
		return nil, &entity.Error{Op: "Update", Kind: entity.ConcurrentModification, ErrorMessage: entity.ErrorMessage("Version conflict"), Severity: logrus.InfoLevel}
	}

	//TODO:handle failure cases
	s.msgRepo.SendMessages(messages)

	Version := int32(version)
	return &Version, nil
}

//Delete product
func (s *Service) Delete(id entity.ID, v int32) *entity.Error {
	Timestamp := time.Now()
	version := entity.Version(v)

	//TODO check Transactions
	p, err := s.storeRepo.FindOneByID(id)
	switch err {
	case entity.ErrNotFound:
		return &entity.Error{Op: "Delete", Kind: entity.NotFound, ErrorMessage: entity.ErrorMessage("Variant with id " + string(id) + " Not found"), Severity: logrus.InfoLevel}
	default:
		if err != nil {
			return &entity.Error{Op: "Delete", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
		}
	}

	if version != p.Version {
		return &entity.Error{Op: "Update", Kind: entity.ConcurrentModification, ErrorMessage: entity.ErrorMessage("Version conflict"), Severity: logrus.InfoLevel}

	}

	c := &entity.Command{AggregateID: string(id), Type: "DeleteVariant", Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	updatedNum, err := s.storeRepo.DeleteOne(id, version)
	if err != nil {
		return &entity.Error{Op: "Delete", Kind: entity.Unexpected, ErrorMessage: "Internal Server Error", Severity: logrus.ErrorLevel}
	}

	if updatedNum != 1 {
		return &entity.Error{Op: "Update", Kind: entity.ConcurrentModification, ErrorMessage: entity.ErrorMessage("Version conflict"), Severity: logrus.InfoLevel}

	}

	version++

	//TODO:handle failure cases
	m := &entity.Message{ID: string(id), Type: "PRODUCT_VARIANT_DELETED", Version: version, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return nil
}
