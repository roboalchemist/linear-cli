package testutils

import (
	"os"
	"testing"
)

// GetTestAPIKey retrieves the test API key from environment
// This allows safe testing without hardcoding credentials
func GetTestAPIKey(t *testing.T) string {
	apiKey := os.Getenv("LINEAR_TEST_API_KEY")
	if apiKey == "" {
		t.Skip("LINEAR_TEST_API_KEY not set, skipping integration test")
	}
	return apiKey
}

// SkipIfNoAuth skips the test if no test API key is available
func SkipIfNoAuth(t *testing.T) string {
	return GetTestAPIKey(t)
}

// IsIntegrationTest checks if integration tests should run
func IsIntegrationTest() bool {
	return os.Getenv("LINEAR_TEST_API_KEY") != ""
}