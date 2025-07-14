package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dorkitude/linctl/pkg/api"
	"github.com/dorkitude/linctl/pkg/auth"
	"github.com/dorkitude/linctl/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Linear users",
	Long: `Manage Linear users including listing users, viewing user details, and showing the current user.

Examples:
  linctl user list              # List all users
  linctl user get john@example.com  # Get user details
  linctl user me                # Show current user`,
}

var userListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List users",
	Long:    `List all users in your Linear workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get filters
		limit, _ := cmd.Flags().GetInt("limit")
		activeOnly, _ := cmd.Flags().GetBool("active")

		// Get sort option
		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				// Use empty string for Linear's default sort
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		// Get users
		users, err := client.GetUsers(context.Background(), limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list users: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Filter active users if requested
		filteredUsers := users.Nodes
		if activeOnly {
			var activeUsers []api.User
			for _, user := range users.Nodes {
				if user.Active {
					activeUsers = append(activeUsers, user)
				}
			}
			filteredUsers = activeUsers
		}

		// Handle output
		if jsonOut {
			output.JSON(filteredUsers)
		} else if plaintext {
			fmt.Println("Name\tEmail\tRole\tActive")
			for _, user := range filteredUsers {
				role := "Member"
				if user.Admin {
					role = "Admin"
				}
				fmt.Printf("%s\t%s\t%s\t%v\n",
					user.Name,
					user.Email,
					role,
					user.Active,
				)
			}
		} else {
			// Table output
			headers := []string{"Name", "Email", "Role", "Status"}
			rows := [][]string{}

			for _, user := range filteredUsers {
				role := "Member"
				roleColor := color.New(color.FgWhite)
				if user.Admin {
					role = "Admin"
					roleColor = color.New(color.FgYellow)
				}
				if user.IsMe {
					role = role + " (You)"
					roleColor = color.New(color.FgCyan, color.Bold)
				}

				status := color.New(color.FgGreen).Sprint("✓ Active")
				if !user.Active {
					status = color.New(color.FgRed).Sprint("✗ Inactive")
				}

				rows = append(rows, []string{
					user.Name,
					color.New(color.FgCyan).Sprint(user.Email),
					roleColor.Sprint(role),
					status,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			if !plaintext && !jsonOut {
				fmt.Printf("\n%s %d users\n",
					color.New(color.FgGreen).Sprint("✓"),
					len(filteredUsers))
			}
		}
	},
}

var userGetCmd = &cobra.Command{
	Use:     "get EMAIL",
	Aliases: []string{"show"},
	Short:   "Get user details",
	Long:    `Get detailed information about a specific user by email.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		email := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get user details
		user, err := client.GetUser(context.Background(), email)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get user: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(user)
		} else if plaintext {
			fmt.Printf("ID: %s\n", user.ID)
			fmt.Printf("Name: %s\n", user.Name)
			fmt.Printf("Email: %s\n", user.Email)
			fmt.Printf("Admin: %v\n", user.Admin)
			fmt.Printf("Active: %v\n", user.Active)
			if user.AvatarURL != "" {
				fmt.Printf("Avatar: %s\n", user.AvatarURL)
			}
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("👤 User:"),
				user.Name)
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Email:"),
				color.New(color.FgCyan).Sprint(user.Email))
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), user.ID)

			role := "Member"
			roleColor := color.New(color.FgWhite)
			if user.Admin {
				role = "Admin"
				roleColor = color.New(color.FgYellow)
			}
			if user.IsMe {
				role = role + " (You)"
				roleColor = color.New(color.FgCyan, color.Bold)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Role:"), roleColor.Sprint(role))

			status := color.New(color.FgGreen).Sprint("✓ Active")
			if !user.Active {
				status = color.New(color.FgRed).Sprint("✗ Inactive")
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), status)

			if user.AvatarURL != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Avatar:"),
					color.New(color.FgBlue).Sprint(user.AvatarURL))
			}
			fmt.Println()
		}
	},
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user",
	Long:  `Display information about the currently authenticated user.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Get current user
		user, err := client.GetViewer(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(user)
		} else if plaintext {
			fmt.Printf("ID: %s\n", user.ID)
			fmt.Printf("Name: %s\n", user.Name)
			fmt.Printf("Email: %s\n", user.Email)
			fmt.Printf("Admin: %v\n", user.Admin)
			fmt.Printf("Active: %v\n", user.Active)
			if user.AvatarURL != "" {
				fmt.Printf("Avatar: %s\n", user.AvatarURL)
			}
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("👤 Current User:"),
				user.Name)
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Email:"),
				color.New(color.FgCyan).Sprint(user.Email))
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), user.ID)

			role := "Member"
			roleColor := color.New(color.FgWhite)
			if user.Admin {
				role = "Admin"
				roleColor = color.New(color.FgYellow, color.Bold)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Role:"), roleColor.Sprint(role))

			status := color.New(color.FgGreen).Sprint("✓ Active")
			if !user.Active {
				status = color.New(color.FgRed).Sprint("✗ Inactive")
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), status)

			if user.AvatarURL != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Avatar:"),
					color.New(color.FgBlue).Sprint(user.AvatarURL))
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userMeCmd)

	// List command flags
	userListCmd.Flags().IntP("limit", "l", 50, "Maximum number of users to return")
	userListCmd.Flags().BoolP("active", "a", false, "Show only active users")
	userListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
}
