package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync"
)

func main() {
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := NewConfig(cfgPath)

	if err != nil {
		log.Println(err)
	}

	database, client := GetDatabase(cfg.MongodbUri)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Println("Error during disconnecting from database", err)
		}
	}(client, context.Background())

	var waitGroup sync.WaitGroup

	for _, value := range cfg.Watches {
		setupWatch(value, database, &waitGroup)
	}
	log.Println("All watcher are set up")

	waitGroup.Wait()
}
