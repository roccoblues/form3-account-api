package form3

import (
	"fmt"
	"net/http"
	"time"
)

// Account represents a bank account that is registered with Form3.
type Account struct {
	ID             string             `json:"id"`
	Type           string             `json:"type,omitempty"`
	OrganisationID string             `json:"organisation_id"`
	Version        int                `json:"version,omitempty"`
	CreatedOn      *time.Time         `json:"created_on,omitempty"`
	ModifiedOn     *time.Time         `json:"modified_on,omitempty"`
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
}

// AccountAttributes contains all attributes of an acccount.
//
// There are some differences between the documentation and the fake account
// implementation (https://github.com/form3tech-oss/interview-accountapi/issues/38).
// This API client currently only supports the non-deprecated working fields.
type AccountAttributes struct {
	Country       string `json:"Country,omitempty"`
	BaseCurrency  string `json:"base_currency,omitempty"`
	AccountNumber string `json:"account_number,omitempty"`
	BankID        string `json:"bank_id,omitempty"`
	BankIDCode    string `json:"bank_id_code,omitempty"`
	Bic           string `json:"bic,omitempty"`
	IBAN          string `json:"iban,omitempty"`
	// Name                    []string `json:"name,omitempty"`
	// AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string `json:"account_classification,omitempty"`
	JointAccount            bool   `json:"joint_account,omitempty"`
	AccountMatchingOptOut   bool   `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string `json:"secondary_identification,omitempty"`
	// Switched                bool   `json:"switched,omitempty"`
	Status string `json:"status,omitempty"`
}

// ListAccountsParams contains parameters to customize the ListAccounts call.
type ListAccountsParams struct {
	PageNumber int
	PageSize   int
}

type accountPayload struct {
	Data  Account `json:"data"`
	Links Links   `json:"links,omitempty"`
}

type accountsPayload struct {
	Data  []*Account `json:"data"`
	Links Links      `json:"links,omitempty"`
}

// GetAccount returns the account for the given ID.
func (c *Client) GetAccount(id string) (*Account, error) {
	req, err := c.NewRequest(http.MethodGet, fmt.Sprintf("/organisation/accounts/%s", id), nil)
	if err != nil {
		return nil, err
	}

	response := &accountPayload{}
	if err := c.DoRequest(req, response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// CreateAccount creates an account with the given attributes.
func (c *Client) CreateAccount(id, organisationID string, attributes *AccountAttributes) (*Account, error) {
	payload := &accountPayload{
		Data: Account{
			ID:             id,
			OrganisationID: organisationID,
			Type:           "accounts",
			Attributes:     attributes,
		},
	}
	req, err := c.NewRequest(http.MethodPost, "/organisation/accounts", payload)
	if err != nil {
		return nil, err
	}

	response := &accountPayload{}
	if err := c.DoRequest(req, response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// ListAccounts fetches accounts.
func (c *Client) ListAccounts(params *ListAccountsParams) ([]*Account, error) {
	path := fmt.Sprintf("/organisation/accounts")
	if params != nil {
		path = fmt.Sprintf("%s?page[number]=%d&page[size]=%d", path, params.PageNumber, params.PageSize)
	}
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &accountsPayload{}
	if err := c.DoRequest(req, response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// DeleteAccount deletes the account with the given ID and version.
func (c *Client) DeleteAccount(id string, version int) error {
	req, err := c.NewRequest(http.MethodDelete, fmt.Sprintf("/organisation/accounts/%s?version=%d", id, version), nil)
	if err != nil {
		return err
	}

	return c.DoRequest(req, nil)
}
