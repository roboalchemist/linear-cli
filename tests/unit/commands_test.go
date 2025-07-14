package unit

import (
	"bytes"
	"os"
	"testing"
	
	"github.com/dorkitude/linctl/cmd"
)

// TestCommandStructure verifies that all commands are properly configured
func TestCommandStructure(t *testing.T) {
	rootCmd := cmd.GetRootCmd()
	
	// Test that root command exists
	if rootCmd == nil {
		t.Fatal("Root command is nil")
	}
	
	// Expected commands
	expectedCommands := []string{
		"auth",
		"issue",
		"team",
		"user",
		"project",
		"comment",
		"docs",
	}
	
	// Check all expected commands exist
	for _, cmdName := range expectedCommands {
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			t.Errorf("Command '%s' not found: %v", cmdName, err)
			continue
		}
		
		// Verify command has required fields
		if cmd.Use == "" {
			t.Errorf("Command '%s' has empty Use field", cmdName)
		}
		if cmd.Short == "" {
			t.Errorf("Command '%s' has empty Short description", cmdName)
		}
	}
}

// TestIssueCommands verifies issue subcommands
func TestIssueCommands(t *testing.T) {
	rootCmd := cmd.GetRootCmd()
	
	issueCmd, _, err := rootCmd.Find([]string{"issue"})
	if err != nil {
		t.Fatalf("Issue command not found: %v", err)
	}
	
	// Expected subcommands
	expectedSubcommands := []struct {
		name    string
		aliases []string
	}{
		{"list", []string{"ls"}},
		{"get", nil},
		{"create", nil},
		{"update", nil},
	}
	
	for _, sub := range expectedSubcommands {
		cmd, _, err := issueCmd.Find([]string{sub.name})
		if err != nil {
			t.Errorf("Subcommand 'issue %s' not found: %v", sub.name, err)
			continue
		}
		
		// Check aliases
		if len(sub.aliases) > 0 {
			if !equalStringSlices(cmd.Aliases, sub.aliases) {
				t.Errorf("Command 'issue %s' aliases mismatch. Expected %v, got %v", 
					sub.name, sub.aliases, cmd.Aliases)
			}
		}
	}
}

// TestAuthCommand verifies auth command functionality
func TestAuthCommand(t *testing.T) {
	// Skip if running in CI without proper setup
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping auth command test in CI")
	}
	
	rootCmd := cmd.GetRootCmd()
	
	// Test auth command exists
	authCmd, _, err := rootCmd.Find([]string{"auth"})
	if err != nil {
		t.Fatalf("Auth command not found: %v", err)
	}
	
	// Verify auth command has expected fields
	if authCmd.Use != "auth" {
		t.Errorf("Expected auth command Use to be 'auth', got '%s'", authCmd.Use)
	}
	
	if authCmd.Short == "" {
		t.Error("Auth command missing short description")
	}
}

// TestHelpOutput verifies help text is generated properly
func TestHelpOutput(t *testing.T) {
	rootCmd := cmd.GetRootCmd()
	
	// Capture help output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	
	// Run help command
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute help command: %v", err)
	}
	
	output := buf.String()
	
	// Verify help contains expected sections
	expectedStrings := []string{
		"linctl",
		"comprehensive CLI tool",
		"Available Commands:",
		"Flags:",
	}
	
	for _, expected := range expectedStrings {
		if !contains(output, expected) {
			t.Errorf("Help output missing expected string: %s", expected)
		}
	}
}

// TestGlobalFlags verifies global flags are properly configured
func TestGlobalFlags(t *testing.T) {
	rootCmd := cmd.GetRootCmd()
	
	// Expected global flags
	expectedFlags := []struct {
		name      string
		shorthand string
		defValue  string
	}{
		{"plaintext", "p", "false"},
		{"json", "j", "false"},
		{"config", "", ""},
	}
	
	for _, flag := range expectedFlags {
		f := rootCmd.PersistentFlags().Lookup(flag.name)
		if f == nil {
			t.Errorf("Global flag '%s' not found", flag.name)
			continue
		}
		
		if f.Shorthand != flag.shorthand {
			t.Errorf("Flag '%s' shorthand mismatch. Expected '%s', got '%s'",
				flag.name, flag.shorthand, f.Shorthand)
		}
		
		if f.DefValue != flag.defValue {
			t.Errorf("Flag '%s' default value mismatch. Expected '%s', got '%s'",
				flag.name, flag.defValue, f.DefValue)
		}
	}
}

// Helper functions
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}