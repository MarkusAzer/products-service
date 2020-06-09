package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MarkusAzer/products-service/config"

	kafkaStore "github.com/MarkusAzer/products-service/lib/kafka"
	"github.com/MarkusAzer/products-service/lib/mongodb"
	"github.com/MarkusAzer/products-service/pkg/brand"
	"github.com/MarkusAzer/products-service/pkg/handler"
	"github.com/MarkusAzer/products-service/pkg/middleware"
	"github.com/MarkusAzer/products-service/pkg/product"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	mongoDatastore := mongodb.NewDatastore(config.DevConfig)

	client, err := kafkaStore.NewKafkaClient()
	check(err)

	fmt.Printf("Kafka Client Created %v\n", client)
	fmt.Printf("Mongo Client Created %v\n", mongoDatastore)
	r := mux.NewRouter()

	productStoreRepo := product.NewMongoRepository(mongoDatastore.Db)
	productMsgRepo := product.NewKafkaRepository(client.Producer)

	brandStoreRepo := brand.NewMongoRepository(mongoDatastore.Db)
	brandMsgRepo := brand.NewKafkaRepository(client.Consumer)

	productService := product.NewService(productMsgRepo, productStoreRepo, brandStoreRepo)
	brandService := brand.NewService(brandStoreRepo)

	//Middlewares
	r.Use(middleware.Logging)
	r.Use(handlers.CORS())
	r.Use(middleware.ValidateHeaderType)

	http.Handle("/", r)
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Route Handlers - Endpoints
	handler.MakeProductHandlers(r, productService)
	handler.MakeBrandHandlers(brandMsgRepo, brandService)

	log.Fatal(http.ListenAndServe(":8080", r))
}
