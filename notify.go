package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Sends notifications to Telegram via webhook
func notify(text string, hook string) error {
	data := url.Values{}
	data.Set("message", text)

	_, err := MakeHTTPRequest("POST", hook, []byte(data.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	})
	if err != nil {
		return fmt.Errorf("webhook error: %v", err)
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

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
