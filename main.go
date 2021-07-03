package main

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Watch struct {
	Name       string `yaml:"name"`
	On         string `yaml:"on"`
	Collection string `yaml:"collection"`
}

type Config struct {
	MongodbUri string  `yaml:"mongodbUri"`
	Watches    []Watch `yaml:"watches"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := Config{}

	// Open config file
	file, err := ioutil.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	cfg, err := NewConfig("./config.yml")

	if err != nil {
		log.Fatal(err)
	}

	cs, err := connstring.ParseAndValidate(cfg.MongodbUri)
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongodbUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	database := client.Database(cs.Database)
	var waitGroup sync.WaitGroup

	for _, value := range cfg.Watches {
		collection := database.Collection(value.Collection)

		stream, err := collection.Watch(ctx, mongo.Pipeline{})
		if err != nil {
			panic(err)
		}

		waitGroup.Add(1)
		routineCtx, _ := context.WithCancel(context.Background())

		go iterateChangeStream(routineCtx, waitGroup, stream)
	}

	waitGroup.Wait()
}

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
