package services

import (
	"encoding/json"
	"fmt"
)

type AdminsAPI struct {
	BaseURL string `json:"baseURL"`
}

type AdminAuthQuery struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Fields   string `json:"fields"`
}

type AdminAuthResponse struct {
	Admin struct {
		Id      string `json:"id"`
		Created string `json:"created"`
		Updated string `json:"updated"`
		Email   string `json:"email"`
		Avatar  int    `json:"avatar"`
	} `json:"admin"`
	Token string `json:"token"`
}

// Authenticates the specified admin account information with the PocketBase server using the "/api/admins/auth-with-password" api route
func (api *AdminsAPI) AuthWithPassword(auth AdminAuthQuery) (authRes AdminAuthResponse, pbErr PocketBaseAPIError) {

	//Setup the request body
	options := map[string]interface{}{
		"identity": auth.Email,
		"password": auth.Password,
	}

	var apiURL string

	//Determine if we should add the "fields" query to the end of the URL
	if auth.Fields == "" {
		apiURL = fmt.Sprintf("%s/api/admins/auth-with-password", api.BaseURL)
	} else {
		apiURL = fmt.Sprintf("%s/api/admins/auth-with-password/?fields=%s", api.BaseURL, auth.Fields)
	}

	//Send the request to the API route
	resp, err := SendHTTPRequest("POST", apiURL, map[string]string{}, options)

	if ok := HandleError(err); !ok {
		return
	}

	defer resp.Body.Close()

	//If the response is successful send back a new AdminAuthResponse
	if resp.StatusCode == 200 {
		var body struct {
			Data AdminAuthResponse
		}

		err = json.NewDecoder(resp.Body).Decode(&body.Data)

		if ok := HandleError(err); !ok {
			return AdminAuthResponse{}, PocketBaseAPIError{}
		}

		return body.Data, PocketBaseAPIError{}

		//If the response is not successful send back a new PocketBaseAPIError
	} else {
		var body struct {
			Data PocketBaseAPIError
		}

		err = json.NewDecoder(resp.Body).Decode(&body.Data)

		if ok := HandleError(err); !ok {
			return AdminAuthResponse{}, PocketBaseAPIError{}
		}

		return AdminAuthResponse{}, body.Data
	}
}

// Returns a new AdminAuthResponse for an already authenticated admin account using the "/api/admins/auth-refresh" api route
func (api *AdminsAPI) AuthRefresh(token string) (authRes AdminAuthResponse, pbErr PocketBaseAPIError) {

	//Setup the request headers
	headers := map[string]string{
		"Authorization": token,
	}

	//Send the request to the API route
	resp, err := SendHTTPRequest("POST", fmt.Sprintf("%s/api/admins/auth-refresh", api.BaseURL), headers, map[string]interface{}{})

	if ok := HandleError(err); !ok {
		return
	}

	defer resp.Body.Close()

	//If the response is successful send back a new AdminAuthResponse
	if resp.StatusCode == 200 {
		var body AdminAuthResponse

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return AdminAuthResponse{}, PocketBaseAPIError{}
		}

		return body, PocketBaseAPIError{}

		//If the response is not successful send back a new PocketBaseAPIError
	} else {
		var body PocketBaseAPIError

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return AdminAuthResponse{}, PocketBaseAPIError{}
		}

		return AdminAuthResponse{}, body
	}
}

// Sends a password reset email to the specified email using the "/api/admins/request-password-reset" api route
func (api *AdminsAPI) RequestPasswordReset(email string) (sent bool) {

	//Setup the request headers
	options := map[string]interface{}{
		"email": email,
	}

	//Send the request to the API route
	resp, err := SendHTTPRequest("POST", fmt.Sprintf("%s/api/admins/request-password-reset", api.BaseURL), map[string]string{}, options)

	if ok := HandleError(err); !ok {
		return
	}

	defer resp.Body.Close()

	//If the response is successful return true
	if resp.StatusCode == 204 {
		return true

		//If the response is not successful send back a new PocketBaseAPIError
	} else {
		var body PocketBaseAPIError

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return false
		}

		return false
	}
}
