package integration

import (
	"context"
	"testing"
	
	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/tests/testutils"
)

func TestListIssues(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	tests := []struct {
		name   string
		filter map[string]interface{}
		desc   string
		limit  int
	}{
		{
			name:   "List all issues",
			filter: map[string]interface{}{},
			limit:  5,
			desc:   "Should list up to 5 issues",
		},
		{
			name: "Filter by state",
			filter: map[string]interface{}{
				"state": map[string]interface{}{
					"name": map[string]interface{}{
						"in": []string{"In Progress", "Todo"},
					},
				},
			},
			limit: 5,
			desc:  "Should list issues in specific states",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := client.GetIssues(ctx, tt.filter, tt.limit, "", "updatedAt")
			if err != nil {
				t.Fatalf("Failed to get issues: %v", err)
			}
			
			t.Logf("✅ %s - Found %d issues", tt.desc, len(issues.Nodes))
			
			// Verify issue structure if we got any
			if len(issues.Nodes) > 0 {
				issue := issues.Nodes[0]
				if issue.ID == "" {
					t.Error("Expected issue ID, got empty string")
				}
				if issue.Title == "" {
					t.Error("Expected issue title, got empty string")
				}
				if issue.Identifier == "" {
					t.Error("Expected issue identifier, got empty string")
				}
			}
		})
	}
}

func TestGetIssue(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	// First, get an issue to test with
	issues, err := client.GetIssues(ctx, map[string]interface{}{}, 1, "", "updatedAt")
	if err != nil {
		t.Fatalf("Failed to get issues: %v", err)
	}
	
	if len(issues.Nodes) == 0 {
		t.Skip("No issues available for testing")
	}
	
	issueID := issues.Nodes[0].ID
	
	// Test getting specific issue
	issue, err := client.GetIssue(ctx, issueID)
	if err != nil {
		t.Fatalf("Failed to get issue %s: %v", issueID, err)
	}
	
	// Verify we got the right issue
	if issue.ID != issueID {
		t.Errorf("Expected issue ID %s, got %s", issueID, issue.ID)
	}
	
	t.Logf("✅ Successfully retrieved issue: %s - %s", issue.Identifier, issue.Title)
}

func TestSearchIssues(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	// Search for issues with a common term
	filter := map[string]interface{}{
		"searchableContent": map[string]string{
			"contains": "test",
		},
	}
	
	issues, err := client.GetIssues(ctx, filter, 5, "", "updatedAt")
	if err != nil {
		t.Fatalf("Failed to search issues: %v", err)
	}
	
	t.Logf("✅ Search found %d issues matching 'test'", len(issues.Nodes))
}