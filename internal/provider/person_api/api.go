package person_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	auth0ClientID     string
	auth0ClientSecret string
	auth0Audience     string
	auth0Endpoint     string
	auth0Scopes       []string
	httpClient        *http.Client
	personEndpoint    string

	auth0AccessToken string
}

func NewClient(auth0ClientID string, auth0ClientSecret string, auth0Audience string, auth0Endpoint string, auth0Scopes []string, personEndpoint string) *Client {
	c := &Client{
		auth0ClientID:     auth0ClientID,
		auth0ClientSecret: auth0ClientSecret,
		auth0Audience:     auth0Audience,
		auth0Endpoint:     auth0Endpoint,
		auth0Scopes:       auth0Scopes,
		httpClient:        &http.Client{},
		personEndpoint:    personEndpoint,
	}

	return c
}

func (client *Client) GetAccessToken(ctx context.Context) error {
	oauth2_config := clientcredentials.Config{
		ClientID:       client.auth0ClientID,
		ClientSecret:   client.auth0ClientSecret,
		EndpointParams: url.Values{"audience": {client.auth0Audience}},
		Scopes:         client.auth0Scopes,
		TokenURL:       client.auth0Endpoint,
	}

	oauth_token, err := oauth2_config.Token(ctx)
	tflog.Info(ctx, fmt.Sprintf("HTTP Request: %#v", oauth_token))

	if err == nil {
		client.auth0AccessToken = oauth_token.AccessToken
	}

	return err
}

func (client *Client) GetPersonByEmail(ctx context.Context, email string) (*Person, error) {
	person := Person{}

	httpReq, err := http.NewRequest("GET", client.personEndpoint+"/v2/user/primary_email/"+email, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+client.auth0AccessToken)

	httpResp, err := client.httpClient.Do(httpReq)
	tflog.Info(ctx, fmt.Sprintf("HTTP Request: %#v", httpReq))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode >= 400 {
		return nil, fmt.Errorf("Person API responded with status code %d", httpResp.StatusCode)
	}

	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &person)
	if err != nil {
		return nil, err
	}

	// Convert map keys into a list of strings
	keys := make([]string, 0, len(person.AccessInformation.Mozilliansorg.Values))
	for key := range person.AccessInformation.Mozilliansorg.Values {
		keys = append(keys, key)
	}
	person.AccessInformation.Mozilliansorg.List = keys

	return &person, nil
}
