package form3

import (
	"errors"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			assert.Equal(account.Version, tt.want.Version)
			assert.False(account.CreatedOn.IsZero())
			assert.False(account.ModifiedOn.IsZero())

			if diff := deep.Equal(account.Attributes, tt.want.Attributes); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestClient_DeleteAccount(t *testing.T) {
	client := newTestClient(t)
	truncateAccounts(client, t)

	account, err := client.CreateAccount(testAccountID, testOrganisationID, &testAccountAttributes)
	require.Nil(t, err)

	tests := []struct {
		name              string
		id                string
		version           int
		wantErr           bool
		wantHTTPErr       bool
		httpErrStatusCode int
	}{
		// The provided fake account API returns 204 instead of 404 for known existing records.
		// See: https://github.com/form3tech-oss/interview-accountapi/issues/30
		// {
		// 	name:    "non-existing account",
		// 	id:      uuid.NewString(),
		// 	version: 0,
		// 	wantErr: true,
		// },
		{
			name:              "existing account, but wrong version",
			id:                account.ID,
			version:           account.Version + 1,
			wantErr:           true,
			wantHTTPErr:       true,
			httpErrStatusCode: 404, // according to the documentation a 409 should be returned
		},
		{
			name:    "existing account",
			id:      account.ID,
			version: account.Version,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteAccount(tt.id, tt.version)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantHTTPErr {
				var e *HTTPError
				if errors.As(err, &e) {
					assert.Equal(t, tt.httpErrStatusCode, e.StatusCode)
				} else {
					t.Errorf("Client.DeleteAccount() wrong error type (%T) %v", err, err)
				}
			}
		})
	}
}

func TestClient_GetAccount(t *testing.T) {
	client := newTestClient(t)
	truncateAccounts(client, t)

	testAccount, err := client.CreateAccount(testAccountID, testOrganisationID, &testAccountAttributes)
	require.Nil(t, err)

	assert := assert.New(t)

	tests := []struct {
		name              string
		id                string
		want              *Account
		wantErr           bool
		wantHTTPErr       bool
		httpErrStatusCode int
	}{
		{
			name:              "non-existing account",
			id:                uuid.NewString(),
			wantErr:           true,
			wantHTTPErr:       true,
			httpErrStatusCode: 404,
		},
		{
			name:    "existing account",
			id:      testAccount.ID,
			want:    testAccount,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := client.GetAccount(tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantHTTPErr {
				var e *HTTPError
				if errors.As(err, &e) {
					assert.Equal(tt.httpErrStatusCode, e.StatusCode)
				} else {
					t.Errorf("Client.GetAccount() wrong error type (%T) %v", err, err)
				}
			}

			if tt.want == nil {
				return
			}

			assert.Equal(account.ID, tt.want.ID)
			assert.Equal(account.OrganisationID, tt.want.OrganisationID)
			assert.Equal(account.Version, tt.want.Version)
			assert.False(account.CreatedOn.IsZero())
			assert.False(account.ModifiedOn.IsZero())

			if diff := deep.Equal(account.Attributes, tt.want.Attributes); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestClient_ListAccounts(t *testing.T) {
	assert := assert.New(t)

	client := newTestClient(t)
	truncateAccounts(client, t)

	testAccount1, err := client.CreateAccount(uuid.NewString(), uuid.NewString(), &testAccountAttributes)
	require.Nil(t, err)

	testAccount2, err := client.CreateAccount(uuid.NewString(), uuid.NewString(), &testAccountAttributes)
	require.Nil(t, err)

	tests := []struct {
		name string
		want []*Account
	}{
		{
			name: "fetch all accounts",
			want: []*Account{
				testAccount1,
				testAccount2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accounts, _ := client.ListAccounts()
			require.Equal(t, len(accounts), len(tt.want))

			for i, a := range tt.want {
				b := accounts[i]

				assert.Equal(a.ID, b.ID)
				assert.Equal(a.OrganisationID, b.OrganisationID)
				assert.Equal(a.Version, b.Version)

				if diff := deep.Equal(a.Attributes, b.Attributes); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}
