package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Package-level state set in TestMain
var (
	binaryPath string
	teamKey    string
	teamUUID   string
	testPrefix string

	// Cross-subtest shared IDs
	testProjectID string
	testLabelID   string
	testLabelName string
	testIssueID   string // identifier like ROB-123
	testIssueUUID string
	testViewID    string
	testCycleID   string
)

// runCLI executes the binary with given args and returns stdout, stderr, exit code.
func runCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run command %v: %v", args, err)
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

// runCLISuccess runs the binary and fails the test on non-zero exit.
func runCLISuccess(t *testing.T, args ...string) string {
	t.Helper()
	stdout, stderr, exitCode := runCLI(t, args...)
	if exitCode != 0 {
		t.Fatalf("command %v failed (exit %d)\nstdout: %s\nstderr: %s", args, exitCode, stdout, stderr)
	}
	return stdout
}

// runCLIFail runs the binary and fails the test if exit IS zero.
func runCLIFail(t *testing.T, args ...string) string {
	t.Helper()
	stdout, stderr, exitCode := runCLI(t, args...)
	if exitCode == 0 {
		t.Fatalf("command %v succeeded but expected failure\nstdout: %s\nstderr: %s", args, stdout, stderr)
	}
	// Return combined output since errors may be on stdout (cobra) or stderr
	return stdout + stderr
}

// parseJSONArray parses a JSON array string into []map[string]interface{}.
func parseJSONArray(t *testing.T, jsonStr string) []map[string]interface{} {
	t.Helper()
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		t.Fatalf("failed to parse JSON array: %v\njson: %s", err, truncate(jsonStr, 500))
	}
	return arr
}

// parseJSONObject parses a JSON object string into map[string]interface{}.
func parseJSONObject(t *testing.T, jsonStr string) map[string]interface{} {
	t.Helper()
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		t.Fatalf("failed to parse JSON object: %v\njson: %s", err, truncate(jsonStr, 500))
	}
	return obj
}

// extractField parses JSON (object or array-first-element) and returns a top-level field as string.
func extractField(t *testing.T, jsonStr string, field string) string {
	t.Helper()
	trimmed := strings.TrimSpace(jsonStr)

	var obj map[string]interface{}
	if strings.HasPrefix(trimmed, "[") {
		arr := parseJSONArray(t, trimmed)
		if len(arr) == 0 {
			t.Fatalf("extractField: empty JSON array, looking for field %q", field)
		}
		obj = arr[0]
	} else {
		obj = parseJSONObject(t, trimmed)
	}

	val, ok := obj[field]
	if !ok {
		t.Fatalf("extractField: field %q not found in JSON", field)
	}
	return fmt.Sprintf("%v", val)
}

// extractID is shorthand for extractField(t, jsonStr, "id").
func extractID(t *testing.T, jsonStr string) string {
	t.Helper()
	return extractField(t, jsonStr, "id")
}

// findByField finds the first element in a JSON array where field==value.
func findByField(t *testing.T, jsonStr, field, value string) map[string]interface{} {
	t.Helper()
	arr := parseJSONArray(t, jsonStr)
	for _, obj := range arr {
		if fmt.Sprintf("%v", obj[field]) == value {
			return obj
		}
	}
	t.Fatalf("findByField: no element with %s=%s found in array of %d elements", field, value, len(arr))
	return nil
}

// jsonArrayLen returns the length of a JSON array.
func jsonArrayLen(t *testing.T, jsonStr string) int {
	t.Helper()
	arr := parseJSONArray(t, jsonStr)
	return len(arr)
}

// truncate limits a string to maxLen chars for error messages.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// assertContains fails if substr is not found in s.
func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("expected output to contain %q, got: %s", substr, truncate(s, 300))
	}
}

// assertNotEmpty fails if s is empty or whitespace-only.
func assertNotEmpty(t *testing.T, s string) {
	t.Helper()
	if strings.TrimSpace(s) == "" {
		t.Errorf("expected non-empty output")
	}
}

func TestMain(m *testing.M) {
	// Build binary
	fmt.Println("Building linear-cli test binary...")
	build := exec.Command("go", "build", "-o", "linear-cli.test", ".")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to build binary: %v\n", err)
		os.Exit(1)
	}

	wd, _ := os.Getwd()
	binaryPath = wd + "/linear-cli.test"
	testPrefix = fmt.Sprintf("crud-test-%d", time.Now().Unix())

	// Verify auth
	fmt.Println("Verifying authentication...")
	cmd := exec.Command(binaryPath, "auth", "status", "--json")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: not authenticated. Run 'linear-cli auth' first.\n")
		os.Exit(1)
	}
	if !strings.Contains(string(out), "true") {
		fmt.Fprintf(os.Stderr, "FATAL: auth status check failed: %s\n", string(out))
		os.Exit(1)
	}

	// Discover team
	fmt.Println("Discovering team...")
	cmd = exec.Command(binaryPath, "team", "list", "--json")
	out, err = cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to list teams: %v\n", err)
		os.Exit(1)
	}
	var teams []map[string]interface{}
	if err := json.Unmarshal(out, &teams); err != nil || len(teams) == 0 {
		fmt.Fprintf(os.Stderr, "FATAL: no teams found or JSON parse error: %v\n", err)
		os.Exit(1)
	}
	teamKey = fmt.Sprintf("%v", teams[0]["key"])
	teamUUID = fmt.Sprintf("%v", teams[0]["id"])

	fmt.Printf("Test config: team=%s (%s), prefix=%s\n", teamKey, teamUUID, testPrefix)

	exitCode := m.Run()

	// Cleanup binary
	os.Remove(binaryPath)
	os.Exit(exitCode)
}

// =============================================================================
// TestCRUD - ordered subtests
// =============================================================================

func TestCRUD(t *testing.T) {
	// Ordered subtests - some produce IDs consumed by later ones.
	// Go runs subtests sequentially within a parent test.
	//
	// Shared resource cleanup is registered at THIS level so resources
	// persist across all subtests, not just the one that created them.

	// Register cleanup for shared resources at parent level.
	// These run in LIFO order after all subtests complete.
	t.Cleanup(func() {
		// Clean up issues (archive)
		if testIssueID != "" {
			runCLI(t, "issue", "archive", testIssueID)
		}
		// Clean up cycle
		if testCycleID != "" {
			runCLI(t, "cycle", "archive", testCycleID)
		}
		// Clean up project (last since other entities may reference it)
		if testProjectID != "" {
			runCLI(t, "project", "delete", testProjectID)
		}
		// Clean up label if still around
		if testLabelID != "" {
			runCLI(t, "label", "delete", testLabelID)
		}
		// Clean up view if still around
		if testViewID != "" {
			runCLI(t, "view", "delete", testViewID)
		}
	})

	t.Run("Team", testTeam)
	t.Run("User", testUser)
	t.Run("Label", testLabel)
	t.Run("Project", testProject)
	t.Run("Initiative", testInitiative)
	t.Run("Milestone", testMilestone)
	t.Run("Document", testDocument)
	t.Run("Cycle", testCycle)
	t.Run("Issue", testIssue)
	t.Run("Comment", testComment)
	t.Run("Relation", testRelation)
	t.Run("Attachment", testAttachment)
	t.Run("View", testView)
	t.Run("Favorite", testFavorite)
	t.Run("StatusUpdate", testStatusUpdate)
	t.Run("Inbox", testInbox)
	t.Run("GraphQL", testGraphQL)
	t.Run("ErrorHandling", testErrorHandling)
}

// =============================================================================
// Team tests (read-only)
// =============================================================================

func testTeam(t *testing.T) {
	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "team", "list")
		assertNotEmpty(t, out)
		assertContains(t, out, teamKey)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "team", "list", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "team", "list", "--json")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("expected at least one team")
		}
		if _, ok := arr[0]["key"]; !ok {
			t.Error("expected 'key' field in team JSON")
		}
	})

	t.Run("List_Limit", func(t *testing.T) {
		out := runCLISuccess(t, "team", "list", "--json", "--limit", "1")
		arr := parseJSONArray(t, out)
		if len(arr) != 1 {
			t.Errorf("expected 1 team with --limit 1, got %d", len(arr))
		}
	})

	t.Run("Get", func(t *testing.T) {
		out := runCLISuccess(t, "team", "get", teamKey)
		assertContains(t, out, teamKey)
	})

	t.Run("Get_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "team", "get", teamKey, "--json")
		key := extractField(t, out, "key")
		if key != teamKey {
			t.Errorf("expected team key %s, got %s", teamKey, key)
		}
	})

	t.Run("Members", func(t *testing.T) {
		out := runCLISuccess(t, "team", "members", teamKey)
		assertNotEmpty(t, out)
	})

	t.Run("Members_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "team", "members", teamKey, "--json")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("expected at least one team member")
		}
	})

	t.Run("States", func(t *testing.T) {
		out := runCLISuccess(t, "team", "states", teamKey)
		assertNotEmpty(t, out)
	})

	t.Run("States_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "team", "states", teamKey, "--json")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("expected at least one team state")
		}
		if _, ok := arr[0]["type"]; !ok {
			t.Error("expected 'type' field in state JSON")
		}
	})

	t.Run("States_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "team", "states", teamKey, "-p")
		assertNotEmpty(t, out)
	})
}

// =============================================================================
// User tests (read-only)
// =============================================================================

func testUser(t *testing.T) {
	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "user", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "user", "list", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "user", "list", "--json")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Fatal("expected at least one user")
		}
		if _, ok := arr[0]["email"]; !ok {
			t.Error("expected 'email' field in user JSON")
		}
	})

	t.Run("List_Active", func(t *testing.T) {
		out := runCLISuccess(t, "user", "list", "--json", "--active")
		arr := parseJSONArray(t, out)
		for _, u := range arr {
			if active, ok := u["active"]; ok && active == false {
				t.Error("found inactive user when --active filter was used")
			}
		}
	})

	t.Run("List_Limit", func(t *testing.T) {
		out := runCLISuccess(t, "user", "list", "--json", "--limit", "1")
		arr := parseJSONArray(t, out)
		if len(arr) != 1 {
			t.Errorf("expected 1 user with --limit 1, got %d", len(arr))
		}
	})

	t.Run("Me_Table", func(t *testing.T) {
		out := runCLISuccess(t, "user", "me")
		assertNotEmpty(t, out)
	})

	t.Run("Me_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "user", "me", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("Me_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "user", "me", "--json")
		id := extractField(t, out, "id")
		if id == "" {
			t.Error("expected non-empty user ID from 'me'")
		}
	})

	t.Run("Get_JSON", func(t *testing.T) {
		// Get our own user ID first
		meOut := runCLISuccess(t, "user", "me", "--json")
		userID := extractField(t, meOut, "id")

		out := runCLISuccess(t, "user", "get", userID, "--json")
		gotID := extractField(t, out, "id")
		if gotID != userID {
			t.Errorf("expected user ID %s, got %s", userID, gotID)
		}
	})
}

// =============================================================================
// Label tests (CRUD)
// =============================================================================

func testLabel(t *testing.T) {
	labelName := testPrefix + "-label"
	var labelID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "label", "create",
			"--name", labelName,
			"--color", "#e11d48",
			"--description", "Test label for CRUD tests",
			"--team-id", teamUUID,
			"--json",
		)
		labelID = extractID(t, out)
		if labelID == "" {
			t.Fatal("expected non-empty label ID")
		}
		testLabelID = labelID
		testLabelName = labelName
		name := extractField(t, out, "name")
		if name != labelName {
			t.Errorf("expected label name %s, got %s", labelName, name)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		if labelID == "" {
			t.Skip("no label created")
		}
		out := runCLISuccess(t, "label", "list", "--json", "--team", teamKey)
		arr := parseJSONArray(t, out)
		found := false
		for _, l := range arr {
			if fmt.Sprintf("%v", l["id"]) == labelID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created label %s not found in list", labelID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "label", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "label", "list", "-p")
		assertContains(t, out, "# Labels")
	})

	t.Run("Update", func(t *testing.T) {
		if labelID == "" {
			t.Skip("no label created")
		}
		updatedName := labelName + "-updated"
		out := runCLISuccess(t, "label", "update", labelID,
			"--name", updatedName,
			"--color", "#2563eb",
			"--description", "Updated description",
			"--json",
		)
		name := extractField(t, out, "name")
		if name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, name)
		}
		color := extractField(t, out, "color")
		if color != "#2563eb" {
			t.Errorf("expected color #2563eb, got %s", color)
		}
		// Keep the original name for issue tests
		testLabelName = updatedName
	})

	t.Run("Delete", func(t *testing.T) {
		if labelID == "" {
			t.Skip("no label created")
		}
		runCLISuccess(t, "label", "delete", labelID)
		// Clear shared state
		testLabelID = ""
		testLabelName = ""
	})
}

// =============================================================================
// Project tests (CRUD)
// =============================================================================

func testProject(t *testing.T) {
	projectName := testPrefix + "-project"
	var projectID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "project", "create",
			"--name", projectName,
			"--description", "Test project for CRUD tests",
			"--state", "planned",
			"--team-ids", teamUUID,
			"--json",
		)
		projectID = extractID(t, out)
		if projectID == "" {
			t.Fatal("expected non-empty project ID")
		}
		testProjectID = projectID
		name := extractField(t, out, "name")
		if name != projectName {
			t.Errorf("expected project name %s, got %s", projectName, name)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		if projectID == "" {
			t.Skip("no project created")
		}
		out := runCLISuccess(t, "project", "list", "--json", "--newer-than", "all_time")
		arr := parseJSONArray(t, out)
		found := false
		for _, p := range arr {
			if fmt.Sprintf("%v", p["id"]) == projectID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created project %s not found in list", projectID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "project", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "project", "list", "-p")
		assertContains(t, out, "# Projects")
	})

	t.Run("List_StateFilter", func(t *testing.T) {
		out := runCLISuccess(t, "project", "list", "--json", "--state", "planned")
		arr := parseJSONArray(t, out)
		for _, p := range arr {
			if fmt.Sprintf("%v", p["state"]) != "planned" {
				t.Errorf("expected state 'planned', got %v", p["state"])
			}
		}
	})

	t.Run("List_Limit", func(t *testing.T) {
		out := runCLISuccess(t, "project", "list", "--json", "--limit", "1", "--newer-than", "all_time")
		arr := parseJSONArray(t, out)
		if len(arr) > 1 {
			t.Errorf("expected at most 1 project with --limit 1, got %d", len(arr))
		}
	})

	t.Run("Get_Table", func(t *testing.T) {
		if projectID == "" {
			t.Skip("no project created")
		}
		out := runCLISuccess(t, "project", "get", projectID)
		assertContains(t, out, "Project:")
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if projectID == "" {
			t.Skip("no project created")
		}
		out := runCLISuccess(t, "project", "get", projectID, "--json")
		id := extractID(t, out)
		if id != projectID {
			t.Errorf("expected project ID %s, got %s", projectID, id)
		}
	})

	t.Run("Get_Plaintext", func(t *testing.T) {
		if projectID == "" {
			t.Skip("no project created")
		}
		out := runCLISuccess(t, "project", "get", projectID, "-p")
		assertContains(t, out, "# ")
	})

	t.Run("Update", func(t *testing.T) {
		if projectID == "" {
			t.Skip("no project created")
		}
		updatedName := projectName + "-updated"
		out := runCLISuccess(t, "project", "update", projectID,
			"--name", updatedName,
			"--description", "Updated project description",
			"--json",
		)
		name := extractField(t, out, "name")
		if name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, name)
		}
	})

	t.Run("Issues_Empty", func(t *testing.T) {
		if projectID == "" {
			t.Skip("no project created")
		}
		// New project should have no issues
		out := runCLISuccess(t, "project", "issues", projectID, "--json")
		arr := parseJSONArray(t, out)
		if len(arr) != 0 {
			t.Errorf("expected 0 issues for new project, got %d", len(arr))
		}
	})

	// Note: project cleanup is handled by TestCRUD's t.Cleanup since
	// later subtests (Milestone, Issue, etc.) depend on testProjectID.
}

// =============================================================================
// Initiative tests (CRUD)
// =============================================================================

func testInitiative(t *testing.T) {
	initName := testPrefix + "-initiative"
	var initID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "initiative", "create",
			"--name", initName,
			"--description", "Test initiative for CRUD tests",
			"--status", "Planned",
			"--json",
		)
		initID = extractID(t, out)
		if initID == "" {
			t.Fatal("expected non-empty initiative ID")
		}
		name := extractField(t, out, "name")
		if name != initName {
			t.Errorf("expected initiative name %s, got %s", initName, name)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "initiative", "list", "--json", "--include-completed")
		arr := parseJSONArray(t, out)
		found := false
		for _, i := range arr {
			if fmt.Sprintf("%v", i["id"]) == initID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created initiative %s not found in list", initID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "initiative", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "initiative", "list", "-p")
		assertContains(t, out, "# Initiatives")
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if initID == "" {
			t.Skip("no initiative created")
		}
		out := runCLISuccess(t, "initiative", "get", initID, "--json")
		id := extractID(t, out)
		if id != initID {
			t.Errorf("expected initiative ID %s, got %s", initID, id)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if initID == "" {
			t.Skip("no initiative created")
		}
		updatedName := initName + "-updated"
		out := runCLISuccess(t, "initiative", "update", initID,
			"--name", updatedName,
			"--description", "Updated initiative description",
			"--json",
		)
		name := extractField(t, out, "name")
		if name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, name)
		}
	})

	t.Run("AddProject", func(t *testing.T) {
		if initID == "" || testProjectID == "" {
			t.Skip("no initiative or project created")
		}
		runCLISuccess(t, "initiative", "add-project", initID, testProjectID)
	})

	t.Run("Projects", func(t *testing.T) {
		if initID == "" {
			t.Skip("no initiative created")
		}
		out := runCLISuccess(t, "initiative", "projects", initID, "--json")
		// Should have at least the project we added
		arr := parseJSONArray(t, out)
		if testProjectID != "" && len(arr) == 0 {
			t.Error("expected at least one project after add-project")
		}
	})

	t.Run("RemoveProject", func(t *testing.T) {
		if initID == "" || testProjectID == "" {
			t.Skip("no initiative or project created")
		}
		runCLISuccess(t, "initiative", "remove-project", initID, testProjectID)
	})

	t.Run("Delete", func(t *testing.T) {
		if initID == "" {
			t.Skip("no initiative created")
		}
		runCLISuccess(t, "initiative", "delete", initID)
	})
}

// =============================================================================
// Milestone tests (CRUD - nested under project)
// =============================================================================

func testMilestone(t *testing.T) {
	if testProjectID == "" {
		t.Skip("no project available for milestone tests")
	}

	msName := testPrefix + "-milestone"
	var msID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "project", "milestone", "create", testProjectID,
			"--name", msName,
			"--description", "Test milestone",
			"--target-date", "2026-12-31",
			"--json",
		)
		msID = extractID(t, out)
		if msID == "" {
			t.Fatal("expected non-empty milestone ID")
		}
		name := extractField(t, out, "name")
		if name != msName {
			t.Errorf("expected milestone name %s, got %s", msName, name)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "project", "milestone", "list", testProjectID, "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, ms := range arr {
			if fmt.Sprintf("%v", ms["id"]) == msID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created milestone %s not found in list", msID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "project", "milestone", "list", testProjectID)
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "project", "milestone", "list", testProjectID, "-p")
		assertContains(t, out, "# Milestones")
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if msID == "" {
			t.Skip("no milestone created")
		}
		out := runCLISuccess(t, "project", "milestone", "get", msID, "--json")
		id := extractID(t, out)
		if id != msID {
			t.Errorf("expected milestone ID %s, got %s", msID, id)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if msID == "" {
			t.Skip("no milestone created")
		}
		updatedName := msName + "-updated"
		out := runCLISuccess(t, "project", "milestone", "update", msID,
			"--name", updatedName,
			"--description", "Updated milestone",
			"--json",
		)
		name := extractField(t, out, "name")
		if name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if msID == "" {
			t.Skip("no milestone created")
		}
		runCLISuccess(t, "project", "milestone", "delete", msID)
	})
}

// =============================================================================
// Document tests (CRUD)
// =============================================================================

func testDocument(t *testing.T) {
	docTitle := testPrefix + "-document"
	var docID string

	t.Run("Create", func(t *testing.T) {
		args := []string{"document", "create",
			"--title", docTitle,
			"--content", "# Test Document\n\nThis is a test document for CRUD tests.",
			"--json",
		}
		if testProjectID != "" {
			args = append(args, "--project", testProjectID)
		} else {
			args = append(args, "--team", teamKey)
		}
		out := runCLISuccess(t, args...)
		docID = extractID(t, out)
		if docID == "" {
			t.Fatal("expected non-empty document ID")
		}
		title := extractField(t, out, "title")
		if title != docTitle {
			t.Errorf("expected document title %s, got %s", docTitle, title)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "document", "list", "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, d := range arr {
			if fmt.Sprintf("%v", d["id"]) == docID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created document %s not found in list", docID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "document", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "document", "list", "-p")
		assertContains(t, out, "# Documents")
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document created")
		}
		out := runCLISuccess(t, "document", "get", docID, "--json")
		id := extractID(t, out)
		if id != docID {
			t.Errorf("expected document ID %s, got %s", docID, id)
		}
	})

	t.Run("Get_Plaintext", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document created")
		}
		out := runCLISuccess(t, "document", "get", docID, "-p")
		assertContains(t, out, "# ")
	})

	t.Run("Search", func(t *testing.T) {
		out := runCLISuccess(t, "document", "search", testPrefix, "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, d := range arr {
			if fmt.Sprintf("%v", d["id"]) == docID {
				found = true
				break
			}
		}
		if !found {
			t.Logf("warning: document search for %q did not find created document (search indexing delay)", testPrefix)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document created")
		}
		updatedTitle := docTitle + "-updated"
		out := runCLISuccess(t, "document", "update", docID,
			"--title", updatedTitle,
			"--content", "# Updated Document\n\nUpdated content.",
			"--json",
		)
		title := extractField(t, out, "title")
		if title != updatedTitle {
			t.Errorf("expected updated title %s, got %s", updatedTitle, title)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if docID == "" {
			t.Skip("no document created")
		}
		runCLISuccess(t, "document", "delete", docID)
	})
}

// =============================================================================
// Cycle tests (CRUD)
// =============================================================================

func testCycle(t *testing.T) {
	cycleName := testPrefix + "-cycle"
	var cycleID string

	// Cycles need future dates to avoid conflicts
	starts := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	ends := time.Now().AddDate(1, 0, 14).Format("2006-01-02")

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "cycle", "create",
			"--team-id", teamUUID,
			"--name", cycleName,
			"--description", "Test cycle for CRUD tests",
			"--starts", starts,
			"--ends", ends,
			"--json",
		)
		cycleID = extractID(t, out)
		if cycleID == "" {
			t.Fatal("expected non-empty cycle ID")
		}
		testCycleID = cycleID
		name := extractField(t, out, "name")
		if name != cycleName {
			t.Errorf("expected cycle name %s, got %s", cycleName, name)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "cycle", "list", "--json", "--team", teamKey)
		arr := parseJSONArray(t, out)
		found := false
		for _, c := range arr {
			if fmt.Sprintf("%v", c["id"]) == cycleID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created cycle %s not found in list", cycleID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "cycle", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "cycle", "list", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if cycleID == "" {
			t.Skip("no cycle created")
		}
		out := runCLISuccess(t, "cycle", "get", cycleID, "--json")
		id := extractID(t, out)
		if id != cycleID {
			t.Errorf("expected cycle ID %s, got %s", cycleID, id)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if cycleID == "" {
			t.Skip("no cycle created")
		}
		updatedName := cycleName + "-updated"
		out := runCLISuccess(t, "cycle", "update", cycleID,
			"--name", updatedName,
			"--description", "Updated cycle description",
			"--json",
		)
		name := extractField(t, out, "name")
		if name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, name)
		}
	})

	// Note: cycle cleanup is handled by TestCRUD's t.Cleanup since
	// later subtests (Issue) may depend on testCycleID.
}

// =============================================================================
// Issue tests (CRUD)
// =============================================================================

func testIssue(t *testing.T) {
	issueTitle := testPrefix + "-issue"
	var issueIdentifier string
	var issueUUID string

	// Create a second issue for relation tests
	issue2Title := testPrefix + "-issue-2"
	var issue2Identifier string

	t.Run("Create", func(t *testing.T) {
		args := []string{"issue", "create",
			"--title", issueTitle,
			"--description", "Test issue for CRUD tests",
			"--team", teamKey,
			"--priority", "2",
			"--json",
		}
		if testProjectID != "" {
			args = append(args, "--project", testProjectID)
		}
		out := runCLISuccess(t, args...)
		issueUUID = extractID(t, out)
		issueIdentifier = extractField(t, out, "identifier")
		if issueUUID == "" || issueIdentifier == "" {
			t.Fatalf("expected non-empty issue ID and identifier, got id=%s identifier=%s", issueUUID, issueIdentifier)
		}
		testIssueID = issueIdentifier
		testIssueUUID = issueUUID
	})

	t.Run("Create_Second", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "create",
			"--title", issue2Title,
			"--description", "Second test issue for relation tests",
			"--team", teamKey,
			"--json",
		)
		issue2Identifier = extractField(t, out, "identifier")
		if issue2Identifier == "" {
			t.Fatal("expected non-empty second issue identifier")
		}
	})

	t.Run("Get_Table", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "get", issueIdentifier)
		assertNotEmpty(t, out)
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "get", issueIdentifier, "--json")
		identifier := extractField(t, out, "identifier")
		if identifier != issueIdentifier {
			t.Errorf("expected identifier %s, got %s", issueIdentifier, identifier)
		}
	})

	t.Run("Get_Plaintext", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "get", issueIdentifier, "-p")
		assertContains(t, out, "# "+issueIdentifier)
	})

	t.Run("Search_JSON", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "search", issueIdentifier, "--json")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Error("search returned no results for created issue")
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "list", "--json", "--team", teamKey, "--limit", "5")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Error("issue list returned no results")
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "list", "--team", teamKey, "--limit", "5")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "list", "-p", "--team", teamKey, "--limit", "5")
		assertContains(t, out, "# Issues")
	})

	t.Run("List_Filters", func(t *testing.T) {
		// Test assignee filter
		runCLISuccess(t, "issue", "list", "--json", "--assignee", "me", "--limit", "5")
		// Test priority filter
		runCLISuccess(t, "issue", "list", "--json", "--priority", "2", "--team", teamKey, "--limit", "5")
		// Test newer-than filter
		runCLISuccess(t, "issue", "list", "--json", "--newer-than", "1_week_ago", "--team", teamKey, "--limit", "5")
		// Test sort
		runCLISuccess(t, "issue", "list", "--json", "--sort", "updated", "--team", teamKey, "--limit", "5")
	})

	t.Run("Update", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		updatedTitle := issueTitle + "-updated"
		out := runCLISuccess(t, "issue", "update", issueIdentifier,
			"--title", updatedTitle,
			"--description", "Updated description",
			"--priority", "1",
			"--json",
		)
		title := extractField(t, out, "title")
		if title != updatedTitle {
			t.Errorf("expected updated title %s, got %s", updatedTitle, title)
		}
	})

	t.Run("Start", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "start", issueIdentifier, "--json")
		// Should show the issue with started state
		assertNotEmpty(t, out)
	})

	t.Run("Done", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "done", issueIdentifier, "--json")
		assertNotEmpty(t, out)
	})

	t.Run("Assign", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "assign", issueIdentifier, "--json")
		assertNotEmpty(t, out)
	})

	t.Run("Activity_JSON", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "activity", issueIdentifier, "--json")
		// Activity should have entries since we just did create/start/done
		assertNotEmpty(t, out)
	})

	t.Run("Activity_Plaintext", func(t *testing.T) {
		if issueIdentifier == "" {
			t.Skip("no issue created")
		}
		out := runCLISuccess(t, "issue", "activity", issueIdentifier, "-p")
		assertContains(t, out, "# Activity:")
	})

	t.Run("Triage", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "triage", teamKey, "--json")
		// May be empty, that's fine
		assertNotEmpty(t, out)
	})

	// Clean up the second issue (local to this subtest).
	// The first issue cleanup is handled by TestCRUD's t.Cleanup.
	t.Cleanup(func() {
		if issue2Identifier != "" {
			runCLI(t, "issue", "archive", issue2Identifier)
		}
	})
}

// =============================================================================
// Comment tests (CRUD - nested under issue)
// =============================================================================

func testComment(t *testing.T) {
	if testIssueID == "" {
		t.Skip("no issue available for comment tests")
	}

	var commentID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "comment", "create", testIssueID,
			"--body", "Test comment from CRUD tests",
			"--json",
		)
		commentID = extractID(t, out)
		if commentID == "" {
			t.Fatal("expected non-empty comment ID")
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "comment", "list", testIssueID, "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, c := range arr {
			if fmt.Sprintf("%v", c["id"]) == commentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created comment %s not found in list", commentID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "comment", "list", testIssueID)
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "comment", "list", testIssueID, "-p")
		assertNotEmpty(t, out)
	})

	t.Run("Update_Body", func(t *testing.T) {
		if commentID == "" {
			t.Skip("no comment created")
		}
		out := runCLISuccess(t, "issue", "comment", "update", commentID,
			"--body", "Updated comment body",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Update_Resolve", func(t *testing.T) {
		if commentID == "" {
			t.Skip("no comment created")
		}
		out := runCLISuccess(t, "issue", "comment", "update", commentID,
			"--resolve",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Update_Unresolve", func(t *testing.T) {
		if commentID == "" {
			t.Skip("no comment created")
		}
		out := runCLISuccess(t, "issue", "comment", "update", commentID,
			"--unresolve",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Delete", func(t *testing.T) {
		if commentID == "" {
			t.Skip("no comment created")
		}
		runCLISuccess(t, "issue", "comment", "delete", commentID)
	})
}

// =============================================================================
// Relation tests (CRUD - nested under issue)
// =============================================================================

func testRelation(t *testing.T) {
	if testIssueID == "" {
		t.Skip("no issue available for relation tests")
	}

	// We need a second issue. Create one if Issue tests didn't provide one.
	relIssueTitle := testPrefix + "-rel-target"
	var targetIdentifier string

	t.Run("Setup_TargetIssue", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "create",
			"--title", relIssueTitle,
			"--description", "Target issue for relation tests",
			"--team", teamKey,
			"--json",
		)
		targetIdentifier = extractField(t, out, "identifier")
		if targetIdentifier == "" {
			t.Fatal("expected non-empty target issue identifier")
		}
	})

	t.Run("Add_Blocks", func(t *testing.T) {
		if targetIdentifier == "" {
			t.Skip("no target issue")
		}
		runCLISuccess(t, "issue", "relation", "add", testIssueID,
			"--type", "blocks",
			"--target", targetIdentifier,
		)
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "relation", "list", testIssueID, "--json")
		assertNotEmpty(t, out)
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "relation", "list", testIssueID)
		assertNotEmpty(t, out)
	})

	t.Run("Remove", func(t *testing.T) {
		if targetIdentifier == "" {
			t.Skip("no target issue")
		}
		runCLISuccess(t, "issue", "relation", "remove", testIssueID,
			"--type", "blocks",
			"--target", targetIdentifier,
		)
	})

	t.Run("Add_Related", func(t *testing.T) {
		if targetIdentifier == "" {
			t.Skip("no target issue")
		}
		runCLISuccess(t, "issue", "relation", "add", testIssueID,
			"--type", "related",
			"--target", targetIdentifier,
		)
	})

	t.Run("Remove_Related", func(t *testing.T) {
		if targetIdentifier == "" {
			t.Skip("no target issue")
		}
		runCLISuccess(t, "issue", "relation", "remove", testIssueID,
			"--type", "related",
			"--target", targetIdentifier,
		)
	})

	// Cleanup target issue
	t.Cleanup(func() {
		if targetIdentifier != "" {
			runCLI(t, "issue", "archive", targetIdentifier)
		}
	})
}

// =============================================================================
// Attachment tests (CRUD - nested under issue)
// =============================================================================

func testAttachment(t *testing.T) {
	if testIssueID == "" {
		t.Skip("no issue available for attachment tests")
	}

	var attachmentID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "attachment", "create", testIssueID,
			"--url", "https://example.com/test-attachment",
			"--title", testPrefix+"-attachment",
			"--subtitle", "Test attachment subtitle",
			"--json",
		)
		attachmentID = extractID(t, out)
		if attachmentID == "" {
			t.Fatal("expected non-empty attachment ID")
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "attachment", "list", testIssueID, "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, a := range arr {
			if fmt.Sprintf("%v", a["id"]) == attachmentID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created attachment %s not found in list", attachmentID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "attachment", "list", testIssueID)
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "attachment", "list", testIssueID, "-p")
		assertNotEmpty(t, out)
	})

	t.Run("Link", func(t *testing.T) {
		out := runCLISuccess(t, "issue", "attachment", "link", testIssueID,
			"--url", "https://github.com/example/repo/pull/1",
			"--title", testPrefix+"-link",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Update", func(t *testing.T) {
		if attachmentID == "" {
			t.Skip("no attachment created")
		}
		out := runCLISuccess(t, "issue", "attachment", "update", attachmentID,
			"--title", testPrefix+"-attachment-updated",
			"--subtitle", "Updated subtitle",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Delete", func(t *testing.T) {
		if attachmentID == "" {
			t.Skip("no attachment created")
		}
		runCLISuccess(t, "issue", "attachment", "delete", attachmentID)
	})
}

// =============================================================================
// View tests (CRUD)
// =============================================================================

func testView(t *testing.T) {
	viewName := testPrefix + "-view"
	var viewID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "view", "create",
			"--name", viewName,
			"--description", "Test view for CRUD tests",
			"--model", "issue",
			"--team", teamKey,
			"--json",
		)
		viewID = extractID(t, out)
		if viewID == "" {
			t.Fatal("expected non-empty view ID")
		}
		testViewID = viewID
		name := extractField(t, out, "name")
		if name != viewName {
			t.Errorf("expected view name %s, got %s", viewName, name)
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "view", "list", "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, v := range arr {
			if fmt.Sprintf("%v", v["id"]) == viewID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created view %s not found in list", viewID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "view", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "view", "list", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if viewID == "" {
			t.Skip("no view created")
		}
		out := runCLISuccess(t, "view", "get", viewID, "--json")
		id := extractID(t, out)
		if id != viewID {
			t.Errorf("expected view ID %s, got %s", viewID, id)
		}
	})

	t.Run("Get_Plaintext", func(t *testing.T) {
		if viewID == "" {
			t.Skip("no view created")
		}
		out := runCLISuccess(t, "view", "get", viewID, "-p")
		assertContains(t, out, "# ")
	})

	t.Run("Run_JSON", func(t *testing.T) {
		if viewID == "" {
			t.Skip("no view created")
		}
		out := runCLISuccess(t, "view", "run", viewID, "--json", "--limit", "5")
		// View run may return empty results, that's OK
		assertNotEmpty(t, out)
	})

	t.Run("Run_Table", func(t *testing.T) {
		if viewID == "" {
			t.Skip("no view created")
		}
		// May return empty results - just ensure no error
		runCLISuccess(t, "view", "run", viewID, "--limit", "5")
	})

	t.Run("Update", func(t *testing.T) {
		if viewID == "" {
			t.Skip("no view created")
		}
		updatedName := viewName + "-updated"
		out := runCLISuccess(t, "view", "update", viewID,
			"--name", updatedName,
			"--description", "Updated view description",
			"--json",
		)
		name := extractField(t, out, "name")
		if name != updatedName {
			t.Errorf("expected updated name %s, got %s", updatedName, name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if viewID == "" {
			t.Skip("no view created")
		}
		runCLISuccess(t, "view", "delete", viewID)
		testViewID = ""
	})
}

// =============================================================================
// Favorite tests (CRUD)
// =============================================================================

func testFavorite(t *testing.T) {
	var favIssueID string
	var favProjectID string

	t.Run("Add_Issue", func(t *testing.T) {
		if testIssueID == "" {
			t.Skip("no issue available")
		}
		out := runCLISuccess(t, "favorite", "add",
			"--issue", testIssueID,
			"--json",
		)
		favIssueID = extractID(t, out)
		if favIssueID == "" {
			t.Fatal("expected non-empty favorite ID")
		}
	})

	t.Run("Add_Project", func(t *testing.T) {
		if testProjectID == "" {
			t.Skip("no project available")
		}
		out := runCLISuccess(t, "favorite", "add",
			"--project", testProjectID,
			"--json",
		)
		favProjectID = extractID(t, out)
		if favProjectID == "" {
			t.Fatal("expected non-empty favorite ID")
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		if favIssueID == "" && favProjectID == "" {
			t.Skip("no favorites added")
		}
		out := runCLISuccess(t, "favorite", "list", "--json")
		arr := parseJSONArray(t, out)
		if len(arr) == 0 {
			t.Error("expected at least one favorite")
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "favorite", "list")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "favorite", "list", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("Update_SortOrder", func(t *testing.T) {
		if favIssueID == "" {
			t.Skip("no issue favorite")
		}
		out := runCLISuccess(t, "favorite", "update", favIssueID,
			"--sort-order", "5.0",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Remove_Issue", func(t *testing.T) {
		if favIssueID == "" {
			t.Skip("no issue favorite")
		}
		runCLISuccess(t, "favorite", "remove", favIssueID)
	})

	t.Run("Remove_Project", func(t *testing.T) {
		if favProjectID == "" {
			t.Skip("no project favorite")
		}
		runCLISuccess(t, "favorite", "remove", favProjectID)
	})
}

// =============================================================================
// StatusUpdate tests (CRUD - nested under project as "status")
// =============================================================================

func testStatusUpdate(t *testing.T) {
	if testProjectID == "" {
		t.Skip("no project available for status update tests")
	}

	var updateID string

	t.Run("Create", func(t *testing.T) {
		out := runCLISuccess(t, "project", "status", "create", testProjectID,
			"--body", "Test status update from CRUD tests",
			"--health", "onTrack",
			"--json",
		)
		updateID = extractID(t, out)
		if updateID == "" {
			t.Fatal("expected non-empty status update ID")
		}
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "project", "status", "list", testProjectID, "--json")
		arr := parseJSONArray(t, out)
		found := false
		for _, u := range arr {
			if fmt.Sprintf("%v", u["id"]) == updateID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("created status update %s not found in list", updateID)
		}
	})

	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "project", "status", "list", testProjectID)
		assertNotEmpty(t, out)
	})

	t.Run("Get_JSON", func(t *testing.T) {
		if updateID == "" {
			t.Skip("no status update created")
		}
		out := runCLISuccess(t, "project", "status", "get", updateID, "--json")
		id := extractID(t, out)
		if id != updateID {
			t.Errorf("expected status update ID %s, got %s", updateID, id)
		}
	})

	t.Run("Update", func(t *testing.T) {
		if updateID == "" {
			t.Skip("no status update created")
		}
		out := runCLISuccess(t, "project", "status", "update", updateID,
			"--body", "Updated status update body",
			"--health", "atRisk",
			"--json",
		)
		assertNotEmpty(t, out)
	})

	t.Run("Delete", func(t *testing.T) {
		if updateID == "" {
			t.Skip("no status update created")
		}
		runCLISuccess(t, "project", "status", "delete", updateID)
	})
}

// =============================================================================
// Inbox tests (read-only)
// =============================================================================

func testInbox(t *testing.T) {
	t.Run("List_Table", func(t *testing.T) {
		out := runCLISuccess(t, "inbox")
		assertNotEmpty(t, out)
	})

	t.Run("List_Plaintext", func(t *testing.T) {
		out := runCLISuccess(t, "inbox", "-p")
		assertNotEmpty(t, out)
	})

	t.Run("List_JSON", func(t *testing.T) {
		out := runCLISuccess(t, "inbox", "--json")
		assertNotEmpty(t, out)
	})

	t.Run("List_Limit", func(t *testing.T) {
		runCLISuccess(t, "inbox", "--json", "--limit", "5")
	})

	t.Run("List_Unread", func(t *testing.T) {
		// May return empty but should not error
		runCLISuccess(t, "inbox", "--json", "--unread")
	})

	t.Run("List_All", func(t *testing.T) {
		runCLISuccess(t, "inbox", "--json", "--all")
	})
}

// =============================================================================
// GraphQL tests
// =============================================================================

func testGraphQL(t *testing.T) {
	t.Run("Viewer_Query", func(t *testing.T) {
		out := runCLISuccess(t, "graphql", `query { viewer { id name email } }`)
		assertContains(t, out, "viewer")
		assertContains(t, out, "email")
	})

	t.Run("Teams_Query", func(t *testing.T) {
		out := runCLISuccess(t, "graphql", `query { teams { nodes { id name key } } }`)
		assertContains(t, out, "teams")
		assertContains(t, out, teamKey)
	})

	t.Run("With_Variables", func(t *testing.T) {
		if testIssueUUID == "" {
			t.Skip("no issue UUID available")
		}
		vars := fmt.Sprintf(`{"id": "%s"}`, testIssueUUID)
		out := runCLISuccess(t, "graphql",
			`query($id: String!) { issue(id: $id) { id title identifier } }`,
			"--variables", vars,
		)
		assertContains(t, out, "issue")
		assertContains(t, out, testIssueID)
	})
}

// =============================================================================
// Error handling tests
// =============================================================================

func testErrorHandling(t *testing.T) {
	t.Run("Unknown_Command", func(t *testing.T) {
		out := runCLIFail(t, "nonexistent-command")
		assertContains(t, out, "unknown command")
	})

	t.Run("Issue_Create_Missing_Title", func(t *testing.T) {
		out := runCLIFail(t, "issue", "create", "--team", teamKey)
		assertContains(t, out, "required")
	})

	t.Run("Issue_Create_Missing_Team", func(t *testing.T) {
		out := runCLIFail(t, "issue", "create", "--title", "test")
		assertContains(t, out, "required")
	})

	t.Run("Issue_Get_Nonexistent", func(t *testing.T) {
		runCLIFail(t, "issue", "get", "NONEXIST-99999")
	})

	t.Run("Project_Get_Nonexistent", func(t *testing.T) {
		runCLIFail(t, "project", "get", "00000000-0000-0000-0000-000000000000")
	})

	t.Run("Label_Create_Missing_Name", func(t *testing.T) {
		out := runCLIFail(t, "label", "create")
		assertContains(t, out, "required")
	})

	t.Run("Label_Update_No_Fields", func(t *testing.T) {
		runCLIFail(t, "label", "update", "00000000-0000-0000-0000-000000000000")
	})

	t.Run("Cycle_Create_Missing_Required", func(t *testing.T) {
		out := runCLIFail(t, "cycle", "create")
		assertContains(t, out, "required")
	})

	t.Run("GraphQL_Missing_Query", func(t *testing.T) {
		out := runCLIFail(t, "graphql")
		assertContains(t, out, "accepts 1 arg")
	})

	t.Run("GraphQL_Invalid_Query", func(t *testing.T) {
		runCLIFail(t, "graphql", "{ invalid_field_that_does_not_exist }")
	})
}
