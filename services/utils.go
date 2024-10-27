package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type PocketBaseAPIError struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    map[string]ErrorInfo `json:"data"`
}

type PaginatedItems interface {
	AdminRecord
}

type PaginatedPocketBaseResponse struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"perPage"`
	TotalItems int           `json:"totalItems"`
	TotalPages int           `json:"totalPages"`
	Items      []interface{} `json:"items"`
}

type PocketBaseQueryOptions struct {
	Page      int    `json:"page"`
	PerPage   int    `json:"perPage"`
	Sort      string `json:"sort"`
	Filter    string `json:"filter"`
	Fields    string `json:"fields"`
	SkipTotal bool   `json:"skipTotal"`
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

func BuildPocketBaseURLWithQueries(startingUrl string, queries PocketBaseQueryOptions) string {
	builtUrl := startingUrl
	queryAdded := false

	//Determine if we should add the "page" query to the end of the URL
	if queries.Page > 0 {
		builtUrl, queryAdded = AddQueryToURL(queryAdded, builtUrl, fmt.Sprintf("page=%s", string(queries.Page)))
	}

	//Determine if we should add the "perPage" query to the end of the URL
	if queries.PerPage > 0 {
		builtUrl, queryAdded = AddQueryToURL(queryAdded, builtUrl, fmt.Sprintf("perPage=%s", string(queries.PerPage)))
	}

	//Determine if we should add the "sort" query to the end of the URL
	if queries.Sort != "" {
		builtUrl, queryAdded = AddQueryToURL(queryAdded, builtUrl, fmt.Sprintf("sort=%s", queries.Sort))
	}

	//Determine if we should add the "fields" query to the end of the URL
	if queries.Filter != "" {
		builtUrl, queryAdded = AddQueryToURL(queryAdded, builtUrl, fmt.Sprintf("filter=(%s)", queries.Filter))
	}

	//Determine if we should add the "fields" query to the end of the URL
	if queries.Fields != "" {
		builtUrl, queryAdded = AddQueryToURL(queryAdded, builtUrl, fmt.Sprintf("fields=%s", queries.Fields))
	}

	//Add the "skipTotal" query to the end of the URL
	builtUrl, queryAdded = AddQueryToURL(queryAdded, builtUrl, fmt.Sprintf("skipTotal=%s", strconv.FormatBool(queries.SkipTotal)))

	return builtUrl
}

func AddQueryToURL(queryAdded bool, currentQueryString, queryAddition string) (string, bool) {
	if queryAdded {
		return fmt.Sprintf("%s&%s", currentQueryString, queryAddition), queryAdded
	}

	queryAdded = true

	return fmt.Sprintf("%s%s", currentQueryString, queryAddition), queryAdded
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

	log.Println(err)
	return HandleError(errors.New(err.Message))
}
