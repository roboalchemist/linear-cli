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

var initiativeCmd = &cobra.Command{
	Use:   "initiative",
	Short: "Manage Linear initiatives",
	Long: `Manage Linear initiatives (high-level strategic objectives that group projects).

Examples:
  linear-cli initiative list
  linear-cli initiative list --status Active
  linear-cli initiative get INITIATIVE-ID
  linear-cli initiative create --name "Q1 Goals"
  linear-cli initiative update INITIATIVE-ID --status Completed
  linear-cli initiative delete INITIATIVE-ID`,
}

var initiativeListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List initiatives",
	Long:    `List Linear initiatives with optional filtering.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		filter := make(map[string]interface{})
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			filter["status"] = map[string]interface{}{"eq": status}
		}

		includeCompleted, _ := cmd.Flags().GetBool("include-completed")
		if !includeCompleted {
			if _, hasStatus := filter["status"]; !hasStatus {
				filter["status"] = map[string]interface{}{
					"nin": []string{"Completed"},
				}
			}
		}

		limit, _ := cmd.Flags().GetInt("limit")

		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		initiatives, err := client.GetInitiatives(context.Background(), filter, limit, "", orderBy, includeCompleted)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch initiatives: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(initiatives.Nodes) == 0 {
			output.Info("No initiatives found", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(initiatives.Nodes)
			return
		}

		if plaintext {
			fmt.Println("# Initiatives")
			for _, init := range initiatives.Nodes {
				fmt.Printf("## %s\n", init.Name)
				fmt.Printf("- **ID**: %s\n", init.ID)
				fmt.Printf("- **Status**: %s\n", init.Status)
				if init.Health != "" {
					fmt.Printf("- **Health**: %s\n", init.Health)
				}
				if init.Owner != nil {
					fmt.Printf("- **Owner**: %s\n", init.Owner.Name)
				}
				if init.TargetDate != nil {
					fmt.Printf("- **Target Date**: %s\n", *init.TargetDate)
				}
				fmt.Printf("- **Created**: %s\n", init.CreatedAt.Format("2006-01-02"))
				fmt.Printf("- **URL**: %s\n", init.URL)
				if init.Description != "" {
					fmt.Printf("- **Description**: %s\n", init.Description)
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d initiatives\n", len(initiatives.Nodes))
			return
		}

		headers := []string{"Name", "Status", "Health", "Owner", "Target Date", "URL"}
		rows := make([][]string, len(initiatives.Nodes))

		for i, init := range initiatives.Nodes {
			owner := color.New(color.FgYellow).Sprint("Unassigned")
			if init.Owner != nil {
				owner = init.Owner.Name
			}

			statusColor := color.New(color.FgWhite)
			switch init.Status {
			case "Planned":
				statusColor = color.New(color.FgCyan)
			case "Active":
				statusColor = color.New(color.FgBlue)
			case "Completed":
				statusColor = color.New(color.FgGreen)
			}

			healthStr := ""
			if init.Health != "" {
				healthStr = init.Health
			}

			targetDate := ""
			if init.TargetDate != nil {
				targetDate = *init.TargetDate
			}

			rows[i] = []string{
				truncateString(init.Name, 30),
				statusColor.Sprint(init.Status),
				healthStr,
				owner,
				targetDate,
				init.URL,
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, false, false)

		fmt.Printf("\n%s %d initiatives\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(initiatives.Nodes))

		if initiatives.PageInfo.HasNextPage {
			fmt.Printf("%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

var initiativeGetCmd = &cobra.Command{
	Use:     "get [initiative-id]",
	Aliases: []string{"show"},
	Short:   "Get initiative details",
	Long:    `Get detailed information about a specific initiative.`,
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
		initiative, err := client.GetInitiative(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(initiative)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n\n", initiative.Name)

			if initiative.Description != "" {
				fmt.Printf("## Description\n%s\n\n", initiative.Description)
			}

			if initiative.Content != "" {
				fmt.Printf("## Content\n%s\n\n", initiative.Content)
			}

			fmt.Printf("## Core Details\n")
			fmt.Printf("- **ID**: %s\n", initiative.ID)
			fmt.Printf("- **Status**: %s\n", initiative.Status)
			if initiative.Health != "" {
				fmt.Printf("- **Health**: %s\n", initiative.Health)
			}
			if initiative.TargetDate != nil {
				fmt.Printf("- **Target Date**: %s\n", *initiative.TargetDate)
			}
			if initiative.TargetDateResolution != "" {
				fmt.Printf("- **Target Date Resolution**: %s\n", initiative.TargetDateResolution)
			}
			if initiative.Icon != nil && *initiative.Icon != "" {
				fmt.Printf("- **Icon**: %s\n", *initiative.Icon)
			}
			fmt.Printf("- **Color**: %s\n", initiative.Color)

			fmt.Printf("\n## People\n")
			if initiative.Owner != nil {
				fmt.Printf("- **Owner**: %s (%s)\n", initiative.Owner.Name, initiative.Owner.Email)
			} else {
				fmt.Printf("- **Owner**: Unassigned\n")
			}
			if initiative.Creator != nil {
				fmt.Printf("- **Creator**: %s (%s)\n", initiative.Creator.Name, initiative.Creator.Email)
			}

			fmt.Printf("\n## Timeline\n")
			fmt.Printf("- **Created**: %s\n", initiative.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", initiative.UpdatedAt.Format("2006-01-02 15:04:05"))
			if initiative.CompletedAt != nil {
				fmt.Printf("- **Completed**: %s\n", initiative.CompletedAt.Format("2006-01-02 15:04:05"))
			}
			if initiative.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", initiative.ArchivedAt.Format("2006-01-02 15:04:05"))
			}

			fmt.Printf("\n## URL\n- %s\n", initiative.URL)

			if initiative.ParentInitiative != nil {
				fmt.Printf("\n## Parent Initiative\n")
				fmt.Printf("- %s (%s)\n", initiative.ParentInitiative.Name, initiative.ParentInitiative.Status)
			}

			if initiative.SubInitiatives != nil && len(initiative.SubInitiatives.Nodes) > 0 {
				fmt.Printf("\n## Sub-Initiatives\n")
				for _, sub := range initiative.SubInitiatives.Nodes {
					fmt.Printf("- %s [%s]", sub.Name, sub.Status)
					if sub.Health != "" {
						fmt.Printf(" (%s)", sub.Health)
					}
					fmt.Println()
				}
			}

			if initiative.Projects != nil && len(initiative.Projects.Nodes) > 0 {
				fmt.Printf("\n## Linked Projects\n")
				for _, proj := range initiative.Projects.Nodes {
					fmt.Printf("- %s [%s] %.0f%%\n", proj.Name, proj.State, proj.Progress*100)
				}
			}

			return
		}

		// Rich display
		fmt.Println()
		fmt.Printf("%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("Initiative:"),
			color.New(color.FgWhite, color.Bold).Sprint(initiative.Name))
		fmt.Println(strings.Repeat("─", 50))

		if initiative.Description != "" {
			fmt.Printf("\n%s\n", initiative.Description)
		}

		statusColor := color.New(color.FgWhite)
		switch initiative.Status {
		case "Planned":
			statusColor = color.New(color.FgCyan)
		case "Active":
			statusColor = color.New(color.FgBlue)
		case "Completed":
			statusColor = color.New(color.FgGreen)
		}
		fmt.Printf("\n%s %s\n", color.New(color.Bold).Sprint("Status:"), statusColor.Sprint(initiative.Status))

		if initiative.Health != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Health:"), initiative.Health)
		}

		if initiative.Owner != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("Owner:"),
				initiative.Owner.Name,
				color.New(color.FgCyan).Sprint(initiative.Owner.Email))
		}

		if initiative.TargetDate != nil {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Target Date:"), *initiative.TargetDate)
		}

		fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Timeline:"))
		fmt.Printf("  Created: %s\n", initiative.CreatedAt.Format("2006-01-02"))
		fmt.Printf("  Updated: %s\n", initiative.UpdatedAt.Format("2006-01-02"))
		if initiative.CompletedAt != nil {
			fmt.Printf("  Completed: %s\n", initiative.CompletedAt.Format("2006-01-02"))
		}

		if initiative.URL != "" {
			fmt.Printf("\n%s %s\n",
				color.New(color.Bold).Sprint("URL:"),
				color.New(color.FgBlue, color.Underline).Sprint(initiative.URL))
		}

		if initiative.ParentInitiative != nil {
			fmt.Printf("\n%s %s [%s]\n",
				color.New(color.Bold).Sprint("Parent:"),
				color.New(color.FgCyan).Sprint(initiative.ParentInitiative.Name),
				initiative.ParentInitiative.Status)
		}

		if initiative.SubInitiatives != nil && len(initiative.SubInitiatives.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Sub-Initiatives:"))
			for _, sub := range initiative.SubInitiatives.Nodes {
				fmt.Printf("  - %s [%s]\n",
					color.New(color.FgCyan).Sprint(sub.Name),
					sub.Status)
			}
		}

		if initiative.Projects != nil && len(initiative.Projects.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Linked Projects:"))
			for _, proj := range initiative.Projects.Nodes {
				progressColor := color.New(color.FgRed)
				if proj.Progress >= 0.75 {
					progressColor = color.New(color.FgGreen)
				} else if proj.Progress >= 0.5 {
					progressColor = color.New(color.FgYellow)
				}
				fmt.Printf("  - %s [%s] %s\n",
					color.New(color.FgCyan).Sprint(proj.Name),
					proj.State,
					progressColor.Sprintf("%.0f%%", proj.Progress*100))
			}
		}

		fmt.Println()
	},
}

var initiativeCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new initiative",
	Long: `Create a new initiative in Linear.

The description can be provided inline via --description or read from a markdown file via --description-file.
Use --description-file - to read from stdin.

Examples:
  linear-cli initiative create --name "Q1 Goals"
  linear-cli initiative create --name "Q1 Goals" --description "Details"
  linear-cli initiative create --name "Q1 Goals" --description-file initiative-brief.md`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			output.Error("Name is required (--name)", plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{
			"name": name,
		}

		// Resolve description from --description or --description-file
		filePath, _ := cmd.Flags().GetString("description-file")
		if cmd.Flags().Changed("description") || filePath != "" {
			descFlag, _ := cmd.Flags().GetString("description")
			desc, err := resolveBodyFromFlags(descFlag, cmd.Flags().Changed("description"), filePath, "description", "description-file")
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["description"] = desc
		}
		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			input["status"] = status
		}
		if cmd.Flags().Changed("target-date") {
			td, _ := cmd.Flags().GetString("target-date")
			input["targetDate"] = td
		}
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}
		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			input["icon"] = icon
		}

		initiative, err := client.CreateInitiative(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(initiative)
		} else if plaintext {
			fmt.Printf("Created initiative: %s (%s)\n", initiative.Name, initiative.ID)
		} else {
			fmt.Printf("%s Created initiative %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(initiative.Name))
			fmt.Printf("  ID: %s\n", initiative.ID)
			if initiative.URL != "" {
				fmt.Printf("  URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(initiative.URL))
			}
		}
	},
}

var initiativeUpdateCmd = &cobra.Command{
	Use:     "update [initiative-id]",
	Aliases: []string{"edit"},
	Short:   "Update an initiative",
	Long: `Update various fields of an initiative.

The description can be provided inline via --description or read from a markdown file via --description-file.
Use --description-file - to read from stdin.

Examples:
  linear-cli initiative update ID --name "New name"
  linear-cli initiative update ID --status Active
  linear-cli initiative update ID --description-file updated-brief.md`,
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

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			input["name"] = name
		}
		// Resolve description from --description or --description-file
		filePath, _ := cmd.Flags().GetString("description-file")
		if cmd.Flags().Changed("description") || filePath != "" {
			descFlag, _ := cmd.Flags().GetString("description")
			desc, err := resolveBodyFromFlags(descFlag, cmd.Flags().Changed("description"), filePath, "description", "description-file")
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["description"] = desc
		}
		if cmd.Flags().Changed("status") {
			status, _ := cmd.Flags().GetString("status")
			input["status"] = status
		}
		if cmd.Flags().Changed("target-date") {
			td, _ := cmd.Flags().GetString("target-date")
			if td == "" {
				input["targetDate"] = nil
			} else {
				input["targetDate"] = td
			}
		}
		if cmd.Flags().Changed("color") {
			c, _ := cmd.Flags().GetString("color")
			input["color"] = c
		}
		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			input["icon"] = icon
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use flags to specify what to update.", plaintext, jsonOut)
			os.Exit(1)
		}

		initiative, err := client.UpdateInitiative(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(initiative)
		} else if plaintext {
			fmt.Printf("Updated initiative: %s\n", initiative.Name)
		} else {
			fmt.Printf("%s Updated initiative %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(initiative.Name))
		}
	},
}

var initiativeDeleteCmd = &cobra.Command{
	Use:     "delete [initiative-id]",
	Aliases: []string{"rm"},
	Short:   "Delete an initiative",
	Long:    `Delete an initiative from Linear.`,
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
		err = client.DeleteInitiative(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete initiative: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{
				"status":  "success",
				"message": "Initiative deleted",
			})
		} else if plaintext {
			fmt.Println("Initiative deleted")
		} else {
			fmt.Printf("%s Initiative deleted\n",
				color.New(color.FgGreen).Sprint("✓"))
		}
	},
}

var initiativeProjectsCmd = &cobra.Command{
	Use:     "projects INITIATIVE-ID",
	Aliases: []string{"project"},
	Short:   "List projects under an initiative",
	Long:    `List all projects that belong to a specific initiative.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		initiativeID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")

		projects, err := client.GetInitiativeProjects(context.Background(), initiativeID, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get initiative projects: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(projects.Nodes)
			return
		}

		if len(projects.Nodes) == 0 {
			if plaintext {
				fmt.Println("No projects found")
			} else {
				fmt.Printf("\n%s No projects in this initiative\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Projects")
			fmt.Println("ID\tName\tState\tProgress\tLead")
			for _, p := range projects.Nodes {
				lead := ""
				if p.Lead != nil {
					lead = p.Lead.Name
				}
				fmt.Printf("%s\t%s\t%s\t%.0f%%\t%s\n",
					p.ID, p.Name, p.State, p.Progress*100, lead)
			}
		} else {
			headers := []string{"ID", "Name", "State", "Progress", "Lead"}
			rows := [][]string{}

			for _, p := range projects.Nodes {
				lead := ""
				if p.Lead != nil {
					lead = p.Lead.Name
				}
				stateColor := color.New(color.FgWhite)
				switch p.State {
				case "started":
					stateColor = color.New(color.FgYellow)
				case "completed":
					stateColor = color.New(color.FgGreen)
				case "canceled":
					stateColor = color.New(color.FgRed)
				case "planned":
					stateColor = color.New(color.FgCyan)
				}

				rows = append(rows, []string{
					p.ID[:8],
					color.New(color.FgWhite, color.Bold).Sprint(p.Name),
					stateColor.Sprint(p.State),
					fmt.Sprintf("%.0f%%", p.Progress*100),
					lead,
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
	},
}

func init() {
	rootCmd.AddCommand(initiativeCmd)
	initiativeCmd.AddCommand(initiativeListCmd)
	initiativeCmd.AddCommand(initiativeGetCmd)
	initiativeCmd.AddCommand(initiativeCreateCmd)
	initiativeCmd.AddCommand(initiativeUpdateCmd)
	initiativeCmd.AddCommand(initiativeDeleteCmd)
	initiativeCmd.AddCommand(initiativeProjectsCmd)

	// Initiative projects flags
	initiativeProjectsCmd.Flags().IntP("limit", "l", 50, "Maximum number of projects to return")

	// List flags
	initiativeListCmd.Flags().StringP("status", "s", "", "Filter by status (Planned, Active, Completed)")
	initiativeListCmd.Flags().IntP("limit", "l", 50, "Maximum number of initiatives to fetch")
	initiativeListCmd.Flags().BoolP("include-completed", "c", false, "Include completed initiatives")
	initiativeListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")

	// Create flags
	initiativeCreateCmd.Flags().String("name", "", "Initiative name (required)")
	initiativeCreateCmd.Flags().StringP("description", "d", "", "Initiative description")
	initiativeCreateCmd.Flags().String("description-file", "", "Read description from a markdown file (use - for stdin)")
	initiativeCreateCmd.Flags().StringP("status", "s", "", "Status (Planned, Active, Completed)")
	initiativeCreateCmd.Flags().String("target-date", "", "Target date (YYYY-MM-DD)")
	initiativeCreateCmd.Flags().String("color", "", "Initiative color (hex)")
	initiativeCreateCmd.Flags().String("icon", "", "Initiative icon")
	_ = initiativeCreateCmd.MarkFlagRequired("name")

	// Update flags
	initiativeUpdateCmd.Flags().String("name", "", "New name")
	initiativeUpdateCmd.Flags().StringP("description", "d", "", "New description")
	initiativeUpdateCmd.Flags().String("description-file", "", "Read description from a markdown file (use - for stdin)")
	initiativeUpdateCmd.Flags().StringP("status", "s", "", "New status (Planned, Active, Completed)")
	initiativeUpdateCmd.Flags().String("target-date", "", "New target date (YYYY-MM-DD, or empty to remove)")
	initiativeUpdateCmd.Flags().String("color", "", "New color (hex)")
	initiativeUpdateCmd.Flags().String("icon", "", "New icon")
}
