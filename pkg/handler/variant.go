package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

		ID, v, errs := service.Create(variant)

		if len(errs) > 0 {
			var errsString = []string{}
			for _, i := range errs {
				j := string(i)
				errsString = append(errsString, j)
			}

			payload, _ := json.Marshal(&Response{Errors: errsString, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		payload, _ := json.Marshal(&Response{Message: "Created Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true})
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	})
}

func updateVariant(service variant.UseCase) http.Handler {
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

		var variant variant.UpdateVariantDTO
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() //WARNNING return only one unknown field

		err = dec.Decode(&variant)

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

		v, errs := service.UpdateOne(ID, int32(version), variant)

		if len(errs) > 0 {
			var errsString = []string{}
			for _, i := range errs {
				j := string(i)
				errsString = append(errsString, j)
			}
			payload, _ := json.Marshal(&Response{Errors: errsString, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		payload, _ := json.Marshal(&Response{Message: "Updated Successfully", Data: map[string]interface{}{"id": ID, "version": v}, Successful: true})
		w.WriteHeader(http.StatusAccepted)
		w.Write(payload)
	})
}

func deleteVariant(service variant.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := entity.ID(vars["id"])
		version, err := strconv.Atoi(vars["version"])
		if err != nil {
			payload, _ := json.Marshal(&Response{Errors: []string{"Provide Valid version value"}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		clientErr := service.Delete(id, int32(version))
		if clientErr != nil {
			payload, _ := json.Marshal(&Response{Errors: []string{string(*clientErr)}, Successful: false})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(payload)
			return
		}

		payload, _ := json.Marshal(&Response{Message: "Deleted Successfully", Data: map[string]interface{}{}, Successful: true})
		w.WriteHeader(http.StatusAccepted)
		w.Write(payload)
	})
}

//MakeVariantHandlers make url handlers
func MakeVariantHandlers(r *mux.Router, service variant.UseCase) {
	r.Handle("/v1/variants/create", createVariant(service)).Methods("POST", "OPTIONS").Name("CreateVariant")
	r.Handle("/v1/variants/{id}/{version}/update", updateVariant(service)).Methods("PATCH", "OPTIONS").Name("UpdateVariant")
	r.Handle("/v1/variants/{id}/{version}/delete", deleteVariant(service)).Methods("DELETE", "OPTIONS").Name("DeleteVariant")
}
