package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
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

//nolint:deadcode,unused // Pretty-prints data as YAML document
func printYAML(data interface{}) {
	text, err := yaml.Marshal(data)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(text))
}

// Renders Watch template and returns its content
func renderTemplate(data interface{}, watch Watch) (string, error) {
	var sprigFuncMap = sprig.GenericFuncMap()

	sprigFuncMap["toYaml"] = func(v interface{}) string {
		output, _ := yaml.Marshal(v)
		return string(output)
	}

	t, err := template.New(watch.Name).Funcs(sprigFuncMap).Parse(watch.Template)

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
