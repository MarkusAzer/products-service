package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markus-azer/products-service/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markus-azer/products-service/api/handler"
	"github.com/markus-azer/products-service/api/metric"
	"github.com/markus-azer/products-service/api/middleware"
	kafkaStore "github.com/markus-azer/products-service/lib/kafka"
	"github.com/markus-azer/products-service/lib/mongodb"
	"github.com/markus-azer/products-service/pkg/brand"
	"github.com/markus-azer/products-service/pkg/product"
	"github.com/markus-azer/products-service/pkg/variant"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Swagger https://medium.com/@ribice/serve-swaggerui-within-your-golang-application-5486748a5ed4
//https://github.com/go-swagger/go-swagger/issues/370

func main() {

	mongoDatastore := mongodb.NewDatastore(config.DevConfig)

	client, err := kafkaStore.NewKafkaClient()
	check(err)

	fmt.Printf("Kafka Client Created %v\n", client)
	fmt.Printf("Mongo Client Created %v\n", mongoDatastore)
	r := mux.NewRouter()

	productStoreRepo := product.NewMongoRepository(mongoDatastore.Db)
	productMsgRepo := product.NewKafkaRepository(client.Producer)

	variantStoreRepo := variant.NewMongoRepository(mongoDatastore.Db)
	variantMsgRepo := variant.NewKafkaRepository(client.Producer)

	brandStoreRepo := brand.NewMongoRepository(mongoDatastore.Db)
	brandMsgRepo := brand.NewKafkaRepository(client.Consumer)

	productService := product.NewService(productMsgRepo, productStoreRepo, brandStoreRepo)
	variantService := variant.NewService(variantMsgRepo, variantStoreRepo, productStoreRepo)
	brandService := brand.NewService(brandStoreRepo)

	metricService, err := metric.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}

	//Middlewares
	r.Use(middleware.Logging)
	r.Use(handlers.CORS())
	r.Use(middleware.Metrics(metricService))
	r.Use(middleware.ValidateHeaderType)
	r.Use(middleware.SetResHeaderType)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Route Handlers - Endpoints
	handler.MakeProductHandlers(r, productService)
	handler.MakeBrandHandlers(brandMsgRepo, brandService)
	handler.MakeVariantHandlers(r, variantService)

	log.Fatal(http.ListenAndServe(":8080", r))
}
