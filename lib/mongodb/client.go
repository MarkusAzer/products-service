package mongodb

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/markus-azer/products-service/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connected = "Successfully connected to database: "

//MongoDatastore contains client db Session
type MongoDatastore struct {
	Db *mongo.Database
}

//NewDatastore Function that creates Mongo DataStore
func NewDatastore(config config.GeneralConfig) *MongoDatastore {

	var mongoDataStore *MongoDatastore
	db := connect(config)
	if db != nil {

		//TODO: log statements here as well

		mongoDataStore = new(MongoDatastore)
		mongoDataStore.Db = db

		return mongoDataStore
	}

	log.Fatal("Failed to connect to database: ", config.DatabaseName)

	return nil
}

func connect(generalConfig config.GeneralConfig) (a *mongo.Database) {
	var connectOnce sync.Once
	var db *mongo.Database
	connectOnce.Do(func() {
		db = connectToMongo(generalConfig)
	})

	return db
}

func connectToMongo(generalConfig config.GeneralConfig) (a *mongo.Database) {

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	session, err := mongo.Connect(ctx, options.Client().ApplyURI(generalConfig.DatabaseHost))
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	err = session.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Failed to connect to Mongo db ", generalConfig.DatabaseHost)
	}

	var DB = session.Database(generalConfig.DatabaseName)
	fmt.Println(connected, generalConfig.DatabaseName)

	return DB
}
