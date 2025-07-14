package integration

import (
	"context"
	"testing"
	
	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/tests/testutils"
)

func TestListUsers(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	users, err := client.GetUsers(ctx, 50, "", "")
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}
	
	if len(users.Nodes) == 0 {
		t.Skip("No users available for testing")
	}
	
	t.Logf("✅ Found %d users", len(users.Nodes))
	
	// Verify user structure
	for i, user := range users.Nodes {
		if user.ID == "" {
			t.Error("Expected user ID, got empty string")
		}
		if user.Email == "" && user.Name == "" {
			t.Error("Expected user to have either email or name")
		}
		
		// Only log first few to avoid spam
		if i < 5 {
			t.Logf("  - %s <%s>", user.Name, user.Email)
		}
	}
}

func TestGetCurrentUser(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	user, err := client.GetViewer(ctx)
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}
	
	if user.ID == "" {
		t.Error("Expected user ID, got empty string")
	}
	if user.Email == "" {
		t.Error("Expected user email, got empty string")
	}
	if user.Name == "" {
		t.Error("Expected user name, got empty string")
	}
	
	t.Logf("✅ Current user: %s <%s>", user.Name, user.Email)
}