package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/MarkusAzer/products-service/pkg/entity"
	"github.com/MarkusAzer/products-service/pkg/product"
	"github.com/gorilla/mux"
)

// Response struct which contains an API Response
type Response struct {
	Message     string                 `json:"message,omitempty"`
	Validations []string               `json:"validations,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Successful  bool                   `json:"successful"`
}

func create(service product.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding Product"
		//TODO: validate body
		var p *entity.Product
		err := json.NewDecoder(r.Body).Decode(&p)
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

		// // Validate product
		// validErrs := p.Validate()
		// if len(validErrs) > 0 {
		// 	response := Response{Message: "Validations Errors", Validations: validErrs, Successful: false}

		// 	payload, err := json.Marshal(response)
		// 	if err != nil {
		// 		log.Println(err)
		// 		w.WriteHeader(http.StatusInternalServerError)
		// 		return
		// 	}

		// 	w.WriteHeader(http.StatusBadRequest)
		// 	//TODO better middleware https://stackoverflow.com/questions/51456253/how-to-set-http-responsewriter-content-type-header-globally-for-all-api-endpoint
		// 	w.Header().Add("Content-Type", "application/json")
		// 	w.Write(payload)
		// 	return
		// }

		p.ID, err = service.Create(p)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(p); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

// func update(service product.UseCase) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// }

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
	// r.Handle("/v1/products/command/{id}/update", update(service)).Methods("GET", "OPTIONS").Name("UpdateProduct")
	r.Handle("/v1/products/command/{id}/{version}/publish", publish(service)).Methods("POST", "OPTIONS").Name("PublishProduct")
	r.Handle("/v1/products/command/{id}/{version}/unpublish", unpublish(service)).Methods("POST", "OPTIONS").Name("UnpublishProduct")
	r.Handle("/v1/products/command/{id}/{version}/update-price", updatePrice(service)).Methods("POST", "OPTIONS").Name("UpdatePrice")
	r.Handle("/v1/products/command/{id}/{version}/delete", delete(service)).Methods("POST", "OPTIONS").Name("DeleteProduct")
}
