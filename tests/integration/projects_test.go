package integration

import (
	"context"
	"testing"
	
	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/tests/testutils"
)

func TestListProjects(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	// Get teams first to filter projects
	teams, err := client.GetTeams(ctx, 1, "", "")
	if err != nil {
		t.Fatalf("Failed to get teams: %v", err)
	}
	
	if len(teams.Nodes) == 0 {
		t.Skip("No teams available for testing")
	}
	
	// Test listing all projects
	projects, err := client.GetProjects(ctx, map[string]interface{}{}, 50, "", "")
	if err != nil {
		t.Fatalf("Failed to get projects: %v", err)
	}
	
	t.Logf("✅ Found %d total projects", len(projects.Nodes))
	
	// Test filtering by team
	teamFilter := map[string]interface{}{
		"teams": map[string]interface{}{
			"in": []string{teams.Nodes[0].ID},
		},
	}
	teamProjects, err := client.GetProjects(ctx, teamFilter, 50, "", "")
	if err != nil {
		t.Fatalf("Failed to get team projects: %v", err)
	}
	
	t.Logf("✅ Found %d projects for team %s", len(teamProjects.Nodes), teams.Nodes[0].Name)
	
	// Verify project structure if we have any
	if len(projects.Nodes) > 0 {
		project := projects.Nodes[0]
		if project.ID == "" {
			t.Error("Expected project ID, got empty string")
		}
		if project.Name == "" {
			t.Error("Expected project name, got empty string")
		}
		
		t.Logf("  Example project: %s (State: %s)", project.Name, project.State)
	}
}

func TestGetProject(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	// First get projects to have a valid ID
	projects, err := client.GetProjects(ctx, map[string]interface{}{}, 1, "", "")
	if err != nil {
		t.Fatalf("Failed to get projects: %v", err)
	}
	
	if len(projects.Nodes) == 0 {
		t.Skip("No projects available for testing")
	}
	
	projectID := projects.Nodes[0].ID
	
	// Test getting specific project
	project, err := client.GetProject(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to get project %s: %v", projectID, err)
	}
	
	if project.ID != projectID {
		t.Errorf("Expected project ID %s, got %s", projectID, project.ID)
	}
	
	t.Logf("✅ Successfully retrieved project: %s", project.Name)
	if project.Description != "" {
		t.Logf("  Description: %s", project.Description)
	}
	if project.StartDate != "" || project.TargetDate != "" {
		t.Logf("  Timeline: %s → %s", project.StartDate, project.TargetDate)
	}
}