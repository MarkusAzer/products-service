package product_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/markus-azer/products-service/pkg/brand"
	"github.com/markus-azer/products-service/pkg/entity"
	"github.com/markus-azer/products-service/pkg/product"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	productRepo := product.NewMockStoreRepository(controller)
	brandRepo := brand.NewMockStoreRepository(controller)
	messagesRepo := product.NewMockMessagesRepository(controller)

	service := product.NewService(messagesRepo, productRepo, brandRepo)

	ID := entity.NewID()
	storeID := entity.NewID()

	cp := product.CreateProductDTO{
		Name:   "Test Product",
		Price:  20,
		Seller: "test",
	}

	invalidCP := product.CreateProductDTO{
		Price:  20,
		Seller: "test",
	}

	productRepo.EXPECT().StoreCommand(gomock.Any()).Return(&storeID, nil)
	productRepo.EXPECT().Create(gomock.Any()).Return(&ID, nil)
	messagesRepo.EXPECT().SendMessages(gomock.Any())

	id, v, err := service.Create(cp)
	// https://godoc.org/golang.org/x/tools/cmd/godoc
	fmt.Println("the current version is ", v)
	// Output:
	// the current version is 3

	assert.Nil(t, err)
	assert.True(t, entity.IsValidUUID(string(*id)))
	assert.Equal(t, entity.Version(3), *v)

	id, v, err = service.Create(invalidCP)

	assert.NotNil(t, err)
	e, ok := err.(*entity.Error)

	// assert.IsType(t, entity.Error, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, entity.ValidationFailed, e.Kind)
	assert.NotNil(t, e.Errors)
	assert.Equal(t, 1, len(e.Errors))
	assert.Equal(t, "Name", e.Errors[0].Field)
	assert.Nil(t, id)
	assert.Nil(t, v)

}

func TestUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	productRepo := product.NewMockStoreRepository(controller)
	brandRepo := brand.NewMockStoreRepository(controller)
	messagesRepo := product.NewMockMessagesRepository(controller)

	service := product.NewService(messagesRepo, productRepo, brandRepo)

	ID := entity.NewID()
	storeID := entity.NewID()

	product := product.UpdateProductDTO{
		Name:  "Updated Test Product",
		Price: 25,
	}

	createdProduct := entity.Product{
		ID:      ID,
		Version: 3,
		Name:    "Test Product",
		Price:   20,
	}

	productRepo.EXPECT().StoreCommand(gomock.Any()).Return(&storeID, nil).MaxTimes(2)
	productRepo.EXPECT().FindOneByID(gomock.Any()).Return(&createdProduct, nil).MaxTimes(2)
	productRepo.EXPECT().UpdateOneP(gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
	messagesRepo.EXPECT().SendMessages(gomock.Any())

	v, err := service.UpdateOne(ID, 3, product)

	assert.Nil(t, err)
	assert.Equal(t, int32(5), *v)

	v, err = service.UpdateOne(ID, 1, product)

	assert.NotNil(t, err)
	assert.Equal(t, entity.ConcurrentModification, err.Kind)
}
