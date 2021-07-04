package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Watch describes structure for collection watch task
type Watch struct {
	Name        string   `yaml:"name"`
	EventTypes  []string `yaml:"event_types"`
	Collections []string `yaml:"collections"`
	NotifyHook  string   `yaml:"notify_hook"`
	Template    string   `yaml:"template"`
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

// ValidateConfigPath just makes sure, that the path provided is a file, that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
