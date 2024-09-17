package mongolib

import (
	"context"

	"github.com/ragpanda/go-toolkit/persistence/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoModel interface {
	GetID() string
	TableName() string
}

type MongoDBRepository[T MongoModel] struct {
	model.Repository[T]

	collection *mongo.Collection
	client     *mongo.Client
}

func NewMongoDBRepository[T MongoModel]() *MongoDBRepository[T] {
	return &MongoDBRepository[T]{}
}

func (self *MongoDBRepository[T]) AttachConnection(ctx context.Context, client *mongo.Client, dbName string) *MongoDBRepository[T] {
	var m T
	self.client = client
	self.collection = client.Database(dbName).Collection(m.TableName())
	return self
}

func (self *MongoDBRepository[T]) GetByID(ctx context.Context, id string) (T, error) {
	var result T
	err := self.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, nil
		}
		return result, err
	}
	return result, nil
}

func (self *MongoDBRepository[T]) GetByIDList(ctx context.Context, ids []string) ([]T, error) {
	cursor, err := self.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var elem T
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (self *MongoDBRepository[T]) Create(ctx context.Context, t T) (*string, error) {
	result, err := self.collection.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}
	id := t.GetID()
	if id == "" {
		id = result.InsertedID.(primitive.ObjectID).Hex()
	}

	return &id, nil
}

func (self *MongoDBRepository[T]) Save(ctx context.Context, t T) error {
	id := t.GetID()
	opts := options.Replace().SetUpsert(true)
	_, err := self.collection.ReplaceOne(ctx, bson.M{"_id": id}, t, opts)
	return err
}

func (self *MongoDBRepository[T]) Update(ctx context.Context, t T) error {
	id := t.GetID()
	_, err := self.collection.ReplaceOne(ctx, bson.M{"_id": id}, t)
	return err
}

func (self *MongoDBRepository[T]) Delete(ctx context.Context, id string) error {
	_, err := self.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (self *MongoDBRepository[T]) DeleteList(ctx context.Context, ids []string) error {
	_, err := self.collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	return err
}

func (self *MongoDBRepository[T]) GetCollection(ctx context.Context) *mongo.Collection {
	return self.collection
}

func (self *MongoDBRepository[T]) GetRawClient(ctx context.Context) *mongo.Client {
	return self.client
}
