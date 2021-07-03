package main

import (
	"bytes"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"text/template"
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
	defer func(stream *mongo.ChangeStream) {
		err := stream.Close(context.Background())
		if err != nil {
			log.Println("Error during change stream closing", err)
		}
	}(stream)
	defer waitGroup.Done()

	for stream.Next(context.Background()) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}

		t, err := template.New("todos").Parse(watch.Template)

		if err != nil {
			panic(err)
		}

		var tpl bytes.Buffer

		err = t.Execute(&tpl, data)
		if err != nil {
			panic(err)
		}

		err = notify(tpl.String(), watch.NotifyHook)
		if err != nil {
			log.Println("Error during change stream closing", err)
		}
	}
}
