package product

import (
	"context"
	"log"

	"github.com/MarkusAzer/products-service/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//MongoRepository mongodb repo
type MongoRepository struct {
	db *mongo.Database
}

//NewMongoRepository create new repository
func NewMongoRepository(db *mongo.Database) StoreRepository {
	return &MongoRepository{
		db: db,
	}
}

//FindOneByID find product by Id
func (r *MongoRepository) FindOneByID(id entity.ID) (entity.Product, error) {
	result := entity.Product{}
	coll := r.db.Collection("products")
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	// if err != nil {
	// 	return result, err
	// }

	return result, err
}

//StoreCommand persistence commands
func (r *MongoRepository) StoreCommand(c *entity.Command) {
	coll := r.db.Collection("commands-product")
	_, err := coll.InsertOne(context.TODO(), c)
	if err != nil {
		log.Println("Error on creating Command", err)
	}
}

//Create create new Product
func (r *MongoRepository) Create(p *entity.Product) {
	coll := r.db.Collection("products")
	_, err := coll.InsertOne(context.TODO(), p)
	if err != nil {
		log.Println("Error on creating Product", err)
	}
}

//UpdateOne update an existing product
func (r *MongoRepository) UpdateOne(id entity.ID, p *entity.Product) {
	coll := r.db.Collection("products")
	_, err := coll.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: p}})
	if err != nil {
		log.Println("Error on updating Product", err)
	}
}

//UpdateOneP update an existing product
func (r *MongoRepository) UpdateOneP(id entity.ID, p *entity.UpdateProduct) {
	coll := r.db.Collection("products")
	_, err := coll.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: p}})
	if err != nil {
		log.Println("Error on updating Product", err)
	}
}

//DeleteOne update an existing Product
func (r *MongoRepository) DeleteOne(id entity.ID) {
	coll := r.db.Collection("products")
	_, err := coll.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		log.Println("Error on deleting Product", err)
	}
}
