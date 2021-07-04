package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"text/template"
)

//nolint:deadcode,unused // Pretty-prints data as JSON document
func printJSON(data interface{}) {
	text, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(text))
}

// Renders Watch template and returns its content
func renderTemplate(data interface{}, watch Watch) (string, error) {
	t, err := template.New(watch.Name).Parse(watch.Template)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	err = t.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
