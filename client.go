package form3

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client represents a Form3 REST API Client.
type Client struct {
	Client  *http.Client
	APIBase string
}

// Links contains links to related endpoints.
type Links struct {
	Self  string `json:"self"`
	First string `json:"first"`
	Last  string `json:"last"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

// ErrorResponse contains additional information about a failed request.
type ErrorResponse struct {
	Response     *http.Response `json:"-"`
	StatusCode   int            `json:"-"`
	Status       string         `json:"-"`
	ErrorCode    string         `json:"error_code"`
	ErrorMessage string         `json:"error_message"`
}

// Error returns a string represantation of the ErrorResponse.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d %s: %s %s", e.Response.StatusCode, e.Response.Status, e.ErrorCode, e.ErrorMessage)
}

// NewClient returns a new Client struct.
func NewClient(APIBase string) (*Client, error) {
	if APIBase == "" {
		return nil, errors.New("APIBase is required")
	}

	return &Client{
		Client:  &http.Client{},
		APIBase: APIBase,
	}, nil
}

// DoRequest makes a request to the API and unmarshales the response into v.
func (c *Client) DoRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Accept", "application/vnd.api+json")

	if req.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/vnd.api+json")
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// only 200 and 201 responses contain a body that can be unmashalled into v
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	errResp := &ErrorResponse{
		Response:   resp,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}

	// 400 responses contain additional information
	if resp.StatusCode == http.StatusBadRequest {
		data, err := ioutil.ReadAll(resp.Body)

		if err != nil && len(data) > 0 {
			if err := json.Unmarshal(data, errResp); err != nil {
				return err
			}
		}
	}

	return errResp
}
