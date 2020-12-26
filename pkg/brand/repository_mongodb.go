package brand

import (
	"context"
	"log"

	"github.com/markus-azer/products-service/pkg/entity"
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

//FindOneByID find brand by Id
func (r *MongoRepository) FindOneByID(id entity.ID) (entity.Brand, error) {
	result := entity.Brand{}
	coll := r.db.Collection("brands")
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

//FindOneByName find brand by Name
func (r *MongoRepository) FindOneByName(name string) (*entity.Brand, error) {
	result := entity.Brand{}
	coll := r.db.Collection("brands")
	err := coll.FindOne(context.TODO(), bson.M{"name": name}).Decode(&result)

	switch err {
	case nil:
		return &result, nil
	case mongo.ErrNoDocuments:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//Create create new Brand
func (r *MongoRepository) Create(e *entity.Brand) {
	coll := r.db.Collection("brands")
	_, err := coll.InsertOne(context.TODO(), e)
	if err != nil {
		log.Println("Error on creating Brand", err)
	}
}

//UpdateOne update an existing brand
func (r *MongoRepository) UpdateOne(id entity.ID, e *entity.Brand) {
	coll := r.db.Collection("brands")
	_, err := coll.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: e}})
	if err != nil {
		log.Println("Error on updating Brand", err)
	}
}

//DeleteOne update an existing Brand
func (r *MongoRepository) DeleteOne(id entity.ID) {
	coll := r.db.Collection("brands")
	_, err := coll.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		log.Println("Error on deleting Brand", err)
	}
}
