package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Watch describes structure for collection watch task
type Watch struct {
	Name       string `yaml:"name"`
	On         string `yaml:"on"`
	Collection string `yaml:"collection"`
	NotifyHook string `yaml:"notify_hook"`
	Template   string `yaml:"template"`
}

// Config describes application configuration structure
type Config struct {
	MongodbUri string  `yaml:"mongodb_uri"`
	NotifyHook string  `yaml:"notify_hook"`
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
		setupCollectionWatch(value, database, ctx, &waitGroup)
	}

	waitGroup.Wait()
}

// Pretty-prints data as JSON document
func printJson(data interface{}) {
	text, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(text))
}

func setupCollectionWatch(watch Watch, database *mongo.Database, ctx context.Context, waitGroup *sync.WaitGroup) {
	collection := database.Collection(watch.Collection)

	stream, err := collection.Watch(ctx, mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		panic(err)
	}

	waitGroup.Add(1)
	routineCtx, _ := context.WithCancel(context.Background())

	go iterateChangeStream(routineCtx, waitGroup, stream, watch)
}

func iterateChangeStream(routineCtx context.Context, waitGroup *sync.WaitGroup, stream *mongo.ChangeStream, watch Watch) {
	defer stream.Close(routineCtx)
	defer waitGroup.Done()
	for stream.Next(routineCtx) {
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

		notify(tpl.String(), watch.NotifyHook)
		printJson(data)
	}
}

func notify(text string, hook string) error {
	data := url.Values{}
	data.Set("message", text)

	_, err := MakeHTTPRequest("POST", hook, []byte(data.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		return fmt.Errorf("Webhook error: %v", err)
	}

	return nil
}

// MakeHTTPRequest - make HTTP request with specified method, body, URL and headers
func MakeHTTPRequest(method string, url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	r := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
