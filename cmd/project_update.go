package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"update-status"},
	Short:   "Manage project status updates",
	Long: `Create, list, view, update, and archive project status updates.

Status updates are periodic check-ins on project health with a body and health indicator.

Health values: onTrack, atRisk, offTrack

Examples:
  linear-cli project status list PROJECT-ID
  linear-cli project status get UPDATE-ID
  linear-cli project status create PROJECT-ID --body "On track for launch" --health onTrack
  linear-cli project status update UPDATE-ID --body "Updated status"
  linear-cli project status delete UPDATE-ID`,
}

var statusListCmd = &cobra.Command{
	Use:     "list [project-id]",
	Aliases: []string{"ls"},
	Short:   "List status updates for a project",
	Long:    `List all status updates for a project, ordered by most recent.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 20
		}

		updates, err := client.GetProjectUpdates(context.Background(), args[0], limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch project updates: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(updates.Nodes) == 0 {
			output.Info("No status updates found for this project", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(updates.Nodes)
			return
		}

		if plaintext {
			fmt.Println("# Project Status Updates")
			for _, u := range updates.Nodes {
				fmt.Printf("\n## %s by %s\n", u.CreatedAt.Format("2006-01-02 15:04"), safeUserName(u.User))
				fmt.Printf("- **ID**: %s\n", u.ID)
				fmt.Printf("- **Health**: %s\n", u.Health)
				if u.EditedAt != nil {
					fmt.Printf("- **Edited**: %s\n", u.EditedAt.Format("2006-01-02 15:04"))
				}
				if u.URL != "" {
					fmt.Printf("- **URL**: %s\n", u.URL)
				}
				fmt.Printf("\n%s\n", u.Body)
			}
			fmt.Printf("\nTotal: %d updates\n", len(updates.Nodes))
			return
		}

		// Table output
		headers := []string{"Date", "Health", "Author", "Preview", "ID"}
		rows := make([][]string, len(updates.Nodes))
		for i, u := range updates.Nodes {
			healthStr := formatHealth(u.Health)
			author := ""
			if u.User != nil {
				author = u.User.Name
			}

			// Preview first line of body
			preview := u.Body
			if idx := strings.Index(preview, "\n"); idx >= 0 {
				preview = preview[:idx]
			}
			preview = truncateString(preview, 50)

			rows[i] = []string{
				u.CreatedAt.Format("2006-01-02"),
				healthStr,
				author,
				preview,
				truncateString(u.ID, 12),
			}
		}

		tableData := output.TableData{
			Headers: headers,
			Rows:    rows,
		}
		output.Table(tableData, false, false)

		fmt.Printf("\n%s %d status updates\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(updates.Nodes))
	},
}

var statusGetCmd = &cobra.Command{
	Use:     "get [update-id]",
	Aliases: []string{"show"},
	Short:   "Get a project status update",
	Long:    `Get full details of a project status update.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		update, err := client.GetProjectUpdate(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(update)
			return
		}

		if plaintext {
			fmt.Printf("# Project Status Update\n\n")
			fmt.Printf("- **ID**: %s\n", update.ID)
			fmt.Printf("- **Health**: %s\n", update.Health)
			if update.User != nil {
				fmt.Printf("- **Author**: %s (%s)\n", update.User.Name, update.User.Email)
			}
			if update.Project != nil {
				fmt.Printf("- **Project**: %s\n", update.Project.Name)
			}
			fmt.Printf("- **Created**: %s\n", update.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", update.UpdatedAt.Format("2006-01-02 15:04:05"))
			if update.EditedAt != nil {
				fmt.Printf("- **Edited**: %s\n", update.EditedAt.Format("2006-01-02 15:04:05"))
			}
			if update.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", update.ArchivedAt.Format("2006-01-02 15:04:05"))
			}
			if update.URL != "" {
				fmt.Printf("- **URL**: %s\n", update.URL)
			}
			fmt.Printf("\n## Body\n%s\n", update.Body)
			return
		}

		// Rich output
		fmt.Println()
		if update.Project != nil {
			fmt.Printf("%s %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("Project:"),
				update.Project.Name)
		}
		fmt.Printf("%s %s\n",
			color.New(color.Bold).Sprint("Health:"),
			formatHealth(update.Health))
		if update.User != nil {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("Author:"),
				color.New(color.FgCyan).Sprint(update.User.Name))
		}
		fmt.Printf("%s %s\n",
			color.New(color.Bold).Sprint("Created:"),
			update.CreatedAt.Format("2006-01-02 15:04:05"))
		if update.EditedAt != nil {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("Edited:"),
				update.EditedAt.Format("2006-01-02 15:04:05"))
		}
		if update.URL != "" {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("URL:"),
				color.New(color.FgBlue, color.Underline).Sprint(update.URL))
		}
		fmt.Printf("%s %s\n",
			color.New(color.Bold).Sprint("ID:"),
			color.New(color.FgWhite, color.Faint).Sprint(update.ID))
		fmt.Printf("\n%s\n", update.Body)
		fmt.Println()
	},
}

var statusCreateCmd = &cobra.Command{
	Use:     "create [project-id]",
	Aliases: []string{"new", "add"},
	Short:   "Create a project status update",
	Long: `Create a new status update on a project.

Health values: onTrack, atRisk, offTrack

Examples:
  linear-cli project status create PROJECT-ID --body "Sprint going well" --health onTrack
  linear-cli project status create PROJECT-ID --body "Blocked on dependency" --health atRisk`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		body, _ := cmd.Flags().GetString("body")
		health, _ := cmd.Flags().GetString("health")

		if body == "" {
			output.Error("--body is required", plaintext, jsonOut)
			os.Exit(1)
		}

		if health != "" {
			if !isValidHealth(health) {
				output.Error(fmt.Sprintf("Invalid health value '%s'. Valid values: onTrack, atRisk, offTrack", health), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		input := map[string]interface{}{
			"projectId": args[0],
			"body":      body,
		}
		if health != "" {
			input["health"] = health
		}

		update, err := client.CreateProjectUpdate(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(update)
		} else if plaintext {
			projectName := args[0]
			if update.Project != nil {
				projectName = update.Project.Name
			}
			fmt.Printf("Created status update on %s (health: %s)\n", projectName, update.Health)
			fmt.Printf("ID: %s\n", update.ID)
		} else {
			projectName := args[0]
			if update.Project != nil {
				projectName = update.Project.Name
			}
			fmt.Printf("%s Created status update on %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(projectName))
			fmt.Printf("  Health: %s\n", formatHealth(update.Health))
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(update.ID))
			if update.URL != "" {
				fmt.Printf("  URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(update.URL))
			}
		}
	},
}

var statusUpdateCmd = &cobra.Command{
	Use:     "update [update-id]",
	Aliases: []string{"edit"},
	Short:   "Update a project status update",
	Long: `Update the body or health of an existing status update.

Examples:
  linear-cli project status update UPDATE-ID --body "Updated status text"
  linear-cli project status update UPDATE-ID --health offTrack
  linear-cli project status update UPDATE-ID --body "New text" --health onTrack`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		input := make(map[string]interface{})

		if cmd.Flags().Changed("body") {
			body, _ := cmd.Flags().GetString("body")
			input["body"] = body
		}

		if cmd.Flags().Changed("health") {
			health, _ := cmd.Flags().GetString("health")
			if !isValidHealth(health) {
				output.Error(fmt.Sprintf("Invalid health value '%s'. Valid values: onTrack, atRisk, offTrack", health), plaintext, jsonOut)
				os.Exit(1)
			}
			input["health"] = health
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use --body and/or --health.", plaintext, jsonOut)
			os.Exit(1)
		}

		update, err := client.UpdateProjectUpdate(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(update)
		} else if plaintext {
			fmt.Printf("Updated project status update %s\n", update.ID)
			fmt.Printf("Health: %s\n", update.Health)
		} else {
			fmt.Printf("%s Updated project status update\n",
				color.New(color.FgGreen).Sprint("✓"))
			fmt.Printf("  Health: %s\n", formatHealth(update.Health))
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(update.ID))
		}
	},
}

var statusDeleteCmd = &cobra.Command{
	Use:     "delete [update-id]",
	Aliases: []string{"archive", "rm"},
	Short:   "Archive a project status update",
	Long: `Archive (soft-delete) a project status update.

Note: Linear does not support permanent deletion of project updates.
This command archives the update instead.

Examples:
  linear-cli project status delete UPDATE-ID`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		err = client.ArchiveProjectUpdate(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to archive project update: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "id": args[0], "action": "archived"})
		} else if plaintext {
			fmt.Printf("Archived project status update %s\n", args[0])
		} else {
			fmt.Printf("%s Archived project status update %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Faint).Sprint(args[0]))
		}
	},
}

func formatHealth(health string) string {
	switch health {
	case "onTrack":
		return color.New(color.FgGreen).Sprint("On Track")
	case "atRisk":
		return color.New(color.FgYellow).Sprint("At Risk")
	case "offTrack":
		return color.New(color.FgRed).Sprint("Off Track")
	default:
		return health
	}
}

func isValidHealth(health string) bool {
	switch health {
	case "onTrack", "atRisk", "offTrack":
		return true
	default:
		return false
	}
}

func init() {
	projectCmd.AddCommand(projectStatusCmd)
	projectStatusCmd.AddCommand(statusListCmd)
	projectStatusCmd.AddCommand(statusGetCmd)
	projectStatusCmd.AddCommand(statusCreateCmd)
	projectStatusCmd.AddCommand(statusUpdateCmd)
	projectStatusCmd.AddCommand(statusDeleteCmd)

	// list flags
	statusListCmd.Flags().IntP("limit", "l", 20, "Maximum number of updates to fetch")

	// create flags
	statusCreateCmd.Flags().StringP("body", "b", "", "Status update body text (required)")
	statusCreateCmd.Flags().String("health", "", "Project health: onTrack, atRisk, offTrack")
	_ = statusCreateCmd.MarkFlagRequired("body")

	// update flags
	statusUpdateCmd.Flags().StringP("body", "b", "", "New body text")
	statusUpdateCmd.Flags().String("health", "", "New health: onTrack, atRisk, offTrack")

	// Aliases so "linear-cli project update-status" also works — handled via Aliases on projectStatusCmd
}
