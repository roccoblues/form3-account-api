package form3

import "testing"

func newTestClient(t *testing.T) *Client {
	client, err := NewClient("http://localhost:8080/v1")
	if err != nil {
		t.Fatal(err)
	}

	// ensure an empty test setup
	accounts, err := client.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range accounts {
		if err := client.DeleteAccount(a.ID, a.Version); err != nil {
			t.Fatal(err)
		}
	}

	return client
}
