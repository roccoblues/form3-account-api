package form3

import (
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testAccountID = uuid.NewString()
var testOrganisationID = uuid.NewString()
var testAccountAttributes = AccountAttributes{
	Country:                 "DE",
	BaseCurrency:            "EUR",
	AccountNumber:           "123",
	BankID:                  "FOOBAR42",
	BankIDCode:              "DEBLZ",
	Bic:                     "XXXXXX42",
	IBAN:                    "DE42100020003000400099",
	AccountClassification:   "Personal",
	JointAccount:            true,
	AccountMatchingOptOut:   true,
	SecondaryIdentification: "666",
}

func TestClient_CreateAccount(t *testing.T) {
	client := newTestClient(t)
	truncateAccounts(client, t)

	assert := assert.New(t)

	tests := []struct {
		name              string
		id                string
		organisationID    string
		attributes        *AccountAttributes
		want              *Account
		wantErr           bool
		wantHTTPErr       bool
		httpErrStatusCode int
		httpErrCode       int
		httpErrMessage    string
	}{
		{
			name:              "empty attributes",
			id:                testAccountID,
			organisationID:    testOrganisationID,
			wantErr:           true,
			wantHTTPErr:       true,
			httpErrStatusCode: 400,
			httpErrCode:       0,
			httpErrMessage:    "validation failure list:\nvalidation failure list:\nattributes in body is required",
		},
		{
			name:           "full attributes",
			id:             testAccountID,
			organisationID: testOrganisationID,
			attributes:     &testAccountAttributes,
			want: &Account{
				ID:             testAccountID,
				OrganisationID: testOrganisationID,
				Attributes:     &testAccountAttributes,
			},
		},
		{
			name:              "duplicate account id",
			id:                testAccountID,
			organisationID:    uuid.NewString(),
			attributes:        &testAccountAttributes,
			wantErr:           true,
			wantHTTPErr:       true,
			httpErrStatusCode: 409,
			httpErrCode:       0,
			httpErrMessage:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := client.CreateAccount(tt.id, tt.organisationID, tt.attributes)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantHTTPErr {
				var e *HTTPError
				if errors.As(err, &e) {
					assert.Equal(tt.httpErrStatusCode, e.StatusCode)
					assert.Equal(tt.httpErrCode, e.ErrorCode)
					assert.Equal(tt.httpErrMessage, e.ErrorMessage)
				} else {
					t.Errorf("Client.CreateAccount() wrong error type (%T) %v", err, err)
				}
			}

			if tt.want == nil {
				return
			}

			assert.Equal(account.ID, tt.id)
			assert.Equal(account.OrganisationID, tt.organisationID)
			assert.False(account.CreatedOn.IsZero())
			assert.False(account.ModifiedOn.IsZero())

			if diff := deep.Equal(account.Attributes, tt.want.Attributes); diff != nil {
				t.Error(diff)
			}
		})
	}
}
