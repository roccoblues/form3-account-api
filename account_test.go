package form3

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestAccountsCrud(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := newTestClient(t)

	// try to fetch account, it should not exist
	account, err := client.GetAccount(uuid.NewString())
	if err == nil {
		t.Fatalf("Client.GetAccount() expected error")
	}
	var e *HTTPError
	if errors.As(err, &e) {
		if e.StatusCode != http.StatusNotFound {
			t.Fatalf("Client.GetAccount() returned wrong status code. Expected %d, got %d", http.StatusNotFound, e.StatusCode)
		}
	} else {
		t.Fatalf("Client.GetAccount() returned error is not an HTTPError: (%T) %v", err, err)
	}
	if account != nil {
		t.Fatalf("Client.GetAccount() returned account, expected nil")
	}

	// create a new account
	attributes := AccountAttributes{
		Country:      "GB",
		BaseCurrency: "GBP",
		BankID:       "400300",
		BankIDCode:   "GBDSC",
		Bic:          "NWBKGB22",
	}
	accountID := uuid.NewString()
	organisationID := uuid.NewString()
	account, err = client.CreateAccount(accountID, organisationID, attributes)
	if err != nil {
		t.Fatalf("Client.CreateAccount() returned unexpected error: (%T) %v", err, err)
	}
	if account == nil {
		t.Fatalf("Client.CreateAccount() didn't return an account")
	}
	if account.ID != accountID {
		t.Errorf("Client.CreateAccount() wrong account id, expected %s, got %s", accountID, account.ID)
	}
	if account.OrganisationID != organisationID {
		t.Errorf("Client.CreateAccount() wrong organisation id, expected %s, got %s", organisationID, account.OrganisationID)
	}
	if !reflect.DeepEqual(account.Attributes, attributes) {
		t.Errorf("Client.CreateAccount() account attributes don't match. Got: %v, expected: %v", account.Attributes, attributes)
	}

	// list accounts
	accounts, err := client.ListAccounts()
	if err != nil {
		t.Fatalf("Client.ListAccounts() returned unexpected error: (%T) %v", err, err)
	}
	if len(accounts) != 1 {
		t.Errorf("Client.ListAccounts() returned wrong number of accounts. Expected 1, got %d", len(accounts))
	}
}
