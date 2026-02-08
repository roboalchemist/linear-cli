#!/bin/bash

# Smoke test for linear-cli - tests all READ commands
# This script runs through all the read-only commands to ensure basic functionality

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_pattern="$3"  # Optional pattern to check in output
    
    TESTS_RUN=$((TESTS_RUN + 1))
    echo -n "Testing: $test_name... "
    
    # Run the command and capture output and exit code
    set +e  # Temporarily disable exit on error
    output=$(eval "$command" 2>&1)
    exit_code=$?
    set -e
    
    if [ $exit_code -eq 0 ]; then
        # If an expected pattern is provided, check for it
        if [ -n "$expected_pattern" ]; then
            if echo "$output" | grep -q "$expected_pattern"; then
                echo -e "${GREEN}PASS${NC}"
                TESTS_PASSED=$((TESTS_PASSED + 1))
            else
                echo -e "${RED}FAIL${NC} - Expected pattern not found: $expected_pattern"
                echo "Output: $output" | head -5
                TESTS_FAILED=$((TESTS_FAILED + 1))
            fi
        else
            echo -e "${GREEN}PASS${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        fi
    else
        echo -e "${RED}FAIL${NC} - Exit code: $exit_code"
        echo "Error: $output" | head -5
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Function to extract first ID from list output
get_first_id() {
    local output="$1"
    # Extract first UUID-like pattern (for projects) or identifier like END-1234 (for issues)
    echo "$output" | grep -E -o '([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}|[A-Z]+-[0-9]+)' | head -1
}

echo "🚀 Starting linear-cli smoke tests..."
echo "================================"

# Check if authenticated
echo -e "\n${YELLOW}Checking authentication...${NC}"
run_test "auth status" "go run main.go auth status" "Authenticated"

# If not authenticated, skip tests
if [ $TESTS_FAILED -gt 0 ]; then
    echo -e "\n${RED}Not authenticated. Please run 'linear-cli auth' first.${NC}"
    exit 1
fi

# Test whoami
echo -e "\n${YELLOW}Testing user commands...${NC}"
run_test "whoami" "go run main.go whoami" "@"

# Test user list
run_test "user list" "go run main.go user list"
run_test "user list (plaintext)" "go run main.go user list -p"
run_test "user list (json)" "go run main.go user list -j" "\"email\""

# Note: user get requires user ID, not email - skipping for now
# Could be implemented by parsing JSON output from user list

# Test team commands
echo -e "\n${YELLOW}Testing team commands...${NC}"
run_test "team list" "go run main.go team list"
run_test "team list (plaintext)" "go run main.go team list -p"
run_test "team list (json)" "go run main.go team list -j" "\"key\""

# Get first team key for additional tests - look for pattern at start of line
team_key=$(go run main.go team list 2>/dev/null | awk 'NR>1 {print $1}' | head -1)
if [ -n "$team_key" ]; then
    run_test "team get" "go run main.go team get $team_key" "$team_key"
    run_test "team members" "go run main.go team members $team_key"
fi

# Test project commands
echo -e "\n${YELLOW}Testing project commands...${NC}"
run_test "project list" "go run main.go project list"
run_test "project list (plaintext)" "go run main.go project list -p" "# Projects"
run_test "project list (json)" "go run main.go project list -j" "\"id\""
# Note: project list doesn't support team filter in the API
# run_test "project list (with team filter)" "go run main.go project list --team $team_key" 
run_test "project list (state filter)" "go run main.go project list --state started"
run_test "project list (time filter)" "go run main.go project list --newer-than 1_month_ago"

# Get first project ID for project get test
project_output=$(go run main.go project list 2>/dev/null || true)
project_id=$(get_first_id "$project_output")
if [ -n "$project_id" ]; then
    run_test "project get" "go run main.go project get $project_id" "Project:"
    run_test "project get (plaintext)" "go run main.go project get $project_id -p" "# "
fi

# Test issue commands
echo -e "\n${YELLOW}Testing issue commands...${NC}"
run_test "issue list" "go run main.go issue list"
run_test "issue list (plaintext)" "go run main.go issue list -p" "# Issues"
run_test "issue list (json)" "go run main.go issue list -j" "\"identifier\""
run_test "issue list (assignee filter)" "go run main.go issue list --assignee me"
run_test "issue list (state filter)" "go run main.go issue list --state Todo"
run_test "issue list (team filter)" "go run main.go issue list --team $team_key"
run_test "issue list (priority filter)" "go run main.go issue list --priority 3"
run_test "issue list (time filter)" "go run main.go issue list --newer-than 2_weeks_ago"
run_test "issue list (sort by updated)" "go run main.go issue list --sort updated"

# Get first issue ID for additional tests
issue_output=$(go run main.go issue list --limit 5 2>/dev/null || true)
issue_id=$(echo "$issue_output" | grep -E -o '[A-Z]+-[0-9]+' | head -1)
if [ -n "$issue_id" ]; then
    run_test "issue search (default)" "go run main.go issue search $issue_id"
    run_test "issue search (json)" "go run main.go issue search $issue_id -j" "$issue_id"
    run_test "issue search (plaintext)" "go run main.go issue search $issue_id -p" "# Search Results"
    run_test "issue get" "go run main.go issue get $issue_id"
    run_test "issue get (plaintext)" "go run main.go issue get $issue_id -p" "# $issue_id"
    
    # Test comment list for this issue
    echo -e "\n${YELLOW}Testing comment commands...${NC}"
    run_test "comment list" "go run main.go comment list $issue_id"
    run_test "comment list (plaintext)" "go run main.go comment list $issue_id -p"
fi

# Test document commands
echo -e "\n${YELLOW}Testing document commands...${NC}"
run_test "document list" "go run main.go document list"
run_test "document list (plaintext)" "go run main.go document list -p" "# Documents"
run_test "document list (json)" "go run main.go document list -j" "\"id\""
run_test "document search" "go run main.go document search test"
run_test "document search (json)" "go run main.go document search test -j"

# Get first document ID for get test
doc_output=$(go run main.go document list --json 2>/dev/null || true)
doc_id=$(echo "$doc_output" | grep -E -o '"id":\s*"[a-f0-9-]+"' | head -1 | grep -E -o '[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}')
if [ -n "$doc_id" ]; then
    run_test "document get" "go run main.go document get $doc_id"
    run_test "document get (plaintext)" "go run main.go document get $doc_id -p" "# "
    run_test "document get (json)" "go run main.go document get $doc_id -j" "\"id\""
fi

# Test help commands
echo -e "\n${YELLOW}Testing help commands...${NC}"
run_test "help" "go run main.go --help" "Usage:"
run_test "issue help" "go run main.go issue --help" "Available Commands:"
run_test "project help" "go run main.go project --help" "Available Commands:"
run_test "project add-team help" "go run main.go project add-team --help" "Add one or more teams"
run_test "project remove-team help" "go run main.go project remove-team --help" "Remove one or more teams"
run_test "project milestone help" "go run main.go project milestone --help" "Available Commands:"
run_test "project milestone list help" "go run main.go project milestone list --help" "milestones"
run_test "project milestone get help" "go run main.go project milestone get --help" "milestone"
run_test "project milestone create help" "go run main.go project milestone create --help" "milestone"
run_test "project milestone update help" "go run main.go project milestone update --help" "milestone"
run_test "project milestone delete help" "go run main.go project milestone delete --help" "milestone"
run_test "project status help" "go run main.go project status --help" "Available Commands:"
run_test "project status list help" "go run main.go project status list --help" "status"
run_test "project status get help" "go run main.go project status get --help" "status"
run_test "project status create help" "go run main.go project status create --help" "status"
run_test "project status update help" "go run main.go project status update --help" "status"
run_test "project status delete help" "go run main.go project status delete --help" "status"
run_test "issue relation help" "go run main.go issue relation --help" "Available Commands:"
run_test "issue relation list help" "go run main.go issue relation list --help" "relation"
run_test "issue relation add help" "go run main.go issue relation add --help" "relation"
run_test "issue relation remove help" "go run main.go issue relation remove --help" "relation"
run_test "issue relation update help" "go run main.go issue relation update --help" "relation"
run_test "team help" "go run main.go team --help" "Available Commands:"
run_test "user help" "go run main.go user --help" "Available Commands:"
run_test "document help" "go run main.go document --help" "Available Commands:"
run_test "view help" "go run main.go view --help" "Available Commands:"
run_test "view list help" "go run main.go view list --help" "custom views"
run_test "view get help" "go run main.go view get --help" "view"
run_test "view run help" "go run main.go view run --help" "Execute"
run_test "view create help" "go run main.go view create --help" "view"
run_test "view update help" "go run main.go view update --help" "view"
run_test "view delete help" "go run main.go view delete --help" "view"
run_test "initiative help" "go run main.go initiative --help" "Available Commands:"
run_test "initiative list help" "go run main.go initiative list --help" "initiatives"
run_test "initiative get help" "go run main.go initiative get --help" "initiative"
run_test "initiative create help" "go run main.go initiative create --help" "initiative"
run_test "initiative update help" "go run main.go initiative update --help" "initiative"
run_test "initiative delete help" "go run main.go initiative delete --help" "initiative"
run_test "issue attachment help" "go run main.go issue attachment --help" "Available Commands:"
run_test "issue attachment list help" "go run main.go issue attachment list --help" "attachment"
run_test "issue attachment create help" "go run main.go issue attachment create --help" "attachment"
run_test "issue attachment link help" "go run main.go issue attachment link --help" "link"
run_test "issue attachment update help" "go run main.go issue attachment update --help" "attachment"
run_test "issue attachment delete help" "go run main.go issue attachment delete --help" "attachment"
run_test "issue activity help" "go run main.go issue activity --help" "activity"

# Test custom view commands
echo -e "\n${YELLOW}Testing custom view commands...${NC}"
run_test "view list" "go run main.go view list"
run_test "view list (plaintext)" "go run main.go view list -p"
run_test "view list (json)" "go run main.go view list -j"

# Get first view ID for additional tests
view_output=$(go run main.go view list --json 2>/dev/null || true)
view_id=$(echo "$view_output" | grep -E -o '"id":\s*"[a-f0-9-]+"' | head -1 | grep -E -o '[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}')
if [ -n "$view_id" ]; then
    run_test "view get" "go run main.go view get $view_id"
    run_test "view get (plaintext)" "go run main.go view get $view_id -p" "# "
    run_test "view get (json)" "go run main.go view get $view_id -j" "\"id\""
    run_test "view run" "go run main.go view run $view_id"
    run_test "view run (plaintext)" "go run main.go view run $view_id -p"
    run_test "view run (json)" "go run main.go view run $view_id -j"
fi

# Test initiative commands
echo -e "\n${YELLOW}Testing initiative commands...${NC}"
run_test "initiative list" "go run main.go initiative list"
run_test "initiative list (plaintext)" "go run main.go initiative list -p"
run_test "initiative list (json)" "go run main.go initiative list -j"
run_test "initiative list (include-completed)" "go run main.go initiative list --include-completed"

# Test issue activity command
echo -e "\n${YELLOW}Testing issue activity command...${NC}"
if [ -n "$issue_id" ]; then
    run_test "issue activity" "go run main.go issue activity $issue_id"
    run_test "issue activity (plaintext)" "go run main.go issue activity $issue_id -p" "# Activity:"
    run_test "issue activity (json)" "go run main.go issue activity $issue_id -j" "\"issue\""
fi

# Test issue attachment commands
echo -e "\n${YELLOW}Testing issue attachment commands...${NC}"
if [ -n "$issue_id" ]; then
    run_test "issue attachment list" "go run main.go issue attachment list $issue_id"
    run_test "issue attachment list (plaintext)" "go run main.go issue attachment list $issue_id -p"
    run_test "issue attachment list (json)" "go run main.go issue attachment list $issue_id -j"
fi

# Test unknown command handling
echo -e "\n${YELLOW}Testing error handling...${NC}"
# This should fail but gracefully
set +e
output=$(go run main.go nonexistent-command 2>&1)
if echo "$output" | grep -q "unknown command"; then
    echo -e "Unknown command handling: ${GREEN}PASS${NC}"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "Unknown command handling: ${RED}FAIL${NC}"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_FAILED=$((TESTS_FAILED + 1))
fi
set -e

# Summary
echo -e "\n================================"
echo "Test Summary:"
echo "  Total tests: $TESTS_RUN"
echo -e "  Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "  Failed: ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}❌ Some tests failed!${NC}"
    exit 1
fi
