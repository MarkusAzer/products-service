package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markus-azer/products-service/pkg/entity"
	"github.com/markus-azer/products-service/pkg/product"
)

//Validation specifies data serialization/deserialization protocol.

// DisallowUnknownFields https://maori.geek.nz/golang-raise-error-if-unknown-field-in-json-with-exceptions-2b0caddecd1

func create(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var p product.CreateProductDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err := dec.Decode(&p)

		if err != nil {
			payload := serializationErrorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		ID, v, err := service.Create(p)

		if err != nil {
			payload := errorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		payload := &response{StatusCode: http.StatusCreated, Message: "Created Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true}
		w.WriteHeader(payload.StatusCode)
		json.NewEncoder(w).Encode(payload)
		return
	})
}

func update(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		ID := entity.ID(vars["id"])
		version, err := strconv.ParseInt(vars["version"], 0, 8)
		if err != nil {
			payload := &response{StatusCode: 500, Message: "Internal Service Error", Successful: false}
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		var p product.UpdateProductDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err = dec.Decode(&p)

		if err != nil {
			payload := serializationErrorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		v, err := service.UpdateOne(ID, int32(version), p)

		if err != nil {
			payload := errorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		payload := &response{StatusCode: http.StatusAccepted, Message: "Updated Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true}
		w.WriteHeader(payload.StatusCode)
		json.NewEncoder(w).Encode(payload)
	})
}

func delete(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID := entity.ID(vars["id"])
		version, err := strconv.ParseInt(vars["version"], 0, 8)
		if err != nil {
			payload := &response{StatusCode: http.StatusBadRequest, Message: "Provide Valid version value", Successful: false}
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		err = service.Delete(ID, int32(version))
		if err != nil {
			payload := errorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		payload := &response{StatusCode: http.StatusAccepted, Message: "Deleted Successfully", Data: map[string]interface{}{}, Successful: true}
		w.WriteHeader(payload.StatusCode)
		json.NewEncoder(w).Encode(payload)
	})
}

//MakeProductHandlers make url handlers
func MakeProductHandlers(r *mux.Router, service product.UseCase) {
	r.Handle("/v1/products/command/create", create(service)).Methods("POST", "OPTIONS").Name("CreateProduct")
	r.Handle("/v1/products/command/{id}/{version}/update", update(service)).Methods("PATCH", "OPTIONS").Name("UpdateProduct")
	r.Handle("/v1/products/command/{id}/{version}/delete", delete(service)).Methods("DELETE", "OPTIONS").Name("DeleteProduct")
}
