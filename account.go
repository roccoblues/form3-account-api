package form3

import (
	"fmt"
	"net/http"
	"time"
)

// Account represents a bank account that is registered with Form3.
type Account struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	OrganisationID string    `json:"organisation_id"`
	Version        int       `json:"version"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
	Attributes     struct {
		Country                 string   `json:"Country"`
		BaseCurrency            string   `json:"base_currency"`
		AccountNumber           string   `json:"account_number"`
		BankID                  string   `json:"bank_id"`
		BankIDCode              string   `json:"bank_id_code"`
		Bic                     string   `json:"bic"`
		IBAN                    string   `json:"iban"`
		Name                    []string `json:"name"`
		AlternativeNames        []string `json:"alternative_names"`
		AccountClassification   string   `json:"account_classification"`
		JointAccount            string   `json:"joint_account"`
		AccountMatchingOptOut   string   `json:"account_matching_opt_out"`
		SecondaryIdentification string   `json:"secondary_identification"`
		Switched                bool     `json:"switched"`
		Status                  string   `json:"status"`
	}
}

type accountResponse struct {
	Data  Account `json:"data"`
	Links Links   `json:"links"`
}

// GetAccount returns the account for the given ID.
func (c *Client) GetAccount(id string) (*Account, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/organisation/accounts/%s", c.APIBase, id), nil)
	if err != nil {
		return nil, err
	}

	accountResponse := &accountResponse{}
	if err := c.DoRequest(req, accountResponse); err != nil {
		return nil, err
	}

	return &accountResponse.Data, nil
}
