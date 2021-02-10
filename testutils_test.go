package form3

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultBaseURL = "http://localhost:8080/v1"

func newTestClient(t *testing.T) *Client {
	baseURL := os.Getenv("API_BASE")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client, err := NewClient(baseURL)
	require.Nil(t, err)

	return client
}

func truncateAccounts(client *Client, t *testing.T) {
	resp, err := client.ListAccounts(&ListAccountsParams{PageSize: 9999})
	require.Nil(t, err)

	for _, a := range resp.Accounts() {
		err := client.DeleteAccount(a.ID, a.Version)
		require.Nil(t, err)
	}
}
