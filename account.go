package form3

import (
	"fmt"
	"net/http"
	"time"
)

// Account represents a bank account that is registered with Form3.
type Account struct {
	ID             string            `json:"id"`
	Type           string            `json:"type,omitempty"`
	OrganisationID string            `json:"organisation_id"`
	Version        int               `json:"version,omitempty"`
	CreatedOn      *time.Time        `json:"created_on,omitempty"`
	ModifiedOn     *time.Time        `json:"modified_on,omitempty"`
	Attributes     AccountAttributes `json:"attributes,omitempty"`
}

// AccountAttributes contains all attributes of an acccount.
type AccountAttributes struct {
	Country                 string   `json:"Country,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	IBAN                    string   `json:"iban,omitempty"`
	Name                    []string `json:"name,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	JointAccount            string   `json:"joint_account,omitempty"`
	AccountMatchingOptOut   string   `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
	Status                  string   `json:"status,omitempty"`
}

type accountPayload struct {
	Data  Account `json:"data"`
	Links Links   `json:"links,omitempty"`
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
func (c *Client) CreateAccount(id, organisationID string, attributes AccountAttributes) (*Account, error) {
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
