package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Check if expand or fields args have values and construct a url query including it
func AddExpandAndFieldsToURL(expand string, fields string) string {
	query := ""

	if len(expand) > 0 {
		query += fmt.Sprintf("?expand=%s", expand)

		if len(fields) > 0 {
			query += fmt.Sprintf("&fields=%s", fields)
		}
	} else if len(fields) > 0 {
		query += fmt.Sprintf("?fields=%s", fields)
	}

	return query
}

// Sends an HTTP request to the provided url
func SendHTTPRequest(method string, url string, headers map[string]string, options map[string]any) (http.Response, error) {

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
		return http.Response{}, fmt.Errorf("request-limit-reached|%s", url)
	}

	return *resp, nil
}

func SendAuthenticatedHTTPRequest(method string, url string, headers map[string]string, options map[string]any, token string) (http.Response, error) {
	headers["Authorization"] = token
	return SendHTTPRequest(method, url, headers, options)
}

func DecodePocketBaseRecord(response http.Response) map[string]any {
	record := map[string]any{}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&record)
	return record
}

func DecodePocketBaseListResponse(response http.Response) PocketBaseListResponse {
	record := PocketBaseListResponse{}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&record)
	return record
}

func DecodePocketBaseErrorResponse(response http.Response) PocketBaseErrorResponse {
	errRes := PocketBaseErrorResponse{}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&errRes)
	return errRes
}
