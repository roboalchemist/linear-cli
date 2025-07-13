package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dorkitude/linctl/pkg/api"
	"github.com/fatih/color"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

type AuthConfig struct {
	APIKey      string    `json:"api_key,omitempty"`
	AccessToken string    `json:"access_token,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// getConfigPath returns the path to the auth config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".linctl-auth.json"), nil
}

// saveAuth saves authentication credentials
func saveAuth(config AuthConfig) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// loadAuth loads authentication credentials
func loadAuth() (*AuthConfig, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not authenticated")
		}
		return nil, err
	}

	var config AuthConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Check if token is expired (for OAuth)
	if !config.ExpiresAt.IsZero() && time.Now().After(config.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &config, nil
}

// GetAuthHeader returns the authorization header value
func GetAuthHeader() (string, error) {
	config, err := loadAuth()
	if err != nil {
		return "", err
	}

	if config.APIKey != "" {
		return config.APIKey, nil
	}

	if config.AccessToken != "" {
		return "Bearer " + config.AccessToken, nil
	}

	return "", fmt.Errorf("no valid authentication found")
}

// Login handles the authentication flow
func Login(plaintext, jsonOut bool) error {
	if !plaintext && !jsonOut {
		fmt.Println("Choose authentication method:")
		fmt.Println("1. Personal API Key (recommended for CLI)")
		fmt.Println("2. OAuth 2.0 (opens browser)")
		fmt.Print("Enter choice (1-2): ")
	}

	var choice string
	if plaintext || jsonOut {
		// Default to API key for non-interactive mode
		choice = "1"
	} else {
		fmt.Scanln(&choice)
	}

	switch choice {
	case "1", "":
		return loginWithAPIKey(plaintext, jsonOut)
	case "2":
		return loginWithOAuth(plaintext, jsonOut)
	default:
		return fmt.Errorf("invalid choice")
	}
}

// loginWithAPIKey handles Personal API Key authentication
func loginWithAPIKey(plaintext, jsonOut bool) error {
	if !plaintext && !jsonOut {
		fmt.Println("\n" + color.New(color.FgYellow).Sprint("📝 Personal API Key Authentication"))
		fmt.Println("Get your API key from: https://linear.app/settings/api")
		fmt.Print("Enter your Personal API Key: ")
	}

	reader := bufio.NewReader(os.Stdin)
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	apiKey = strings.TrimSpace(apiKey)

	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Test the API key
	client := api.NewClient(apiKey)
	user, err := client.GetViewer(context.Background())
	if err != nil {
		return fmt.Errorf("invalid API key: %v", err)
	}

	// Save the API key
	config := AuthConfig{
		APIKey: apiKey,
	}
	err = saveAuth(config)
	if err != nil {
		return err
	}

	if !plaintext && !jsonOut {
		fmt.Printf("\n%s Authenticated as %s (%s)\n",
			color.New(color.FgGreen).Sprint("✅"),
			color.New(color.FgCyan).Sprint(user.Name),
			color.New(color.FgCyan).Sprint(user.Email))
	}

	return nil
}

// loginWithOAuth handles OAuth 2.0 authentication
func loginWithOAuth(plaintext, jsonOut bool) error {
	if plaintext {
		return fmt.Errorf("OAuth authentication not supported in plaintext mode")
	}

	// This would require registering an OAuth app with Linear
	// For now, we'll return an error suggesting API key usage
	return fmt.Errorf("OAuth authentication not yet implemented. Please use Personal API Key (option 1)")
}

// GetCurrentUser returns the current authenticated user
func GetCurrentUser() (*User, error) {
	authHeader, err := GetAuthHeader()
	if err != nil {
		return nil, err
	}

	client := api.NewClient(authHeader)
	apiUser, err := client.GetViewer(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert api.User to auth.User
	return &User{
		ID:        apiUser.ID,
		Name:      apiUser.Name,
		Email:     apiUser.Email,
		AvatarURL: apiUser.AvatarURL,
	}, nil
}

// Logout clears stored credentials
func Logout() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.Remove(configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// openBrowser opens the default browser to the given URL
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}