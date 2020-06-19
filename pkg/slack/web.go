package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
)

// See https://api.slack.com/web

// Get is used to query the Slack Web API
func Get(ctx context.Context, token, apiMethod, query string, response interface{}) error {
	url := SlackEndpoint + apiMethod + "?" + query

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// unmarshal the response
	return json.NewDecoder(resp.Body).Decode(&response)
}

// Post is used to invoke a Slack Web API method by posting a JSON payload.
func Post(ctx context.Context, token, apiMethod string, request interface{}) (*StandardResponse, error) {
	url := SlackEndpoint + apiMethod

	m, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// unmarshal the response
	var apiResponse StandardResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)

	return &apiResponse, err
}

// CustomPost is used to invoke a Slack Web API method that respondes with a non-standard payload. See https://api.slack.com/web
func CustomPost(ctx context.Context, token, apiMethod string, request, response interface{}) error {
	url := SlackEndpoint + apiMethod

	m, err := json.Marshal(&request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	// post the request to Slack
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// unmarshal the response
	err = json.NewDecoder(resp.Body).Decode(&response)

	return err
}
