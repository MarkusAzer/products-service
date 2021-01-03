package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/markus-azer/products-service/pkg/entity"
	"github.com/markus-azer/products-service/pkg/product"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ID := entity.NewID()
	v := entity.Version(3)
	service := product.NewMockUseCase(controller)
	service.EXPECT().Create(gomock.Any()).Return(&ID, &v, nil)

	// test routing
	r := mux.NewRouter()
	MakeProductHandlers(r, service)
	path, err := r.GetRoute("CreateProduct").GetPathTemplate()

	assert.Nil(t, err)
	assert.Equal(t, "/v1/products", path)

	payload := []byte(`{
		"name": "Test product",
		"description": "Test product description"
	  }`)

	req, err := http.NewRequest("POST", "localhost:8080/v1/products", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	rec := httptest.NewRecorder()

	create(service).ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	var resp *response
	json.NewDecoder(res.Body).Decode(&resp)
	assert.Equal(t, string(ID), resp.Data["id"])
	assert.Equal(t, float64(v), resp.Data["version"])

}
