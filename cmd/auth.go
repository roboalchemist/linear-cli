package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dorkitude/linear-cli/pkg/api"
	"github.com/dorkitude/linear-cli/pkg/auth"
	"github.com/dorkitude/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Linear",
	Long: `Authenticate with Linear using Personal API Key.

Examples:
  linear-cli auth              # Interactive authentication
  linear-cli auth login        # Same as above
  linear-cli auth status       # Check authentication status
  linear-cli auth logout       # Clear stored credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior is to run login
		loginCmd.Run(cmd, args)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Linear",
	Long:  `Authenticate with Linear using Personal API Key.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		if !plaintext && !jsonOut {
			fmt.Println(color.New(color.FgCyan, color.Bold).Sprint("🔐 Linear Authentication"))
			fmt.Println()
		}

		err := auth.Login(plaintext, jsonOut)
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if !plaintext && !jsonOut {
			fmt.Println(color.New(color.FgGreen).Sprint("✅ Successfully authenticated with Linear!"))
		} else if jsonOut {
			output.JSON(map[string]interface{}{
				"status":  "success",
				"message": "Successfully authenticated with Linear",
			})
		} else {
			fmt.Println("Successfully authenticated with Linear")
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Check if you are currently authenticated with Linear.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		user, err := auth.GetCurrentUser()
		if err != nil {
			if !plaintext && !jsonOut {
				fmt.Println(color.New(color.FgRed).Sprint("❌ Not authenticated"))
			} else if jsonOut {
				output.JSON(map[string]interface{}{
					"authenticated": false,
					"error":         err.Error(),
				})
			} else {
				fmt.Println("Not authenticated")
			}
			os.Exit(1)
		}

		authSource := auth.GetAuthSource()
		sourceLabel := "config file"
		if authSource == "env:LINEAR_API_KEY" {
			sourceLabel = "LINEAR_API_KEY env var"
		} else if authSource == "env:LINCTL_API_KEY" {
			sourceLabel = "LINCTL_API_KEY env var"
		}

		if jsonOut {
			output.JSON(map[string]interface{}{
				"authenticated": true,
				"user":          user,
				"auth_source":   authSource,
			})
		} else if plaintext {
			fmt.Printf("Authenticated as: %s (%s)\n", user.Name, user.Email)
			fmt.Printf("Auth source: %s\n", sourceLabel)
		} else {
			fmt.Println(color.New(color.FgGreen).Sprint("✅ Authenticated"))
			fmt.Printf("User: %s\n", color.New(color.FgCyan).Sprint(user.Name))
			fmt.Printf("Email: %s\n", color.New(color.FgCyan).Sprint(user.Email))
			fmt.Printf("Source: %s\n", color.New(color.FgYellow).Sprint(sourceLabel))
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Linear",
	Long:  `Clear stored Linear credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		err := auth.Logout()
		if err != nil {
			output.Error(fmt.Sprintf("Logout failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{
				"status":  "success",
				"message": "Successfully logged out",
			})
		} else if plaintext {
			fmt.Println("Successfully logged out")
		} else {
			fmt.Println(color.New(color.FgGreen).Sprint("✅ Successfully logged out"))
		}
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user",
	Long:  `Display information about the currently authenticated user.`,
	Run: func(cmd *cobra.Command, args []string) {
		statusCmd.Run(cmd, args)
	},
}

var rateLimitCmd = &cobra.Command{
	Use:   "rate-limit",
	Short: "Show API rate limit status",
	Long:  `Check current Linear API rate limit usage (requests and complexity).`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		rl, err := client.GetRateLimit(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get rate limit: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(rl)
			return
		}

		reqResetIn := time.Until(rl.RequestReset).Round(time.Second)
		cplxResetIn := time.Until(rl.ComplexityReset).Round(time.Second)

		if plaintext {
			fmt.Printf("Requests: %d/%d remaining (resets in %s)\n",
				rl.RequestRemaining, rl.RequestLimit, reqResetIn)
			fmt.Printf("Complexity: %d/%d remaining (resets in %s)\n",
				rl.ComplexityRemaining, rl.ComplexityLimit, cplxResetIn)
			fmt.Printf("Last query complexity: %d\n", rl.Complexity)
		} else {
			fmt.Printf("\n%s API Rate Limits\n\n",
				color.New(color.FgCyan, color.Bold).Sprint("📊"))

			reqPct := float64(rl.RequestRemaining) / float64(max(rl.RequestLimit, 1)) * 100
			cplxPct := float64(rl.ComplexityRemaining) / float64(max(rl.ComplexityLimit, 1)) * 100

			reqColor := color.New(color.FgGreen)
			if reqPct < 20 {
				reqColor = color.New(color.FgRed)
			} else if reqPct < 50 {
				reqColor = color.New(color.FgYellow)
			}

			cplxColor := color.New(color.FgGreen)
			if cplxPct < 20 {
				cplxColor = color.New(color.FgRed)
			} else if cplxPct < 50 {
				cplxColor = color.New(color.FgYellow)
			}

			fmt.Printf("  Requests:   %s / %d  (resets in %s)\n",
				reqColor.Sprintf("%d", rl.RequestRemaining),
				rl.RequestLimit, reqResetIn)
			fmt.Printf("  Complexity: %s / %d  (resets in %s)\n",
				cplxColor.Sprintf("%d", rl.ComplexityRemaining),
				rl.ComplexityLimit, cplxResetIn)
			fmt.Printf("  Last query: %d complexity points\n", rl.Complexity)
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(rateLimitCmd)

	// Add whoami as a top-level command too
	rootCmd.AddCommand(whoamiCmd)
}
