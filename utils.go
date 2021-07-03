package main

import (
	"encoding/json"
	"fmt"
	"log"
)

//nolint:deadcode,unused // Pretty-prints data as JSON document
func printJSON(data interface{}) {
	text, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(text))
}
