# linctl Tests 🧪

This directory contains the test suite for linctl, organized into unit and integration tests.

## Test Structure

```
tests/
├── integration/     # End-to-end tests with real Linear API
├── unit/           # Unit tests with mocked dependencies
└── testutils/      # Shared test utilities
```

## Running Tests

### All Tests
```bash
make test
```

### Unit Tests Only
```bash
go test ./tests/unit/...
```

### Integration Tests
Integration tests require a Linear API key. Set it as an environment variable:

```bash
export LINEAR_TEST_API_KEY="your-test-api-key"
go test ./tests/integration/...
```

⚠️ **Important**: Integration tests use READ-ONLY operations to ensure safety when using real API tokens.

## Test Categories

### Unit Tests (`tests/unit/`)
- **api_client_test.go**: Tests API client with mocked HTTP server
- **commands_test.go**: Tests command structure and configuration

### Integration Tests (`tests/integration/`)
- **auth_test.go**: Tests authentication flow
- **issues_test.go**: Tests issue listing and retrieval
- **teams_test.go**: Tests team operations
- **users_test.go**: Tests user queries
- **projects_test.go**: Tests project management

### Test Utilities (`tests/testutils/`)
- **config.go**: Test configuration and environment helpers
- **mock_client.go**: HTTP server mocking for unit tests

## Writing New Tests

### Unit Test Example
```go
func TestNewFeature(t *testing.T) {
    // Create mock server
    server := testutils.MockLinearServer(t, responses)
    defer server.Close()
    
    // Test your feature
    client := api.NewClientWithURL(server.URL, "test-token")
    // ... test logic
}
```

### Integration Test Example
```go
func TestRealAPICall(t *testing.T) {
    apiKey := testutils.SkipIfNoAuth(t)  // Skip if no API key
    client := api.NewClient(apiKey)
    
    // Make real API calls (read-only!)
    result, err := client.GetSomething()
    // ... assertions
}
```

## CI/CD Configuration

For CI environments, set the `LINEAR_TEST_API_KEY` as a secret to enable integration tests.

## Safety Guidelines

1. Integration tests should NEVER modify data
2. Use test-specific Linear workspace if possible
3. Always check for `LINEAR_TEST_API_KEY` before running integration tests
4. Mock all write operations in unit tests