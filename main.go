package main

import (
	"log"
	"net/http"

	"github.com/MarkusAzer/products-service/pkg/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	r := mux.NewRouter()

	//Middlewares
	r.Use(middleware.Logging)
	r.Use(handlers.CORS())
	r.Use(middleware.SetHeaderJSON)

	http.Handle("/", r)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
