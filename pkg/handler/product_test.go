package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/markus-azer/products-service/pkg/entity"
	"github.com/markus-azer/products-service/pkg/product"
	"github.com/stretchr/testify/assert"
)

func TestProductIndex(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	service := product.NewMockUseCase(controller)
	r := mux.NewRouter()
	MakeProductHandlers(r, service)
	path, err := r.GetRoute("CreateProduct").GetPathTemplate()

	assert.Nil(t, err)
	assert.Equal(t, "/v1/products/command/create", path)

	ID := entity.NewID()
	p := product.CreateProductDTO{
		Name:        "Test product",
		Description: "Test product description",
	}

	v := entity.Version(3)
	service.EXPECT().Create(p).Return(&ID, &v, nil)
	create := create(service)

	ts := httptest.NewServer(create)
	defer ts.Close()

	payload := fmt.Sprintf(`{
		"name": "Test product",
		"description": "Test product description"
	  }`)

	_, err = http.Post(ts.URL+"/v1/products/command/create", "application/json", strings.NewReader(payload))
	// var res *handler.Response
	// json.NewDecoder(resp.Body).Decode(&res)
	// assert.Equal(t, true, res.Successful)
}
