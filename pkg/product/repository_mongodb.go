package product

import (
	"context"

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

//TODO: USE Sessions

//FindOneByID find product by Id
func (r *MongoRepository) FindOneByID(id entity.ID) (*entity.Product, error) {
	result := entity.Product{}
	coll := r.db.Collection("products")
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	if err != nil { //TODO: check if we should remove it
		return &result, err
	}

	return &result, err
}

//StoreCommand persistence commands
func (r *MongoRepository) StoreCommand(c *entity.Command) (*entity.ID, error) {
	coll := r.db.Collection("commands-product")

	result, err := coll.InsertOne(context.TODO(), c)

	if err != nil {
		return nil, err
	}

	str := result.InsertedID.(primitive.ObjectID).Hex()
	id := entity.ID(str)

	return &id, nil
}

//Create create new Product
func (r *MongoRepository) Create(p *entity.Product) (*entity.ID, error) {
	coll := r.db.Collection("products")

	result, err := coll.InsertOne(context.TODO(), p)

	if err != nil {
		return nil, err
	}

	str := result.InsertedID.(primitive.ObjectID).Hex()
	id := entity.ID(str)

	return &id, err
}

//UpdateOne update an existing product
func (r *MongoRepository) UpdateOne(id entity.ID, p *entity.Product, v entity.Version) (int, error) {

	coll := r.db.Collection("products")

	result, err := coll.UpdateOne(
		context.TODO(),
		bson.D{primitive.E{Key: "_id", Value: id}, primitive.E{Key: "_V", Value: v}},
		bson.D{primitive.E{Key: "$set", Value: p}},
	)

	if err != nil {
		return int(result.ModifiedCount), err
	}

	return int(result.ModifiedCount), nil
}

//UpdateOneP update an existing product
func (r *MongoRepository) UpdateOneP(id entity.ID, p *entity.UpdateProduct, v entity.Version) (int, error) {
	coll := r.db.Collection("products")

	result, err := coll.UpdateOne(
		context.TODO(),
		bson.D{primitive.E{Key: "_id", Value: id}, primitive.E{Key: "_V", Value: v}},
		bson.D{primitive.E{Key: "$set", Value: p}},
	)

	if err != nil {
		return int(result.ModifiedCount), err
	}

	return int(result.ModifiedCount), nil
}

//DeleteOne update an existing Product
func (r *MongoRepository) DeleteOne(id entity.ID, v entity.Version) (int, error) {
	coll := r.db.Collection("products")

	result, err := coll.DeleteOne(
		context.TODO(),
		bson.D{primitive.E{Key: "_id", Value: id}, primitive.E{Key: "_V", Value: v}},
	)

	if err != nil {
		return int(result.DeletedCount), err
	}

	return int(result.DeletedCount), nil

}
