# form3-account-api

This repository contains a client for the Form3 [interview account API](https://github.com/form3tech-oss/interview-accountapi).

## Usage

```Go
import "github.com/roccoblues/form3-account-api"

func main() {
  // error handling has been omitted for brevity.

  client, err := form3.NewClient("http://localhost:8080")

  // Create
  account, err := client.CreateAccount(organisationID, attributes)

  // Fetch
  account, err := client.GetAccount(id)

  // Delete
  err := client.DeleteAccount(id, version)

  // List
  res := client.ListAccounts(nil)
  for res.NextPage() {
    accounts, err := res.Accounts()
  }
}
```

The `http.Client` used to make the actual HTTP requests can be changed. It just needs to fullfil the `form3.HTTPClient` interface.

```Go
httpClient := &http.Client{Timeout: time.Second * 10}
options := []form3.ClientOption{form3.WithHTTPClient(httpClient)}
client, err := form3.NewClient("http://localhost:8080", options...)
```

### Error handling

All API methods return an error in case something went wrong. If the error is based on an HTTP response a `form3.HTTPError` is returned. It contains the actual `http.Response` and [additional information](https://api-docs.form3.tech/api.html#introduction-and-api-conventions-errors-and-status-codes) if provided.

```Go
account, err := client.CreateAccount(id, organisationId, attributes)
if err != nil {
  var e *form3.HTTPError
  if errors.As(err, &e) {
    fmt.Printf("status_code=%d error_code=%d error_message=%s", e.StatusCode, e.ErrorCode, e.ErrorMessage)
  }
  ...
}
```

## Tests

The tests can be run with `docker-compose up` in the base directory.

They can also be run standalone against a custom API Endpoint with:

```Bash
API_BASE="http://somewhere.else" go test
```

## Development Notes

1. I've chosen to write integration tests only. They cover most of the code and all of the happy paths. For a production-ready solution unit tests to cover the error cases could be added.

2. UUIDs are handled as simple strings. I chose to rely on the server to validate the correct format. This makes the client easier to use. Also the underlying format can be changed without changing the library interface.

3. There are some differences in the attributes between the documentation and the fake account implementation (see: https://github.com/form3tech-oss/interview-accountapi/issues/38).
This API client currently only supports the non-deprecated working attributes.

4. Another (known) issue is that the delete endpoint currently always returns 204, even for non-existing accounts. (see: https://github.com/form3tech-oss/interview-accountapi/issues/30)
