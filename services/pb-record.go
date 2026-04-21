package services

import (
	"errors"
	"fmt"
	"net/http"
)

type PBRecord struct {
	BaseURL string
}

// ### CREATE RECORDS ###
func (record *PBRecord) CreateAuthRecord(collection string, email string, password string, passwordConfirm string, token string) (map[string]any, error) {
	apiURL := fmt.Sprintf("%s/api/collections/%s/records", record.BaseURL, collection)

	body := map[string]any{
		"email":           email,
		"password":        password,
		"passwordConfirm": passwordConfirm,
	}

	res, err := SendAuthenticatedHTTPRequest("POST", apiURL, map[string]string{}, body, token)

	if err != nil {
		return map[string]any{}, err
	}

	status := res.StatusCode

	if status == http.StatusOK {
		createdRecord := DecodePocketBaseRecord(res)

		return createdRecord, nil
	}

	errRes := DecodePocketBaseErrorResponse(res)
	return map[string]any{}, errors.New(errRes.Message)
}

func (record *PBRecord) CreateNewRecord(collection string, token string, data map[string]any) (map[string]any, error) {
	apiURL := fmt.Sprintf("%s/api/collections/%s/records", record.BaseURL, collection)

	res, err := SendAuthenticatedHTTPRequest("POST", apiURL, map[string]string{}, data, token)

	if err != nil {
		return map[string]any{}, err
	}

	status := res.StatusCode

	if status == http.StatusOK {
		createdRecord := DecodePocketBaseRecord(res)
		return createdRecord, nil
	}

	errRes := DecodePocketBaseErrorResponse(res)
	return map[string]any{}, errors.New(errRes.Message)
}

//### VIEW RECORDS ###

func (record *PBRecord) ListRecords(collection string, token string, queryOptions PocketBaseListOptions) (PocketBaseListResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/records", record.BaseURL, collection)

	queryString, hasOptions := ConstructQueryStringForAPI(queryOptions)

	if hasOptions {
		apiUrl += queryString
	}

	res, err := SendAuthenticatedHTTPRequest("GET", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return PocketBaseListResponse{}, err
	}

	status := res.StatusCode

	if status == http.StatusOK {
		record := DecodePocketBaseListResponse(res)

		return record, nil
	}

	errRes := DecodePocketBaseErrorResponse(res)
	return PocketBaseListResponse{}, errors.New(errRes.Message)
}

func (record *PBRecord) ViewRecord(collection string, recordId string, token string) (map[string]any, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/records/%s", record.BaseURL, collection, recordId)

	res, err := SendAuthenticatedHTTPRequest("GET", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return map[string]any{}, err
	}

	status := res.StatusCode

	if status == http.StatusOK {
		record := DecodePocketBaseRecord(res)

		return record, nil
	}

	errRes := DecodePocketBaseErrorResponse(res)
	return map[string]any{}, errors.New(errRes.Message)
}

// ### DELETE RECORDS ###
func (record *PBRecord) DeleteRecord(collection string, recordId string, token string) (bool, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/records/%s", record.BaseURL, collection, recordId)

	res, err := SendAuthenticatedHTTPRequest("DELETE", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return false, err
	}

	status := res.StatusCode

	if status == http.StatusNoContent {
		return true, nil
	}

	pbErr := DecodePocketBaseErrorResponse(res)

	return false, errors.New(pbErr.Message)
}

// ### UPDATE RECORDS ###
func (record *PBRecord) UpdateRecord(collection string, recordId string, token string, updatedData map[string]any) (map[string]any, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/records/%s", record.BaseURL, collection, recordId)

	res, err := SendAuthenticatedHTTPRequest("PATCH", apiUrl, map[string]string{}, updatedData, token)

	if err != nil {
		return map[string]any{}, err
	}

	status := res.StatusCode

	if status == http.StatusOK {
		decodedRecord := DecodePocketBaseRecord(res)
		return decodedRecord, nil
	}

	pbErr := DecodePocketBaseErrorResponse(res)

	return map[string]any{}, errors.New(pbErr.Message)
}
