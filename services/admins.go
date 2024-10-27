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

type AdminRecord struct {
	Id      string `json:"id"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Email   string `json:"email"`
	Avatar  int    `json:"avatar"`
}

type AdminAuthResponse struct {
	Admin AdminRecord `json:"admin"`
	Token string      `json:"token"`
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

// Confirms the password reset with the token sent to the users email using the "/api/admins/confirm-password-reset" api route
func (api *AdminsAPI) ConfirmPasswordReset(token, password, passwordConfirm string) (valid bool) {

	//Setup the request headers
	options := map[string]interface{}{
		"token":           token,
		"password":        password,
		"passwordConfirm": passwordConfirm,
	}

	//Send the request to the API route
	resp, err := SendHTTPRequest("POST", fmt.Sprintf("%s/api/admins/confirm-password-reset", api.BaseURL), map[string]string{}, options)

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

// Get a single admin account based on their ID using the "/api/admins/{id}" api route
func (api *AdminsAPI) GetAdmin(token, id, fields string) (admin AdminRecord, pbError PocketBaseAPIError) {

	//Setup the request body
	headers := map[string]string{
		"Authorization": token,
	}

	var apiURL string

	//Determine if we should add the "fields" query to the end of the URL
	if fields == "" {
		apiURL = fmt.Sprintf("%s/api/admins/%s", api.BaseURL, id)
	} else {
		apiURL = fmt.Sprintf("%s/api/admins/%s/?fields=%s", api.BaseURL, id, fields)
	}

	//Send the request to the API route
	resp, err := SendHTTPRequest("GET", apiURL, headers, map[string]interface{}{})

	if ok := HandleError(err); !ok {
		return
	}

	//If the response is successful send back a new AdminAuthResponse
	if resp.StatusCode == 200 {
		var body AdminRecord

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return AdminRecord{}, PocketBaseAPIError{}
		}

		return body, PocketBaseAPIError{}

		//If the response is not successful send back a new PocketBaseAPIError
	} else {
		var body PocketBaseAPIError

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return AdminRecord{}, PocketBaseAPIError{}
		}

		return AdminRecord{}, body
	}
}

// Get a list of admin accounts using the "/api/admins/" api route and the provided query options
func (api *AdminsAPI) GetList(page, perPage int, sort, filter, fields string, skipTotal bool, token string) (admin PaginatedPocketBaseResponse, pbError PocketBaseAPIError) {

	//Setup the request body
	headers := map[string]string{
		"Authorization": token,
	}

	queries := PocketBaseQueryOptions{
		Page:      page,
		PerPage:   perPage,
		Sort:      sort,
		Filter:    filter,
		Fields:    fields,
		SkipTotal: skipTotal,
	}

	apiURL := BuildPocketBaseURLWithQueries(fmt.Sprintf("%s/api/admins/?", api.BaseURL), queries)

	//Send the request to the API route
	resp, err := SendHTTPRequest("GET", apiURL, headers, map[string]interface{}{})

	if ok := HandleError(err); !ok {
		return
	}

	//If the response is successful send back a new AdminAuthResponse
	if resp.StatusCode == 200 {
		var body PaginatedPocketBaseResponse

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return PaginatedPocketBaseResponse{}, PocketBaseAPIError{}
		}

		return body, PocketBaseAPIError{}

		//If the response is not successful send back a new PocketBaseAPIError
	} else {
		var body PocketBaseAPIError

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return PaginatedPocketBaseResponse{}, PocketBaseAPIError{}
		}

		return PaginatedPocketBaseResponse{}, body
	}
}

// Get a list of admin accounts using the "/api/admins/" api route and the provided query options
func (api *AdminsAPI) CreateAdmin(id, email, password, passwordConfirm string, avatar int, token, fields string) (admin AdminRecord, pbError PocketBaseAPIError) {

	//Setup the request headers
	headers := map[string]string{
		"Authorization": token,
	}
	//Setup the request body
	options := map[string]interface{}{
		"id":              id,
		"email":           email,
		"password":        password,
		"passwordConfirm": passwordConfirm,
		"avatar":          avatar,
	}

	queries := PocketBaseQueryOptions{
		Page:      0,
		PerPage:   0,
		Sort:      "",
		Filter:    "",
		Fields:    fields,
		SkipTotal: false,
	}

	apiURL := BuildPocketBaseURLWithQueries(fmt.Sprintf("%s/api/admins/?", api.BaseURL), queries)

	//Send the request to the API route
	resp, err := SendHTTPRequest("POST", apiURL, headers, options)

	if ok := HandleError(err); !ok {
		return
	}

	//If the response is successful send back a new AdminAuthResponse
	if resp.StatusCode == 200 {
		var body AdminRecord

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return AdminRecord{}, PocketBaseAPIError{}
		}

		return body, PocketBaseAPIError{}

		//If the response is not successful send back a new PocketBaseAPIError
	} else {
		var body PocketBaseAPIError

		err = json.NewDecoder(resp.Body).Decode(&body)

		if ok := HandleError(err); !ok {
			return AdminRecord{}, PocketBaseAPIError{}
		}

		return AdminRecord{}, body
	}
}
