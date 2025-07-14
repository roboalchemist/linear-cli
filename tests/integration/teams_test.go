package integration

import (
	"context"
	"testing"
	
	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/tests/testutils"
)

func TestListTeams(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	teams, err := client.GetTeams(ctx, 50, "", "")
	if err != nil {
		t.Fatalf("Failed to get teams: %v", err)
	}
	
	if len(teams.Nodes) == 0 {
		t.Skip("No teams available for testing")
	}
	
	t.Logf("✅ Found %d teams", len(teams.Nodes))
	
	// Verify team structure
	for _, team := range teams.Nodes {
		if team.ID == "" {
			t.Error("Expected team ID, got empty string")
		}
		if team.Name == "" {
			t.Error("Expected team name, got empty string")
		}
		if team.Key == "" {
			t.Error("Expected team key, got empty string")
		}
		
		t.Logf("  - %s (%s): %s", team.Name, team.Key, team.ID)
	}
}

func TestGetTeam(t *testing.T) {
	apiKey := testutils.SkipIfNoAuth(t)
	client := api.NewClient(apiKey)
	ctx := context.Background()
	
	// First get teams to have a valid ID
	teams, err := client.GetTeams(ctx, 1, "", "")
	if err != nil {
		t.Fatalf("Failed to get teams: %v", err)
	}
	
	if len(teams.Nodes) == 0 {
		t.Skip("No teams available for testing")
	}
	
	teamID := teams.Nodes[0].ID
	
	// Test getting specific team
	team, err := client.GetTeam(ctx, teamID)
	if err != nil {
		t.Fatalf("Failed to get team %s: %v", teamID, err)
	}
	
	if team.ID != teamID {
		t.Errorf("Expected team ID %s, got %s", teamID, team.ID)
	}
	
	t.Logf("✅ Successfully retrieved team: %s (%s)", team.Name, team.Key)
}