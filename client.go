package form3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client represents a Form3 REST API Client.
type Client struct {
	httpClient HTTPClient
	baseURL    string
}

// HTTPClient models the http client interface.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Links contains links to related endpoints.
type Links struct {
	Self  string `json:"self,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// HTTPError contains additional information about a failed http request.
type HTTPError struct {
	Response     *http.Response `json:"-"`
	StatusCode   int            `json:"-"` // convinience accessor for Response.StatusCode
	Status       string         `json:"-"` // convinience accessor for Response.Status
	ErrorCode    int            `json:"error_code,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
}

// Error returns a string representation of the HTTPError.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d %s: %d %s", e.Response.StatusCode, e.Response.Status, e.ErrorCode, e.ErrorMessage)
}

// NewClient returns a new Client struct.
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is required")
	}

	c := &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
	for _, option := range options {
		option(c)
	}
	return c, nil
}

// ClientOption sets some additional options on a client.
type ClientOption func(*Client)

// WithHTTPClient sets the http client used to make the actual requests.
func WithHTTPClient(httpClient HTTPClient) ClientOption {
	return func(c *Client) { c.httpClient = httpClient }
}

// NewRequest returns a http.Request for the given path.
// If a payload is provided it will get JSON encoded.
func (c *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	if payload == nil {
		return http.NewRequest(method, fmt.Sprintf("%s%s", c.baseURL, path), nil)
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)
	return http.NewRequest(method, fmt.Sprintf("%s%s", c.baseURL, path), buf)
}

// DoRequest makes a request to the API and unmarshales the response into v.
func (c *Client) DoRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/vnd.api+json")

	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// successful delete requests have no return value
	if req.Method == http.MethodDelete && resp.StatusCode == http.StatusNoContent {
		return nil
	}

	// only 200 and 201 responses contain a body that can be unmashalled into v
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	errResp := &HTTPError{
		Response:   resp,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}

	// 400 responses contain additional information
	if resp.StatusCode == http.StatusBadRequest {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if len(data) > 0 {
			if err := json.Unmarshal(data, errResp); err != nil {
				return err
			}
		}
	}

	return errResp
}
