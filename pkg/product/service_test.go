package product

import (
	"testing"

	brandMock "github.com/MarkusAzer/products-service/pkg/brand/mock"
	"github.com/MarkusAzer/products-service/pkg/entity"
	productMock "github.com/MarkusAzer/products-service/pkg/product/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	productRepo := productMock.NewMockStoreRepository(controller)
	brandRepo := brandMock.NewMockStoreRepository(controller)
	messagesRepo := productMock.NewMockMessagesRepository(controller)

	service := NewService(messagesRepo, productRepo, brandRepo)

	ID := entity.NewID()
	storeID := entity.NewID()

	p := CreateProductDTO{
		Name:  "Test Product",
		Price: 20,
	}

	productRepo.EXPECT().StoreCommand(gomock.Any()).Return(&storeID, nil)
	productRepo.EXPECT().Create(gomock.Any()).Return(&ID, nil)
	messagesRepo.EXPECT().SendMessage(gomock.Any())

	id, _, err := service.Create(p)

	assert.Nil(t, err)
	assert.True(t, entity.IsValidUUID(string(id)))
}
