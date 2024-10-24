package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type PocketBaseAPIError struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    map[string]ErrorInfo `json:"data"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Sends an HTTP request to the provided url
func SendHTTPRequest(method string, url string, headers map[string]string, options map[string]interface{}) (http.Response, error) {

	//Marshal the provided into JSON for the body of the request.
	body, err := json.Marshal(options)
	if err != nil {
		fmt.Println(err)
		return http.Response{}, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return http.Response{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return http.Response{}, err
	}

	if resp.StatusCode == 429 {
		return http.Response{}, errors.New(fmt.Sprintf("request-limit-reached|%s", url))
	}

	return *resp, nil
}

func HandleError(err error) (cont bool) {
	if err != nil {
		log.Println("[ERROR] ", err)
		return false
	}

	return true
}

func HandlePocketBaseError(err PocketBaseAPIError) (cont bool) {
	if err.Message == "" {
		return true
	}

	return HandleError(errors.New(err.Message))
}
