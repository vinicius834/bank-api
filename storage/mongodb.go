package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DBName string
	Client *mongo.Client
}

type IMongoDB interface {
	GetDatabase(databaseName string) *mongo.Database
	GetCollection(collectionName string) *mongo.Collection
	InsertOne(collectionName string, document interface{}) (*mongo.InsertOneResult, error)
	FindOne(collectionName string, filter bson.M) *mongo.SingleResult
}

func NewMongoDB(url, dbName string) *MongoDB {
	var client *mongo.Client
	var err error

	client, err = newClient(url)

	if err != nil {
		log.Fatal(err)
	}
	return &MongoDB{dbName, client}
}

func newClient(url string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return mongo.Connect(ctx, options.Client().ApplyURI(url))
}

func (mongodb *MongoDB) GetDatabase(databaseName string) *mongo.Database {
	client := mongodb.Client
	return client.Database(mongodb.DBName)
}

func (mongodb *MongoDB) GetCollection(collectionName string) *mongo.Collection {
	client := mongodb.Client
	return client.Database(mongodb.DBName).Collection(collectionName)
}

func (mongodb *MongoDB) InsertOne(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	ctx := context.Background()
	return mongodb.GetCollection(collectionName).InsertOne(ctx, document)
}

func (mongodb *MongoDB) FindOne(collectionName string, filter bson.M) *mongo.SingleResult {
	ctx := context.Background()
	return mongodb.GetCollection(collectionName).FindOne(ctx, filter)
}
