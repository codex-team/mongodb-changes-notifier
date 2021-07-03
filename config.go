package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
