package form3

import "testing"

func newTestClient(t *testing.T) *Client {
	client, err := NewClient("http://localhost:8080/v1")
	if err != nil {
		t.Fatal(err)
	}
	return client
}
