package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
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
  linear-cli user list              # List all users
  linear-cli user get john@example.com  # Get user details
  linear-cli user me                # Show current user`,
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
			fmt.Println("Name\tEmail\tRole\tActive\tGuest\tTimezone\tStatus")
			for _, user := range filteredUsers {
				role := "Member"
				if user.Owner {
					role = "Owner"
				} else if user.Admin {
					role = "Admin"
				}
				status := ""
				if user.StatusEmoji != "" || user.StatusLabel != "" {
					status = user.StatusEmoji + " " + user.StatusLabel
				}
				fmt.Printf("%s\t%s\t%s\t%v\t%v\t%s\t%s\n",
					user.Name,
					user.Email,
					role,
					user.Active,
					user.Guest,
					user.Timezone,
					status,
				)
			}
		} else {
			// Table output
			headers := []string{"Name", "Email", "Role", "Status", "Timezone"}
			rows := [][]string{}

			for _, user := range filteredUsers {
				role := "Member"
				roleColor := color.New(color.FgWhite)
				if user.Owner {
					role = "Owner"
					roleColor = color.New(color.FgMagenta, color.Bold)
				} else if user.Admin {
					role = "Admin"
					roleColor = color.New(color.FgYellow)
				}
				if user.Guest {
					role = "Guest"
					roleColor = color.New(color.FgHiBlack)
				}
				if user.IsMe {
					role = role + " (You)"
					roleColor = color.New(color.FgCyan, color.Bold)
				}

				status := color.New(color.FgGreen).Sprint("✓ Active")
				if !user.Active {
					status = color.New(color.FgRed).Sprint("✗ Inactive")
				}
				// Show user status if set
				if user.StatusEmoji != "" || user.StatusLabel != "" {
					userStatus := strings.TrimSpace(user.StatusEmoji + " " + user.StatusLabel)
					status = status + " " + color.New(color.FgHiBlack).Sprint("("+userStatus+")")
				}

				timezone := user.Timezone
				if timezone == "" {
					timezone = "-"
				}

				rows = append(rows, []string{
					user.Name,
					color.New(color.FgCyan).Sprint(user.Email),
					roleColor.Sprint(role),
					status,
					timezone,
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
	Use:     "get EMAIL_OR_ID",
	Aliases: []string{"show"},
	Short:   "Get user details",
	Long:    `Get detailed information about a specific user by email or ID.`,
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
			if user.DisplayName != "" && user.DisplayName != user.Name {
				fmt.Printf("Display Name: %s\n", user.DisplayName)
			}
			fmt.Printf("Email: %s\n", user.Email)
			fmt.Printf("Admin: %v\n", user.Admin)
			fmt.Printf("Owner: %v\n", user.Owner)
			fmt.Printf("Active: %v\n", user.Active)
			fmt.Printf("Guest: %v\n", user.Guest)
			if user.Timezone != "" {
				fmt.Printf("Timezone: %s\n", user.Timezone)
			}
			if user.Description != "" {
				fmt.Printf("Description: %s\n", user.Description)
			}
			if user.StatusEmoji != "" || user.StatusLabel != "" {
				fmt.Printf("Status: %s %s\n", user.StatusEmoji, user.StatusLabel)
			}
			if user.URL != "" {
				fmt.Printf("URL: %s\n", user.URL)
			}
			if user.AvatarURL != "" {
				fmt.Printf("Avatar: %s\n", user.AvatarURL)
			}
			if user.LastSeen != nil {
				fmt.Printf("Last Seen: %s\n", user.LastSeen.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("Issues Created: %d\n", user.CreatedIssueCount)
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("👤 User:"),
				user.Name)
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Email:"),
				color.New(color.FgCyan).Sprint(user.Email))
			if user.DisplayName != "" && user.DisplayName != user.Name {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Display Name:"), user.DisplayName)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), user.ID)

			role := "Member"
			roleColor := color.New(color.FgWhite)
			if user.Owner {
				role = "Owner"
				roleColor = color.New(color.FgMagenta, color.Bold)
			} else if user.Admin {
				role = "Admin"
				roleColor = color.New(color.FgYellow)
			}
			if user.Guest {
				role = "Guest"
				roleColor = color.New(color.FgHiBlack)
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
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Active:"), status)

			// User status (emoji + label)
			if user.StatusEmoji != "" || user.StatusLabel != "" {
				userStatus := strings.TrimSpace(user.StatusEmoji + " " + user.StatusLabel)
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), userStatus)
				if user.StatusUntilAt != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status Until:"),
						user.StatusUntilAt.Format("2006-01-02 15:04"))
				}
			}

			if user.Timezone != "" {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Timezone:"), user.Timezone)
			}

			if user.Description != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Description:"),
					color.New(color.FgHiBlack).Sprint(user.Description))
			}

			if user.LastSeen != nil {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Last Seen:"),
					user.LastSeen.Format("2006-01-02 15:04:05"))
			}

			fmt.Printf("%s %d\n", color.New(color.Bold).Sprint("Issues Created:"), user.CreatedIssueCount)

			if user.URL != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Profile URL:"),
					color.New(color.FgBlue).Sprint(user.URL))
			}

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
			if user.DisplayName != "" && user.DisplayName != user.Name {
				fmt.Printf("Display Name: %s\n", user.DisplayName)
			}
			fmt.Printf("Email: %s\n", user.Email)
			fmt.Printf("Admin: %v\n", user.Admin)
			fmt.Printf("Owner: %v\n", user.Owner)
			fmt.Printf("Active: %v\n", user.Active)
			fmt.Printf("Guest: %v\n", user.Guest)
			if user.Timezone != "" {
				fmt.Printf("Timezone: %s\n", user.Timezone)
			}
			if user.Description != "" {
				fmt.Printf("Description: %s\n", user.Description)
			}
			if user.StatusEmoji != "" || user.StatusLabel != "" {
				fmt.Printf("Status: %s %s\n", user.StatusEmoji, user.StatusLabel)
			}
			if user.URL != "" {
				fmt.Printf("URL: %s\n", user.URL)
			}
			if user.AvatarURL != "" {
				fmt.Printf("Avatar: %s\n", user.AvatarURL)
			}
			if user.LastSeen != nil {
				fmt.Printf("Last Seen: %s\n", user.LastSeen.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("Issues Created: %d\n", user.CreatedIssueCount)
		} else {
			// Formatted output
			fmt.Println()
			fmt.Printf("%s %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("👤 Current User:"),
				user.Name)
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Email:"),
				color.New(color.FgCyan).Sprint(user.Email))
			if user.DisplayName != "" && user.DisplayName != user.Name {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Display Name:"), user.DisplayName)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), user.ID)

			role := "Member"
			roleColor := color.New(color.FgWhite)
			if user.Owner {
				role = "Owner"
				roleColor = color.New(color.FgMagenta, color.Bold)
			} else if user.Admin {
				role = "Admin"
				roleColor = color.New(color.FgYellow, color.Bold)
			}
			if user.Guest {
				role = "Guest"
				roleColor = color.New(color.FgHiBlack)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Role:"), roleColor.Sprint(role))

			status := color.New(color.FgGreen).Sprint("✓ Active")
			if !user.Active {
				status = color.New(color.FgRed).Sprint("✗ Inactive")
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Active:"), status)

			// User status (emoji + label)
			if user.StatusEmoji != "" || user.StatusLabel != "" {
				userStatus := strings.TrimSpace(user.StatusEmoji + " " + user.StatusLabel)
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), userStatus)
				if user.StatusUntilAt != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status Until:"),
						user.StatusUntilAt.Format("2006-01-02 15:04"))
				}
			}

			if user.Timezone != "" {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Timezone:"), user.Timezone)
			}

			if user.Description != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Description:"),
					color.New(color.FgHiBlack).Sprint(user.Description))
			}

			if user.LastSeen != nil {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Last Seen:"),
					user.LastSeen.Format("2006-01-02 15:04:05"))
			}

			fmt.Printf("%s %d\n", color.New(color.Bold).Sprint("Issues Created:"), user.CreatedIssueCount)

			if user.URL != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Profile URL:"),
					color.New(color.FgBlue).Sprint(user.URL))
			}

			if user.AvatarURL != "" {
				fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Avatar:"),
					color.New(color.FgBlue).Sprint(user.AvatarURL))
			}
			fmt.Println()
		}
	},
}

var userUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Short:   "Update your user profile",
	Long: `Update your own user profile settings.

Note: You can only update your own profile. Linear does not allow modifying other users' profiles via API.

Examples:
  linear-cli user update --name "John Doe"
  linear-cli user update --status-emoji "🏖️" --status-label "On vacation"
  linear-cli user update --timezone "America/New_York"
  linear-cli user update --status-emoji "" --status-label ""  # Clear status`,
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

		// Get current user to get their ID
		viewer, err := client.GetViewer(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get current user: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Build update input
		input := api.UserUpdateInput{}
		hasUpdates := false

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			input.Name = &name
			hasUpdates = true
		}
		if cmd.Flags().Changed("display-name") {
			displayName, _ := cmd.Flags().GetString("display-name")
			input.DisplayName = &displayName
			hasUpdates = true
		}
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			input.Description = &description
			hasUpdates = true
		}
		if cmd.Flags().Changed("avatar-url") {
			avatarURL, _ := cmd.Flags().GetString("avatar-url")
			input.AvatarURL = &avatarURL
			hasUpdates = true
		}
		if cmd.Flags().Changed("timezone") {
			timezone, _ := cmd.Flags().GetString("timezone")
			input.Timezone = &timezone
			hasUpdates = true
		}
		if cmd.Flags().Changed("status-emoji") {
			statusEmoji, _ := cmd.Flags().GetString("status-emoji")
			input.StatusEmoji = &statusEmoji
			hasUpdates = true
		}
		if cmd.Flags().Changed("status-label") {
			statusLabel, _ := cmd.Flags().GetString("status-label")
			input.StatusLabel = &statusLabel
			hasUpdates = true
		}
		if cmd.Flags().Changed("status-until") {
			statusUntil, _ := cmd.Flags().GetString("status-until")
			if statusUntil != "" {
				// Parse as date (YYYY-MM-DD) or datetime (YYYY-MM-DDTHH:MM:SS)
				var t time.Time
				var parseErr error
				if len(statusUntil) == 10 {
					t, parseErr = time.Parse("2006-01-02", statusUntil)
				} else {
					t, parseErr = time.Parse(time.RFC3339, statusUntil)
				}
				if parseErr != nil {
					output.Error(fmt.Sprintf("Invalid status-until format: %v. Use YYYY-MM-DD or RFC3339 format.", parseErr), plaintext, jsonOut)
					os.Exit(1)
				}
				input.StatusUntilAt = &t
			}
			hasUpdates = true
		}

		if !hasUpdates {
			output.Error("No updates specified. Use --help to see available flags.", plaintext, jsonOut)
			os.Exit(1)
		}

		// Update user
		user, err := client.UpdateUser(context.Background(), viewer.ID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update user: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Handle output
		if jsonOut {
			output.JSON(user)
		} else if plaintext {
			fmt.Printf("Updated user: %s\n", user.Name)
			fmt.Printf("ID: %s\n", user.ID)
			fmt.Printf("Email: %s\n", user.Email)
			if user.DisplayName != "" {
				fmt.Printf("Display Name: %s\n", user.DisplayName)
			}
			if user.Timezone != "" {
				fmt.Printf("Timezone: %s\n", user.Timezone)
			}
			if user.StatusEmoji != "" || user.StatusLabel != "" {
				fmt.Printf("Status: %s %s\n", user.StatusEmoji, user.StatusLabel)
			}
		} else {
			fmt.Printf("\n%s Updated user profile\n",
				color.New(color.FgGreen).Sprint("✅"))
			fmt.Println(strings.Repeat("─", 50))

			fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Name:"), user.Name)
			if user.DisplayName != "" && user.DisplayName != user.Name {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Display Name:"), user.DisplayName)
			}
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Email:"),
				color.New(color.FgCyan).Sprint(user.Email))

			if user.Description != "" {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Description:"), user.Description)
			}

			if user.Timezone != "" {
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Timezone:"), user.Timezone)
			}

			if user.StatusEmoji != "" || user.StatusLabel != "" {
				userStatus := strings.TrimSpace(user.StatusEmoji + " " + user.StatusLabel)
				fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status:"), userStatus)
				if user.StatusUntilAt != nil {
					fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Status Until:"),
						user.StatusUntilAt.Format("2006-01-02 15:04"))
				}
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
	userCmd.AddCommand(userUpdateCmd)

	// List command flags
	userListCmd.Flags().IntP("limit", "l", 50, "Maximum number of users to return")
	userListCmd.Flags().BoolP("active", "a", false, "Show only active users")
	userListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")

	// Update command flags (all fields from UserUpdateInput)
	userUpdateCmd.Flags().String("name", "", "Update name")
	userUpdateCmd.Flags().String("display-name", "", "Update display name")
	userUpdateCmd.Flags().String("description", "", "Update description/bio")
	userUpdateCmd.Flags().String("avatar-url", "", "Update avatar URL")
	userUpdateCmd.Flags().String("timezone", "", "Update timezone (e.g., America/New_York)")
	userUpdateCmd.Flags().String("status-emoji", "", "Set status emoji (empty string to clear)")
	userUpdateCmd.Flags().String("status-label", "", "Set status label (empty string to clear)")
	userUpdateCmd.Flags().String("status-until", "", "Status expiration (YYYY-MM-DD or RFC3339)")
}
