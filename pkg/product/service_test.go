package product

import (
	"testing"

	"github.com/MarkusAzer/products-service/pkg/brand"
	"github.com/MarkusAzer/products-service/pkg/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	productRepo := NewMockStoreRepository(controller)
	brandRepo := brand.NewMockStoreRepository(controller)
	messagesRepo := NewMockMessagesRepository(controller)

	service := NewService(messagesRepo, productRepo, brandRepo)

	ID := entity.NewID()
	storeID := entity.NewID()
	// Timestamp := time.Now()

	cp := CreateProductDTO{
		Name:  "Test Product",
		Price: 20,
	}

	// p := &entity.Product{
	// 	ID:        ID,
	// 	Version:   entity.Version(3),
	// 	Name:      cp.Name,
	// 	Price:     cp.Price,
	// 	Status:    "unpublish", // Init the product as unpublish
	// 	CreatedAt: Timestamp,
	// }

	// // c := &entity.Command{AggregateID: string(ID), Type: "CreateProduct", Payload: structs.Map(cp), Timestamp: Timestamp}

	// var messages []*entity.Message
	// messages = append(messages, &entity.Message{
	// 	ID:        string(ID),
	// 	Type:      "PRODUCT_DRAFT_CREATED",
	// 	Version:   entity.Version(1),
	// 	Payload:   make(map[string]interface{}),
	// 	Timestamp: Timestamp},
	// )
	productRepo.EXPECT().StoreCommand(gomock.Any()).Return(&storeID, nil)
	productRepo.EXPECT().Create(gomock.Any()).Return(&ID, nil)
	messagesRepo.EXPECT().SendMessages(gomock.Any())

	id, _, err := service.Create(cp)

	assert.Nil(t, err)
	assert.True(t, entity.IsValidUUID(string(id)))
}
