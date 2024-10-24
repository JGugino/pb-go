package services

import (
	"encoding/json"
	"fmt"
)

type HealthAPI struct {
	BaseURL string `json:"baseURL"`
}

type HealthCheckResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CanBackup bool `json:"canBackup"`
	} `json:"data"`
}

func (api *HealthAPI) CheckAPIHealth() HealthCheckResponse {
	headers := map[string]string{}

	options := map[string]interface{}{}

	resp, err := SendHTTPRequest("GET", fmt.Sprintf("%s/api/health", api.BaseURL), headers, options)

	if ok := HandleError(err); !ok {
		return HealthCheckResponse{}
	}

	var body HealthCheckResponse

	err = json.NewDecoder(resp.Body).Decode(&body)

	if ok := HandleError(err); !ok {
		return HealthCheckResponse{}
	}

	return body
}
