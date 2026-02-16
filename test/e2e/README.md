# E2E Tests

End-to-end tests for the IP Verifier API.

## Running E2E Tests

### 1. Start the server
```bash
go run cmd/ip-verifier-api/main.go
```

### 2. Run E2E tests (in another terminal)
```bash
# Run all e2e tests
go test -tags=e2e ./test/e2e/

# Run with verbose output
go test -tags=e2e -v ./test/e2e/

# Run specific test
go test -tags=e2e -v ./test/e2e/ -run TestVerifyIP_ValidIPInAllowedCountries
```

## Running All Tests

```bash
# Run only unit tests (default - no build tags)
go test ./...

# Run all tests including e2e
go test -tags=e2e ./...
```

## Why Build Tags?

The `//go:build e2e` tag prevents these tests from running during normal `go test ./...` because:
- E2E tests require a running server
- They're slower than unit tests
- They should run separately in CI/CD pipelines
