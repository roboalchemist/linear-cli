package integration

import (
	"context"
	"os"
	"testing"
	
	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/tests/testutils"
)

func TestAuthValidation(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	
	// Create a new client with the test API key
	client := api.NewClient(apiKey)
	
	// Test viewer query
	viewer, err := client.GetViewer(context.Background())
	if err != nil {
		t.Fatalf("Failed to get viewer: %v", err)
	}
	
	// Verify we got valid user data
	if viewer.ID == "" {
		t.Error("Expected viewer ID, got empty string")
	}
	if viewer.Email == "" {
		t.Error("Expected viewer email, got empty string")
	}
	if viewer.Name == "" {
		t.Error("Expected viewer name, got empty string")
	}
	
	t.Logf("✅ Successfully authenticated as: %s <%s>", viewer.Name, viewer.Email)
}

func TestInvalidAuth(t *testing.T) {
	// Only run if we have a test environment
	if !testutils.IsIntegrationTest() {
		t.Skip("Skipping integration test - no test API key")
	}
	
	// Test with invalid API key
	client := api.NewClient("invalid-api-key")
	
	_, err := client.GetViewer(context.Background())
	if err == nil {
		t.Fatal("Expected error with invalid API key, got nil")
	}
	
	t.Logf("✅ Invalid auth correctly rejected: %v", err)
}