package form3

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultAPIBase = "http://localhost:8080/v1"

func newTestClient(t *testing.T) *Client {
	apiBase := os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = defaultAPIBase
	}

	client, err := NewClient(apiBase)
	require.Nil(t, err)

	return client
}

func truncateAccounts(client *Client, t *testing.T) {
	accounts, err := client.ListAccounts()
	require.Nil(t, err)

	for _, a := range accounts {
		err := client.DeleteAccount(a.ID, a.Version)
		require.Nil(t, err)
	}
}
