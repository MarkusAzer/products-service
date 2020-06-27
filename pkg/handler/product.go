package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/MarkusAzer/products-service/pkg/entity"
	"github.com/MarkusAzer/products-service/pkg/product"
	"github.com/gorilla/mux"
)

//Validation specifies data serialization/deserialization protocol.

// DisallowUnknownFields https://maori.geek.nz/golang-raise-error-if-unknown-field-in-json-with-exceptions-2b0caddecd1

// Response struct which contains an API Response
type Response struct {
	Message    string                 `json:"message,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Errors     []string               `json:"errors,omitempty"`
	Successful bool                   `json:"successful"`
}

func create(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var p product.CreateProductDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err := dec.Decode(&p)

		switch {
		case err == io.EOF:
			payload, _ := json.Marshal(&Response{Errors: []string{"Provide valid Body"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		case err != nil && strings.Contains(err.Error(), "json: unknown field"):
			m := regexp.MustCompile(`\"(.*)\"`)
			field := m.FindString(err.Error())

			payload, _ := json.Marshal(&Response{Errors: []string{field + " : Not Allowed"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		case err != nil:
			payload, _ := json.Marshal(&Response{Errors: []string{"Provide valid Body"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		ID, v, errs := service.Create(p)

		if len(errs) > 0 {
			payload, _ := json.Marshal(&Response{Errors: errs, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		payload, _ := json.Marshal(&Response{Message: "Created Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true})
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	})
}

func update(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		ID := entity.ID(vars["id"])
		version, err := strconv.Atoi(vars["version"])
		if err != nil {
			payload, _ := json.Marshal(&Response{Errors: []string{"Internal Server Error"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		var p product.UpdateProductDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err = dec.Decode(&p)

		switch {
		case err == io.EOF:
			payload, _ := json.Marshal(&Response{Errors: []string{"Provide valid Body"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		case err != nil && strings.Contains(err.Error(), "json: unknown field"):
			m := regexp.MustCompile(`\"(.*)\"`)
			field := m.FindString(err.Error())

			payload, _ := json.Marshal(&Response{Errors: []string{field + " : Not Allowed"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		case err != nil:
			payload, _ := json.Marshal(&Response{Errors: []string{"Provide valid Body"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		v, errs := service.UpdateOne(ID, int32(version), p)

		if len(errs) > 0 {
			payload, _ := json.Marshal(&Response{Errors: errs, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		payload, _ := json.Marshal(&Response{Message: "Updated Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true})
		w.WriteHeader(http.StatusAccepted)
		w.Write(payload)
	})
}

func publish(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := entity.ID(vars["id"])
		version, err := strconv.Atoi(vars["version"])
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide Valid version value"))
			return
		}

		updatedVersion, err := service.Publish(id, int32(version))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(strconv.Itoa(int(updatedVersion))))
	})
}

func unpublish(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := entity.ID(vars["id"])
		version, err := strconv.Atoi(vars["version"])
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide Valid version value"))
			return
		}

		updatedVersion, err := service.Unpublish(id, int32(version))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(strconv.Itoa(int(updatedVersion))))
	})
}

func updatePrice(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := entity.ID(vars["id"])
		version, err := strconv.Atoi(vars["version"])
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide Valid version value"))
			return
		}

		//TODO extract price
		var p *entity.Product
		err = json.NewDecoder(r.Body).Decode(&p)
		//TODO: notify not allowed fields
		switch {
		case err == io.EOF:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide valid Body"))
			return
		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide valid JSON"))
			return
		}

		updatedVersion, err := service.UpdatePrice(id, int32(version), int(p.Price))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(strconv.Itoa(int(updatedVersion))))
	})
}

func delete(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := entity.ID(vars["id"])
		version, err := strconv.Atoi(vars["version"])
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Provide Valid version value"))
			return
		}

		err = service.Delete(id, int32(version))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	})
}

//MakeProductHandlers make url handlers
func MakeProductHandlers(r *mux.Router, service product.UseCase) {
	r.Handle("/v1/products/command/create", create(service)).Methods("POST", "OPTIONS").Name("CreateProduct")
	r.Handle("/v1/products/command/{id}/{version}/update", update(service)).Methods("POST", "OPTIONS").Name("UpdateProduct")
	r.Handle("/v1/products/command/{id}/{version}/publish", publish(service)).Methods("POST", "OPTIONS").Name("PublishProduct")
	r.Handle("/v1/products/command/{id}/{version}/unpublish", unpublish(service)).Methods("POST", "OPTIONS").Name("UnpublishProduct")
	r.Handle("/v1/products/command/{id}/{version}/update-price", updatePrice(service)).Methods("POST", "OPTIONS").Name("UpdatePrice")
	r.Handle("/v1/products/command/{id}/{version}/delete", delete(service)).Methods("POST", "OPTIONS").Name("DeleteProduct")
}
