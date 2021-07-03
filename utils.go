package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Pretty-prints data as JSON document
func printJson(data interface{}) {
	text, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(text))
}
