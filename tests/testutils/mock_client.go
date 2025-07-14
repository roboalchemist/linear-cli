package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// MockResponse represents a mocked GraphQL response
type MockResponse struct {
	Query    string
	Response interface{}
	Error    error
}

// MockLinearServer creates a test server that mocks Linear API responses
func MockLinearServer(t *testing.T, responses map[string]MockResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}
		
		// Find matching response
		for queryPrefix, mockResp := range responses {
			if contains(req.Query, queryPrefix) {
				if mockResp.Error != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"errors": []map[string]string{
							{"message": mockResp.Error.Error()},
						},
					})
					return
				}
				
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"data": mockResp.Response,
				})
				return
			}
		}
		
		// No matching response found
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]string{
				{"message": fmt.Sprintf("No mock response for query: %s", req.Query)},
			},
		})
	}))
}

func contains(s, substr string) bool {
	// Normalize whitespace for comparison
	normalizedS := strings.ReplaceAll(strings.TrimSpace(s), "\t", " ")
	normalizedS = strings.ReplaceAll(normalizedS, "\n", " ")
	normalizedSubstr := strings.ReplaceAll(strings.TrimSpace(substr), "\t", " ")
	normalizedSubstr = strings.ReplaceAll(normalizedSubstr, "\n", " ")
	
	return strings.Contains(normalizedS, normalizedSubstr)
}

// SampleIssue creates a sample issue for testing
func SampleIssue() map[string]interface{} {
	return map[string]interface{}{
		"id":          "ISSUE-123",
		"identifier":  "LIN-123",
		"title":       "Test Issue",
		"description": "This is a test issue",
		"state": map[string]interface{}{
			"id":   "STATE-1",
			"name": "In Progress",
		},
		"assignee": map[string]interface{}{
			"id":    "USER-1",
			"name":  "Test User",
			"email": "test@example.com",
		},
		"team": map[string]interface{}{
			"id":   "TEAM-1",
			"key":  "TEST",
			"name": "Test Team",
		},
		"priority":  1,
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-02T00:00:00Z",
	}
}

// SampleTeam creates a sample team for testing
func SampleTeam() map[string]interface{} {
	return map[string]interface{}{
		"id":          "TEAM-1",
		"key":         "TEST",
		"name":        "Test Team",
		"description": "Test team description",
	}
}

// SampleUser creates a sample user for testing
func SampleUser() map[string]interface{} {
	return map[string]interface{}{
		"id":          "USER-1",
		"name":        "Test User",
		"email":       "test@example.com",
		"displayName": "Test User",
		"active":      true,
	}
}

// SampleProject creates a sample project for testing
func SampleProject() map[string]interface{} {
	return map[string]interface{}{
		"id":          "PROJECT-1",
		"name":        "Test Project",
		"description": "Test project description",
		"state":       "started",
		"startDate":   "2024-01-01",
		"targetDate":  "2024-12-31",
		"team": map[string]interface{}{
			"id":   "TEAM-1",
			"name": "Test Team",
		},
	}
}