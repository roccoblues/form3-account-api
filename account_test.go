package form3

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestAccountsCrud(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := newTestClient(t)

	// 404 Not Found
	account, err := client.GetAccount(uuid.NewString())
	if err == nil {
		t.Fatalf("Client.GetAccount() expected error")
	}
	var e *ErrorResponse
	if errors.As(err, &e) {
		if e.StatusCode != http.StatusNotFound {
			t.Fatalf("Client.GetAccount() returned wrong status code. Expected %d, got %d", http.StatusNotFound, e.StatusCode)
		}
	} else {
		t.Fatalf("Client.GetAccount() returned error is not an ErrorResponse: (%T) %v", err, err)
	}
	if account != nil {
		t.Fatalf("Client.GetAccount() returned account, expected nil")
	}

	// Create account

}
