# form3-account-api

This repository contains a client for the Form3 fake account API.

## Development Notes

1. I opted to write integration tests only. They cover nearly all of the code paths. Ideally we should replace the `http.Client` in `Client` with an interface and write unit tests for `Client.DoRequest()` as this is doing quite some work.

## Submitted by

Dennis Schön (mail@dennis-schoen.de)
