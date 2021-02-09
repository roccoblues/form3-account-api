# form3-account-api

This repository contains a client for the Form3 fake account API.

## Usage


### Initialize a Client

```Go
import "github.com/roccoblues/form3-account-api"

func main() {
  client, err := form3.NewClient("http://localhost:8080/v1")
  ...
}
```

The `http.Client` used to make the actual HTTP requests can be configured.

```Go
httpClient := &http.Client{Timeout: 10}
options := []form3.ClientOption(form3.WithHTTPClient(httpClient))
client, err := form3.NewClient("http://localhost:8080/v1", options...)
```

### Error handling

All API methods return an error in case something went wrong. If the error is based on an HTTP response a `form3.HTTPError` is returned. It contains the actual HTTP response and [additional information](https://api-docs.form3.tech/api.html#introduction-and-api-conventions-errors-and-status-codes) (if provided).

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
API_BASE="http://somewhere.else" go test -v
```

## Development Notes

1. I've choosen to write integration tests only. They cover nearly all of the code paths. Ideally we should replace the `http.Client` in `Client` with an interface and write unit tests for `Client.DoRequest()`, as this is doing quite some work.

2. There are some differences between the documentation and the fake account implementation (see: https://github.com/form3tech-oss/interview-accountapi/issues/38).
This API client currently only supports the non-deprecated working fields.

## Submitted by

Dennis Schön (mail@dennis-schoen.de)
