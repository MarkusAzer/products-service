package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/markus-azer/products-service/pkg/entity"
	"github.com/markus-azer/products-service/pkg/variant"
)

func createVariant(service variant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var variant variant.CreateVariantDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err := dec.Decode(&variant)

		if err != nil {
			payload := serializationErrorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		ID, v, err := service.Create(variant)

		if err != nil {
			payload := errorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		payload := &response{StatusCode: http.StatusCreated, Message: "Created Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true}
		w.WriteHeader(payload.StatusCode)
		json.NewEncoder(w).Encode(payload)
	})
}

func updateVariant(service variant.UseCase) http.Handler {
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

		var variant variant.UpdateVariantDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err = dec.Decode(&variant)

		if err != nil {
			payload := serializationErrorHandler(err)
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		v, err := service.UpdateOne(ID, int32(version), variant)

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

func deleteVariant(service variant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := entity.ID(vars["id"])
		version, err := strconv.ParseInt(vars["version"], 0, 8)
		if err != nil {
			payload := &response{StatusCode: 500, Message: "Internal Service Error", Successful: false}
			w.WriteHeader(payload.StatusCode)
			json.NewEncoder(w).Encode(payload)
			return
		}

		err = service.Delete(id, int32(version))
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

//MakeVariantHandlers make url handlers
func MakeVariantHandlers(r *mux.Router, service variant.UseCase) {
	r.Handle("/v1/variants/create", createVariant(service)).Methods("POST", "OPTIONS").Name("CreateVariant")
	r.Handle("/v1/variants/{id}/{version}/update", updateVariant(service)).Methods("PATCH", "OPTIONS").Name("UpdateVariant")
	r.Handle("/v1/variants/{id}/{version}/delete", deleteVariant(service)).Methods("DELETE", "OPTIONS").Name("DeleteVariant")
}
