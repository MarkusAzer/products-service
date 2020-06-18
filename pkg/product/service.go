package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MarkusAzer/products-service/pkg/brand"
	"github.com/MarkusAzer/products-service/pkg/entity"
	"github.com/fatih/structs"
)

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

//TODO struct to map funchttps://stackoverflow.com/questions/23589564/function-for-converting-a-struct-to-map-in-golang

//Create new product
func (s *Service) Create(p *entity.Product) (entity.ID, error) {
	ID := entity.NewID()
	Timestamp := time.Now()

	p.ID = ID
	p.Version = 1
	p.CreatedAt = entity.Time(Timestamp)

	data, err := json.Marshal(p) // Convert to a json string

	if err != nil {
		return "", err
	}
	var newMap map[string]interface{}
	err = json.Unmarshal(data, &newMap) // Convert to a map

	c := &entity.Command{AggregateID: string(ID), Type: "CreateProduct", Payload: newMap, Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	s.storeRepo.Create(p)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_CREATED", Version: 1, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return ID, err
}

//UpdateOne product
func (s *Service) UpdateOne(id entity.ID, version int32, p *entity.UpdateProduct) (int32, error) {
	Timestamp := time.Now()

	//TODO check Transactions
	product, err := s.storeRepo.FindOneByID(id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if version != product.Version {
		fmt.Println("miss matching version")
		return 0, errors.New("miss matching version")
	}

	var errs []error
	var messages []*entity.Message

	if p.Brand != "" {
		if p.Brand == product.Brand {
			errs = append(errs, errors.New("Brand already updated"))
		}

		_, err := s.brandRepo.FindOneByName(p.Brand)

		if err != nil {
			errs = append(errs, errors.New("Brand doesnt exist"))
		}

		version++

		newMap := make(map[string]interface{})
		newMap["brand"] = p.Brand

		messages = append(messages, &entity.Message{ID: string(id), Type: "PRODUCT_BRAND_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp})

	}

	if p.Category != "" {
		if p.Category == product.Category {
			errs = append(errs, errors.New("Category already updated"))
		}

		version++

		newMap := make(map[string]interface{})
		newMap["category"] = p.Category

		messages = append(messages, &entity.Message{ID: string(id), Type: "PRODUCT_CATEGORY_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp})

	}

	if p.Description != "" {
		if p.Description == product.Description {
			errs = append(errs, errors.New("Description already updated"))
		}

		version++

		newMap := make(map[string]interface{})
		newMap["description"] = p.Description

		messages = append(messages, &entity.Message{ID: string(id), Type: "PRODUCT_DESCRIPTION_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp})

	}

	if p.Image != "" {
		if p.Image == product.Image {
			errs = append(errs, errors.New("Image already updated"))
		}

		version++

		newMap := make(map[string]interface{})
		newMap["image"] = p.Image

		messages = append(messages, &entity.Message{ID: string(id), Type: "PRODUCT_IMAGE_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp})

	}

	if p.Name != "" {
		if p.Name == product.Name {
			errs = append(errs, errors.New("Name already updated"))
		}

		version++

		newMap := make(map[string]interface{})
		newMap["name"] = p.Name

		messages = append(messages, &entity.Message{ID: string(id), Type: "PRODUCT_NAME_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp})

	}

	if p.Slug != "" {
		if p.Slug == product.Slug {
			errs = append(errs, errors.New("Slug already updated"))
		}

		version++

		newMap := make(map[string]interface{})
		newMap["slug"] = p.Slug

		messages = append(messages, &entity.Message{ID: string(id), Type: "PRODUCT_SLUG_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp})

	}

	if len(errs) > 0 {
		return 0, errs[0]
	}

	c := &entity.Command{AggregateID: string(id), Type: "UpdateProduct", Payload: structs.Map(p), Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	p.Version = version
	//TODO Patch the update
	s.storeRepo.UpdateOneP(id, p)

	//TODO:handle failure cases
	s.msgRepo.SendMessages(messages)

	return version, nil
}

//Publish publish product
func (s *Service) Publish(ID entity.ID, version int32) (int32, error) {
	Timestamp := time.Now()

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
	if p.Status == "Publish" {
		fmt.Println("Product is already published")
		return 0, errors.New("Product is already published")
	}

	newMap := make(map[string]interface{})
	newMap["Status"] = "Publish"

	c := &entity.Command{AggregateID: string(ID), Type: "PublishProduct", Payload: newMap, Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	version++

	//TODO Patch the update
	p.Status = "Publish"
	p.Version = version
	s.storeRepo.UpdateOne(ID, &p)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_PUBLISHED", Version: version, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return version, nil
}

//Unpublish unpublish product
func (s *Service) Unpublish(ID entity.ID, version int32) (int32, error) {
	Timestamp := time.Now()

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
	s.storeRepo.UpdateOne(ID, &p)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_UNPUBLISHED", Version: version, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return version, nil
}

//UpdatePrice product price
func (s *Service) UpdatePrice(ID entity.ID, version int32, price int) (int32, error) {
	Timestamp := time.Now()

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
	s.storeRepo.UpdateOne(ID, &p)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_PRICE_UPDATED", Version: version, Payload: newMap, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return version, nil
}

//Delete product
func (s *Service) Delete(ID entity.ID, version int32) error {
	Timestamp := time.Now()

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

	s.storeRepo.DeleteOne(ID)

	version++

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_DELETED", Version: version, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return nil
}
