package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CollectionType int

const (
	BaseCollection CollectionType = iota
	AuthCollection
	ViewCollection
)

func (c CollectionType) String() string {
	switch c {
	case 0:
		return "base"
	case 1:
		return "auth"
	case 2:
		return "view"
	default:
		return "base"
	}
}

type CollectionAuthAlerts struct {
	Enabled       bool                    `json:"enabled"`
	EmailTemplate CollectionEmailTemplate `json:"emailTemplate"`
}

type OAuthProviderOptions struct {
	Name         string `json:"name"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AuthURL      string `json:"authURL"`
	TokenURL     string `json:"tokenURL"`
	UserInfoURL  string `json:"userInfoURL"`
	DisplayName  string `json:"displayName"`
	PKCE         string `json:"pkce"`
	Extra        string `json:"extra"`
}

type CollectionOAuthOptions struct {
	Enabled      bool `json:"enabled"`
	MappedFields struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatarURL"`
	} `json:"mappedFields"`
	Providers OAuthProviderOptions
}

type CollectionTokenOptions struct {
	Duration int    `json:"duration"`
	Secret   string `json:""`
}

type CollectionEmailTemplate struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type CollectionPasswordAuth struct {
	Enabled        bool     `json:"enabled"`
	IdentityFields []string `json:"identityFields"`
}
type CollectionMFA struct {
	Enabled  bool `json:"enabled"`
	Duration int  `json:"duration"`
	Rule     bool `json:"rule"`
}
type CollectionOTP struct {
	Enabled       bool                    `json:"enabled"`
	Duration      int                     `json:"duration"`
	Length        int                     `json:"length"`
	EmailTemplate CollectionEmailTemplate `json:"emailTemplate"`
}

type CollectionOptions struct {
	Name       string               `json:"name"`
	Type       CollectionType       `json:"type"`
	Fields     []map[string]any     `json:"fields"`
	System     bool                 `json:"system"`
	ListRule   string               `json:"listRule"`
	ViewRule   string               `json:"viewRule"`
	CreateRule string               `json:"createRule"`
	UpdateRule string               `json:"updateRule"`
	DeleteRule string               `json:"deleteRule"`
	ViewQuery  string               `json:"viewQuery"`
	AuthAlerts CollectionAuthAlerts `json:"authAlerts"`

	//Auth Options
	OAuth2       CollectionOAuthOptions `json:"oauth2"`
	PasswordAuth CollectionPasswordAuth `json:"passwordAuth"`
	MFA          CollectionMFA          `json:"mfa"`
	OTP          CollectionOTP          `json:"otp"`

	//Auth Tokens
	AuthToken          CollectionTokenOptions `json:"authToken"`
	PasswordResetToken CollectionTokenOptions `json:"passwordResetToken"`
	EmailChangeToken   CollectionTokenOptions `json:"emailChangeToken"`
	VerificationToken  CollectionTokenOptions `json:"verificationToken"`
	FileToken          CollectionTokenOptions `json:"fileToken"`

	//Email Templates
	VerificationTemplate       CollectionEmailTemplate `json:"verificationTemplate"`
	ResetPasswordTemplate      CollectionEmailTemplate `json:"resetPasswordTemplate"`
	ConfirmEmailChangeTemplate CollectionEmailTemplate `json:"confirmEmailChangeTemplate"`
}

type PocketBaseCollectionResponse struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Fields     []map[string]any `json:"fields"`
	System     bool             `json:"system"`
	ListRule   string           `json:"listRule"`
	ViewRule   string           `json:"viewRule"`
	CreateRule string           `json:"createRule"`
	UpdateRule string           `json:"updateRule"`
	DeleteRule string           `json:"deleteRule"`
}

type PBCollection struct {
	BaseURL string
}

// ### CREATE COLLECTION ###
func (collection *PBCollection) CreateNewCollection(token string, options CollectionOptions) (PocketBaseCollectionResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/collections", collection.BaseURL)

	collectionOptions := map[string]any{
		"name":   options.Name,
		"type":   options.Type.String(),
		"fields": options.Fields,
		"system": options.System,
	}

	switch options.Type.String() {
	//Create Auth Collection
	case AuthCollection.String():
		collectionOptions["createRule"] = options.CreateRule
		collectionOptions["updateRule"] = options.UpdateRule
		collectionOptions["deleteRule"] = options.DeleteRule

		if len(options.PasswordAuth.IdentityFields) > 0 {
			collectionOptions["passwordAuth"] = options.PasswordAuth
		}

	//Create View Collection
	case ViewCollection.String():
		collectionOptions["listRule"] = options.ListRule
		collectionOptions["viewRule"] = options.ViewRule
		collectionOptions["viewQuery"] = options.ViewQuery
	}

	res, err := SendAuthenticatedHTTPRequest("POST", apiUrl, map[string]string{}, collectionOptions, token)

	if err != nil {
		return PocketBaseCollectionResponse{}, err
	}

	status := res.StatusCode

	if status != http.StatusOK {
		pbErr := DecodePocketBaseErrorResponse(res)

		return PocketBaseCollectionResponse{}, errors.New(pbErr.Message)
	}

	collectionRes := PocketBaseCollectionResponse{}
	json.NewDecoder(res.Body).Decode(&collectionRes)
	defer res.Body.Close()

	return collectionRes, nil

}

// ### UPDATE COLLECTION ###
func (collection *PBCollection) UpdateCollection(token string, desiredCollection string) (map[string]any, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s", collection.BaseURL, desiredCollection)
	res, err := SendAuthenticatedHTTPRequest("PATCH", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return map[string]any{}, err
	}

	fmt.Println(res)
	return map[string]any{}, nil
}

// ### VIEW COLLECTION ###

func (collection *PBCollection) ScaffoldCollections(token string) (map[string]any, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/meta/scaffolds", collection.BaseURL)
	res, err := SendAuthenticatedHTTPRequest("GET", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return map[string]any{}, err
	}

	if res.StatusCode != http.StatusOK {
		pbErr := DecodePocketBaseErrorResponse(res)

		return map[string]any{}, errors.New(pbErr.Message)
	}

	scaffoldRes := map[string]any{}
	json.NewDecoder(res.Body).Decode(&scaffoldRes)
	defer res.Body.Close()

	return scaffoldRes, nil
}

func (collection *PBCollection) ViewCollection(token string, desiredCollection string) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s", collection.BaseURL, desiredCollection)
	res, err := SendAuthenticatedHTTPRequest("GET", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return
	}

	fmt.Println(res)

}

func (collection *PBCollection) ListCollections(token string) {
	apiUrl := fmt.Sprintf("%s/api/collections", collection.BaseURL)
	res, err := SendAuthenticatedHTTPRequest("GET", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return
	}

	fmt.Println(res)

}

// ### DELETE COLLECTION ###
func (collection *PBCollection) DeleteCollection(token string, desiredCollection string) (bool, error) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s", collection.BaseURL, desiredCollection)
	res, err := SendAuthenticatedHTTPRequest("DELETE", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusNoContent {
		pbErr := DecodePocketBaseErrorResponse(res)

		return false, errors.New(pbErr.Message)
	}

	return true, nil
}

func (collection *PBCollection) TruncateCollection(token string, desiredCollection string) {
	apiUrl := fmt.Sprintf("%s/api/collections/%s/truncate", collection.BaseURL, desiredCollection)
	res, err := SendAuthenticatedHTTPRequest("DELETE", apiUrl, map[string]string{}, map[string]any{}, token)

	if err != nil {
		return
	}

	fmt.Println(res)
}
