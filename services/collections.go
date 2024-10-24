package services

import (
	"fmt"
	"log"
)

type CollectionsAPI struct {
	BaseURL string `json:"baseURL"`
}

type ListCollectionQuery struct {
	Page      int    `json:"page"`
	PerPage   int    `json:"perPage"`
	Sort      string `json:"sort"`
	Filter    string `json:"filter"`
	Fields    string `json:"fields"`
	SkipTotal bool   `json:"skipTotal"`
}

func (api *CollectionsAPI) ListCollections(authorization string, listCollectionQuery ListCollectionQuery) {

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": authorization,
	}

	options := map[string]interface{}{
		"page":      listCollectionQuery.Page,
		"perPage":   listCollectionQuery.PerPage,
		"sort":      listCollectionQuery.Sort,
		"filter":    listCollectionQuery.Filter,
		"fields":    listCollectionQuery.Fields,
		"skipTotal": listCollectionQuery.SkipTotal,
	}

	resp, err := SendHTTPRequest("POST", fmt.Sprintf("%s/api/collections", api.BaseURL), headers, options)

	if err != nil {
		log.Println("Collection Error:", err.Error())
		return
	}

	fmt.Println(resp)

}
