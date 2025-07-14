package unit

import (
	"context"
	"fmt"
	"strings"
	"testing"
	
	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/tests/testutils"
)

func TestClientWithMockServer(t *testing.T) {
	// Create mock responses
	responses := map[string]testutils.MockResponse{
		"query Me": {
			Response: map[string]interface{}{
				"viewer": testutils.SampleUser(),
			},
		},
		"query Issues": {
			Response: map[string]interface{}{
				"issues": map[string]interface{}{
					"nodes": []interface{}{
						testutils.SampleIssue(),
					},
				},
			},
		},
	}
	
	// Create mock server
	server := testutils.MockLinearServer(t, responses)
	defer server.Close()
	
	// Create client pointing to mock server
	client := api.NewClientWithURL(server.URL, "test-token")
	ctx := context.Background()
	
	// Test viewer query
	t.Run("GetViewer", func(t *testing.T) {
		viewer, err := client.GetViewer(ctx)
		if err != nil {
			t.Fatalf("Failed to get viewer: %v", err)
		}
		
		if viewer.Name != "Test User" {
			t.Errorf("Expected name 'Test User', got '%s'", viewer.Name)
		}
		if viewer.Email != "test@example.com" {
			t.Errorf("Expected email 'test@example.com', got '%s'", viewer.Email)
		}
	})
	
	// Test issues query
	t.Run("GetIssues", func(t *testing.T) {
		issues, err := client.GetIssues(ctx, map[string]interface{}{}, 10, "", "")
		if err != nil {
			t.Fatalf("Failed to get issues: %v", err)
		}
		
		if len(issues.Nodes) != 1 {
			t.Fatalf("Expected 1 issue, got %d", len(issues.Nodes))
		}
		
		issue := issues.Nodes[0]
		if issue.Identifier != "LIN-123" {
			t.Errorf("Expected identifier 'LIN-123', got '%s'", issue.Identifier)
		}
		if issue.Title != "Test Issue" {
			t.Errorf("Expected title 'Test Issue', got '%s'", issue.Title)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	// Create mock server that returns errors
	responses := map[string]testutils.MockResponse{
		"Me": {
			Error: fmt.Errorf("Unauthorized"),
		},
	}
	
	server := testutils.MockLinearServer(t, responses)
	defer server.Close()
	
	client := api.NewClientWithURL(server.URL, "invalid-token")
	ctx := context.Background()
	
	_, err := client.GetViewer(ctx)
	if err == nil {
		t.Fatal("Expected error for unauthorized request, got nil")
	}
	
	if !strings.Contains(err.Error(), "Unauthorized") {
		t.Errorf("Expected error to contain 'Unauthorized', got: %v", err)
	}
}