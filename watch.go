package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

// Set up change stream watching based on Watch config
func setupCollectionWatch(watch Watch, database *mongo.Database, waitGroup *sync.WaitGroup) {
	collection := database.Collection(watch.Collection)

	stream, err := collection.Watch(
		context.Background(),
		mongo.Pipeline{},
		options.ChangeStream().SetFullDocument(options.UpdateLookup),
	)
	if err != nil {
		panic(err)
	}

	waitGroup.Add(1)

	go iterateChangeStream(waitGroup, stream, watch)
}

// Listens to incoming change events and performs actions with them
func iterateChangeStream(waitGroup *sync.WaitGroup, stream *mongo.ChangeStream, watch Watch) {
	defer func() {
		err := stream.Close(context.Background())
		if err != nil {
			log.Println("Error during change stream closing", err)
		}
	}()
	defer waitGroup.Done()

	for stream.Next(context.Background()) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}

		printJSON(data)
		text, err := renderTemplate(data, watch)

		if err != nil {
			log.Println("Error during template rendering", err)
		}

		err = notify(text, watch.NotifyHook)
		if err != nil {
			log.Println("Error during sending notification", err)
		}
	}
}
