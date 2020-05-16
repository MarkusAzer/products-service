package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/MarkusAzer/products-service/pkg/entity"
)

//Service service interface
type Service struct {
	msgRepo   MessagesRepository
	storeRepo StoreRepository
}

//NewService create new service
func NewService(msgR MessagesRepository, storeR StoreRepository) *Service {
	return &Service{
		msgRepo:   msgR,
		storeRepo: storeR,
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

// //UpdateOne product
// func (s *Service) UpdateOne(id entity.ID, e *entity.Product) (int, error) {
// }

//Publish publish product
func (s *Service) Publish(ID entity.ID, version int32) (int32, error) {
	Timestamp := time.Now()

	//TODO check Transactions
	p, err := s.storeRepo.FindOneByID(ID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	if version != p.Version+1 {
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

	if version != p.Version+1 {
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

	if version != p.Version+1 {
		fmt.Println("miss matching version")
		return 0, errors.New("miss matching version")
	}

	newMap := make(map[string]interface{})
	newMap["Price"] = price

	c := &entity.Command{AggregateID: string(ID), Type: "UpdateProductPrice", Payload: newMap, Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

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

	if version != p.Version+1 {
		fmt.Println("miss matching version")
		return errors.New("miss matching version")
	}

	c := &entity.Command{AggregateID: string(ID), Type: "DeleteProduct", Timestamp: Timestamp}
	s.storeRepo.StoreCommand(c)

	s.storeRepo.DeleteOne(ID)

	//TODO:handle failure cases
	m := &entity.Message{ID: string(ID), Type: "PRODUCT_DELETED", Version: version, Timestamp: Timestamp}
	s.msgRepo.SendMessage(m)

	return nil
}
