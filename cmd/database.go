package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Client

func GetConnection() (*mongo.Client, error) {

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Error("Mongo URI Required")
	}

	var err error

	if DB == nil {
		DB, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	}

	if err != nil {
		log.Fatal("Error while connecting with DB." + err.Error())
	}
	err = DB.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Error while connecting with DB." + err.Error())
	}
	return DB, nil
}

func GetAllCollections(dbName string) ([]string, error) {
	client, err := GetConnection()
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	collectionNames, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	return collectionNames, nil
}

func CollectionExists(dbName string, collectionName string) (bool, error) {

	collectionNames, err := GetAllCollections(dbName)
	if err != nil {
		return false, err
	}

	for _, name := range collectionNames {
		if name == collectionName {
			return true, nil
		}
	}
	return false, nil
}

func GetCollection(dbName string, collectionName string) (*mongo.Collection, error) {
	client, err := GetConnection()
	if err != nil {
		log.Fatal("Error while getting collection " + collectionName + " : " + err.Error())
	}

	exists, err := CollectionExists(dbName, collectionName)
	if err != nil {
		log.Fatal("Error while getting collection " + collectionName + " : " + err.Error())
	}

	if !exists {
		log.Fatal("Collection " + collectionName + " does not exists")
	}

	collection := client.Database(dbName).Collection(collectionName)

	return collection, nil
}
