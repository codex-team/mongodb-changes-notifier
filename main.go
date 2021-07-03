package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func iterateChangeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()
	for stream.Next(routineCtx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}
		fmt.Printf("%v\n", data)
	}
}

func main() {
	_, err := connstring.ParseAndValidate(os.Getenv("mongodb://127.0.0.1:27017/retrospect?readPreference=primary&replicaSet=rs0"))
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017/retrospect?readPreference=primary&replicaSet=rs0"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	/*
	   List databases
	*/
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	database := client.Database("retrospect")

	collection := database.Collection("persons")

	stream, err := collection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		panic(err)
	}

	defer stream.Close(context.TODO())
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	routineCtx, _ := context.WithCancel(context.Background())

	go iterateChangeStream(routineCtx, waitGroup, stream)

	waitGroup.Wait()
}
