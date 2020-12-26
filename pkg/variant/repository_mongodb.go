package variant

import (
	"context"

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

//TODO: USE Sessions

//FindOneByID find Variant by Id
func (r *MongoRepository) FindOneByID(id entity.ID) (*entity.Variant, error) {
	result := entity.Variant{}
	coll := r.db.Collection("variants")
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	switch err {
	case nil:
		return &result, nil
	case mongo.ErrNoDocuments:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//FindOneByAttribute find Variant by product id and attributes
func (r *MongoRepository) FindOneByAttribute(product entity.ID, attributes map[string]string) (*entity.Variant, error) {
	result := entity.Variant{}
	coll := r.db.Collection("variants")

	query := bson.M{}
	query["product"] = product
	// query["$size"] = bson.M{"$objectToArray": "$purchase_record"}
	//TODO: Add attributes
	for k, v := range attributes {
		query["attributes."+k] = v
	}
	err := coll.FindOne(context.TODO(), query).Decode(&result)

	switch err {
	case nil:
		return &result, nil
	case mongo.ErrNoDocuments:
		return nil, entity.ErrNotFound
	default:
		return nil, err
	}
}

//StoreCommand persistence commands
func (r *MongoRepository) StoreCommand(c *entity.Command) (*entity.ID, error) {
	coll := r.db.Collection("commands-variant")

	result, err := coll.InsertOne(context.TODO(), c)

	if err != nil {
		return nil, err
	}

	str := result.InsertedID.(primitive.ObjectID).String()
	id := entity.ID(str)

	return &id, nil
}

//Create create new Variant
func (r *MongoRepository) Create(variant *entity.Variant) (*entity.ID, error) {
	coll := r.db.Collection("variants")

	result, err := coll.InsertOne(context.TODO(), variant)

	if err != nil {
		return nil, err
	}

	str := result.InsertedID.(string)
	id := entity.ID(str)

	return &id, err
}

//UpdateOne update an existing Variant
func (r *MongoRepository) UpdateOne(id entity.ID, variant *entity.UpdateVariant, version entity.Version) (int, error) {

	coll := r.db.Collection("variants")

	result, err := coll.UpdateOne(
		context.TODO(),
		bson.D{primitive.E{Key: "_id", Value: id}, primitive.E{Key: "_V", Value: version}},
		bson.D{primitive.E{Key: "$set", Value: variant}},
	)

	if err != nil {
		return int(result.ModifiedCount), err
	}

	return int(result.ModifiedCount), nil
}

//DeleteOne update an existing Variant
func (r *MongoRepository) DeleteOne(id entity.ID, version entity.Version) (int, error) {
	coll := r.db.Collection("variants")

	result, err := coll.DeleteOne(
		context.TODO(),
		bson.D{primitive.E{Key: "_id", Value: id}, primitive.E{Key: "_V", Value: version}},
	)

	if err != nil {
		return int(result.DeletedCount), err
	}

	return int(result.DeletedCount), nil

}
