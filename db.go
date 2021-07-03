package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"log"
	"time"
)

// GetDatabase - connects to MongoDB instance and returns database with name from URI and Client
func GetDatabase(uri string) (*mongo.Database, *mongo.Client) {
	cs, err := connstring.ParseAndValidate(uri)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	log.Println("Successfully connected to MongoDB at", uri)

	database := client.Database(cs.Database)

	return database, client
}
