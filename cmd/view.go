package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dorkitude/linear-cli/pkg/api"
	"github.com/dorkitude/linear-cli/pkg/auth"
	"github.com/dorkitude/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Manage Linear custom views (saved filters)",
	Long: `Create, list, run, and manage Linear custom views (saved filters).

Custom views are saved filter configurations that can be shared across teams.
The key feature is 'view run', which executes a view and returns matching issues or projects.

Examples:
  linear-cli view list
  linear-cli view get VIEW-ID
  linear-cli view run VIEW-ID
  linear-cli view create --name "My Bugs" --model issue --filter-json '{"state":{"type":{"eq":"started"}}}'
  linear-cli view update VIEW-ID --name "Renamed View"
  linear-cli view delete VIEW-ID`,
}

var viewListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List custom views",
	Long:    `List all custom views in your Linear workspace.`,
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
			limit = 50
		}

		filter := make(map[string]interface{})

		if modelName, _ := cmd.Flags().GetString("model"); modelName != "" {
			filter["modelName"] = map[string]interface{}{"eq": modelName}
		}

		if shared, _ := cmd.Flags().GetBool("shared"); cmd.Flags().Changed("shared") && shared {
			filter["shared"] = map[string]interface{}{"eq": true}
		}

		if teamKey, _ := cmd.Flags().GetString("team"); teamKey != "" {
			team, err := client.GetTeam(context.Background(), teamKey)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
				os.Exit(1)
			}
			filter["team"] = map[string]interface{}{"id": map[string]interface{}{"eq": team.ID}}
		}

		var filterArg map[string]interface{}
		if len(filter) > 0 {
			filterArg = filter
		}

		views, err := client.GetCustomViews(context.Background(), filterArg, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list views: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(views.Nodes) == 0 {
			output.Info("No custom views found", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(views.Nodes)
			return
		}

		if plaintext {
			fmt.Println("# Custom Views")
			for _, v := range views.Nodes {
				fmt.Printf("## %s\n", v.Name)
				fmt.Printf("- **ID**: %s\n", v.ID)
				fmt.Printf("- **Model**: %s\n", v.ModelName)
				fmt.Printf("- **Shared**: %v\n", v.Shared)
				if v.Creator != nil {
					fmt.Printf("- **Creator**: %s\n", v.Creator.Name)
				}
				if v.Team != nil {
					fmt.Printf("- **Team**: %s\n", v.Team.Key)
				}
				fmt.Printf("- **Updated**: %s\n", v.UpdatedAt.Format("2006-01-02"))
				if v.Description != nil && *v.Description != "" {
					fmt.Printf("- **Description**: %s\n", *v.Description)
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d views\n", len(views.Nodes))
			return
		}

		headers := []string{"Name", "Model", "Shared", "Creator", "Team", "Updated", "ID"}
		rows := make([][]string, len(views.Nodes))

		for i, v := range views.Nodes {
			creator := ""
			if v.Creator != nil {
				creator = v.Creator.Name
			}

			team := ""
			if v.Team != nil {
				team = v.Team.Key
			}

			shared := ""
			if v.Shared {
				shared = color.New(color.FgGreen).Sprint("yes")
			} else {
				shared = color.New(color.FgWhite, color.Faint).Sprint("no")
			}

			modelColor := color.New(color.FgCyan)
			if v.ModelName == "project" {
				modelColor = color.New(color.FgMagenta)
			}

			rows[i] = []string{
				truncateString(v.Name, 30),
				modelColor.Sprint(v.ModelName),
				shared,
				creator,
				team,
				v.UpdatedAt.Format("2006-01-02"),
				v.ID,
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, false, false)

		fmt.Printf("\n%s %d views\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(views.Nodes))

		if views.PageInfo.HasNextPage {
			fmt.Printf("%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("i"))
		}
	},
}

var viewGetCmd = &cobra.Command{
	Use:     "get [view-id]",
	Aliases: []string{"show"},
	Short:   "Get view details",
	Long:    `Get detailed information about a custom view including its filter configuration.`,
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
		view, err := client.GetCustomView(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch view: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(view)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n\n", view.Name)
			fmt.Printf("## Details\n")
			fmt.Printf("- **ID**: %s\n", view.ID)
			fmt.Printf("- **Slug**: %s\n", view.SlugId)
			fmt.Printf("- **Model**: %s\n", view.ModelName)
			fmt.Printf("- **Shared**: %v\n", view.Shared)
			if view.Description != nil && *view.Description != "" {
				fmt.Printf("- **Description**: %s\n", *view.Description)
			}
			if view.Icon != nil && *view.Icon != "" {
				fmt.Printf("- **Icon**: %s\n", *view.Icon)
			}
			if view.Color != nil && *view.Color != "" {
				fmt.Printf("- **Color**: %s\n", *view.Color)
			}
			if view.Creator != nil {
				fmt.Printf("- **Creator**: %s (%s)\n", view.Creator.Name, view.Creator.Email)
			}
			if view.Owner != nil {
				fmt.Printf("- **Owner**: %s (%s)\n", view.Owner.Name, view.Owner.Email)
			}
			if view.Team != nil {
				fmt.Printf("- **Team**: %s (%s)\n", view.Team.Name, view.Team.Key)
			}
			fmt.Printf("- **Created**: %s\n", view.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", view.UpdatedAt.Format("2006-01-02 15:04:05"))

			if len(view.FilterData) > 0 {
				filterJSON, _ := json.MarshalIndent(view.FilterData, "", "  ")
				fmt.Printf("\n## Filter Data\n```json\n%s\n```\n", string(filterJSON))
			}
			if len(view.ProjectFilterData) > 0 {
				filterJSON, _ := json.MarshalIndent(view.ProjectFilterData, "", "  ")
				fmt.Printf("\n## Project Filter Data\n```json\n%s\n```\n", string(filterJSON))
			}
			return
		}

		// Rich display
		fmt.Println()
		nameStr := view.Name
		if view.Icon != nil && *view.Icon != "" {
			nameStr = *view.Icon + " " + nameStr
		}
		fmt.Printf("%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("View:"),
			color.New(color.FgWhite, color.Bold).Sprint(nameStr))
		fmt.Println(strings.Repeat("-", 50))

		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), view.ID)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Model:"), view.ModelName)

		if view.Description != nil && *view.Description != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Description:"), *view.Description)
		}

		shared := "No"
		if view.Shared {
			shared = color.New(color.FgGreen).Sprint("Yes")
		}
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Shared:"), shared)

		if view.Creator != nil {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Creator:"), view.Creator.Name)
		}
		if view.Owner != nil {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Owner:"), view.Owner.Name)
		}
		if view.Team != nil {
			fmt.Printf("%s %s (%s)\n", color.New(color.Bold).Sprint("Team:"), view.Team.Name, view.Team.Key)
		}

		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Created:"), view.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Updated:"), view.UpdatedAt.Format("2006-01-02 15:04:05"))

		if len(view.FilterData) > 0 {
			filterJSON, _ := json.MarshalIndent(view.FilterData, "", "  ")
			fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Filter Data:"), string(filterJSON))
		}
		if len(view.ProjectFilterData) > 0 {
			filterJSON, _ := json.MarshalIndent(view.ProjectFilterData, "", "  ")
			fmt.Printf("\n%s\n%s\n", color.New(color.Bold).Sprint("Project Filter Data:"), string(filterJSON))
		}

		fmt.Println()
	},
}

var viewRunCmd = &cobra.Command{
	Use:     "run [view-id]",
	Aliases: []string{"exec"},
	Short:   "Execute a view and show matching items",
	Long: `Execute a custom view and display the matching issues or projects.

This is the primary feature of custom views — run a saved filter and see results.
The view's modelName determines whether issues or projects are returned.

Examples:
  linear-cli view run VIEW-ID
  linear-cli view run VIEW-ID --limit 100
  linear-cli view run VIEW-ID --json`,
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

		// First get the view to determine its model type
		view, err := client.GetCustomView(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch view: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 50
		}

		switch strings.ToLower(view.ModelName) {
		case "issue":
			issues, err := client.GetCustomViewIssues(context.Background(), view.ID, limit, "")
			if err != nil {
				output.Error(fmt.Sprintf("Failed to run view: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			viewLabel := fmt.Sprintf("issues in view %q", view.Name)
			renderIssueCollection(issues, plaintext, jsonOut, "No issues match this view", viewLabel, fmt.Sprintf("# %s", view.Name))

		case "project":
			projects, err := client.GetCustomViewProjects(context.Background(), view.ID, limit, "")
			if err != nil {
				output.Error(fmt.Sprintf("Failed to run view: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			renderViewProjects(projects, view.Name, plaintext, jsonOut)

		default:
			output.Error(fmt.Sprintf("Unsupported view model type: %s", view.ModelName), plaintext, jsonOut)
			os.Exit(1)
		}
	},
}

func renderViewProjects(projects *api.Projects, viewName string, plaintext, jsonOut bool) {
	if len(projects.Nodes) == 0 {
		output.Info("No projects match this view", plaintext, jsonOut)
		return
	}

	if jsonOut {
		output.JSON(projects.Nodes)
		return
	}

	if plaintext {
		fmt.Printf("# %s\n", viewName)
		for _, project := range projects.Nodes {
			fmt.Printf("## %s\n", project.Name)
			fmt.Printf("- **ID**: %s\n", project.ID)
			fmt.Printf("- **State**: %s\n", project.State)
			fmt.Printf("- **Progress**: %.0f%%\n", project.Progress*100)
			if project.Lead != nil {
				fmt.Printf("- **Lead**: %s\n", project.Lead.Name)
			}
			if project.Teams != nil && len(project.Teams.Nodes) > 0 {
				teams := ""
				for i, team := range project.Teams.Nodes {
					if i > 0 {
						teams += ", "
					}
					teams += team.Key
				}
				fmt.Printf("- **Teams**: %s\n", teams)
			}
			fmt.Printf("- **Created**: %s\n", project.CreatedAt.Format("2006-01-02"))
			fmt.Printf("- **Updated**: %s\n", project.UpdatedAt.Format("2006-01-02"))
			fmt.Printf("- **URL**: %s\n", constructProjectURL(project.ID, project.URL))
			fmt.Println()
		}
		fmt.Printf("\nTotal: %d projects in view %q\n", len(projects.Nodes), viewName)
		return
	}

	headers := []string{"Name", "State", "Lead", "Teams", "Created", "Updated", "URL"}
	rows := [][]string{}

	for _, project := range projects.Nodes {
		lead := color.New(color.FgYellow).Sprint("Unassigned")
		if project.Lead != nil {
			lead = project.Lead.Name
		}

		teams := ""
		if project.Teams != nil && len(project.Teams.Nodes) > 0 {
			for i, team := range project.Teams.Nodes {
				if i > 0 {
					teams += ", "
				}
				teams += team.Key
			}
		}

		stateColor := color.New(color.FgGreen)
		switch project.State {
		case "planned":
			stateColor = color.New(color.FgCyan)
		case "started":
			stateColor = color.New(color.FgBlue)
		case "paused":
			stateColor = color.New(color.FgYellow)
		case "completed":
			stateColor = color.New(color.FgGreen)
		case "canceled":
			stateColor = color.New(color.FgRed)
		}

		rows = append(rows, []string{
			truncateString(project.Name, 25),
			stateColor.Sprint(project.State),
			lead,
			teams,
			project.CreatedAt.Format("2006-01-02"),
			project.UpdatedAt.Format("2006-01-02"),
			constructProjectURL(project.ID, project.URL),
		})
	}

	output.Table(output.TableData{
		Headers: headers,
		Rows:    rows,
	}, plaintext, jsonOut)

	fmt.Printf("\n%s %d projects in view %q\n",
		color.New(color.FgGreen).Sprint("✓"),
		len(projects.Nodes),
		viewName)

	if projects.PageInfo.HasNextPage {
		fmt.Printf("%s Use --limit to see more results\n",
			color.New(color.FgYellow).Sprint("i"))
	}
}

var viewCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new custom view",
	Long: `Create a new custom view with a name and optional filter configuration.

The --filter-json flag accepts raw JSON matching Linear's IssueFilter or ProjectFilter schema.

Examples:
  linear-cli view create --name "My Bugs" --model issue
  linear-cli view create --name "Active Projects" --model project --shared
  linear-cli view create --name "Urgent Issues" --filter-json '{"priority":{"eq":1}}'
  linear-cli view create --name "Team Bugs" --team ENG --filter-json '{"state":{"type":{"eq":"started"}}}'`,
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
		description, _ := cmd.Flags().GetString("description")
		modelName, _ := cmd.Flags().GetString("model")
		teamKey, _ := cmd.Flags().GetString("team")
		shared, _ := cmd.Flags().GetBool("shared")
		filterJSON, _ := cmd.Flags().GetString("filter-json")

		if name == "" {
			output.Error("Name is required (--name)", plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{
			"name":   name,
			"shared": shared,
		}

		if description != "" {
			input["description"] = description
		}

		if teamKey != "" {
			team, err := client.GetTeam(context.Background(), teamKey)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["teamId"] = team.ID
		}

		if filterJSON != "" {
			var filterData map[string]interface{}
			if err := json.Unmarshal([]byte(filterJSON), &filterData); err != nil {
				output.Error(fmt.Sprintf("Invalid filter JSON: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			if modelName == "project" {
				input["projectFilterData"] = filterData
			} else {
				input["filterData"] = filterData
			}
		}

		view, err := client.CreateCustomView(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create view: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(view)
		} else if plaintext {
			fmt.Printf("Created view: %s (ID: %s)\n", view.Name, view.ID)
		} else {
			fmt.Printf("%s Created view %s (ID: %s)\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(view.Name),
				color.New(color.FgWhite, color.Faint).Sprint(view.ID))
		}
	},
}

var viewUpdateCmd = &cobra.Command{
	Use:     "update [view-id]",
	Aliases: []string{"edit"},
	Short:   "Update a custom view",
	Long: `Update an existing custom view's name, description, sharing, or filters.

Examples:
  linear-cli view update VIEW-ID --name "Renamed View"
  linear-cli view update VIEW-ID --shared
  linear-cli view update VIEW-ID --filter-json '{"priority":{"eq":1}}'`,
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

		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			input["description"] = description
		}

		if cmd.Flags().Changed("shared") {
			shared, _ := cmd.Flags().GetBool("shared")
			input["shared"] = shared
		}

		if cmd.Flags().Changed("filter-json") {
			filterJSON, _ := cmd.Flags().GetString("filter-json")
			var filterData map[string]interface{}
			if err := json.Unmarshal([]byte(filterJSON), &filterData); err != nil {
				output.Error(fmt.Sprintf("Invalid filter JSON: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			// Need to check the view's model to decide which filter field to use
			view, err := client.GetCustomView(context.Background(), args[0])
			if err != nil {
				output.Error(fmt.Sprintf("Failed to fetch view: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			if strings.ToLower(view.ModelName) == "project" {
				input["projectFilterData"] = filterData
			} else {
				input["filterData"] = filterData
			}
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use flags to specify what to update.", plaintext, jsonOut)
			os.Exit(1)
		}

		view, err := client.UpdateCustomView(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update view: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(view)
		} else if plaintext {
			fmt.Printf("Updated view: %s\n", view.Name)
			fmt.Printf("ID: %s\n", view.ID)
			fmt.Printf("Model: %s\n", view.ModelName)
			fmt.Printf("Shared: %v\n", view.Shared)
		} else {
			fmt.Printf("%s Updated view %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(view.Name))
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(view.ID))
			fmt.Printf("  Model: %s\n", view.ModelName)
			shared := "No"
			if view.Shared {
				shared = color.New(color.FgGreen).Sprint("Yes")
			}
			fmt.Printf("  Shared: %s\n", shared)
		}
	},
}

var viewDeleteCmd = &cobra.Command{
	Use:     "delete [view-id]",
	Aliases: []string{"rm"},
	Short:   "Delete a custom view",
	Long:    `Delete an existing custom view.`,
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

		err = client.DeleteCustomView(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete view: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "id": args[0], "action": "deleted"})
		} else if plaintext {
			fmt.Printf("Deleted view %s\n", args[0])
		} else {
			fmt.Printf("%s Deleted view %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Faint).Sprint(args[0]))
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.AddCommand(viewListCmd)
	viewCmd.AddCommand(viewGetCmd)
	viewCmd.AddCommand(viewRunCmd)
	viewCmd.AddCommand(viewCreateCmd)
	viewCmd.AddCommand(viewUpdateCmd)
	viewCmd.AddCommand(viewDeleteCmd)

	// List flags
	viewListCmd.Flags().IntP("limit", "l", 50, "Maximum number of views to fetch")
	viewListCmd.Flags().Bool("shared", false, "Show only shared views")
	viewListCmd.Flags().StringP("model", "m", "", "Filter by model type (issue, project)")
	viewListCmd.Flags().StringP("team", "t", "", "Filter by team key")

	// Run flags
	viewRunCmd.Flags().IntP("limit", "l", 50, "Maximum number of results to fetch")

	// Create flags
	viewCreateCmd.Flags().String("name", "", "View name (required)")
	viewCreateCmd.Flags().StringP("description", "d", "", "View description")
	viewCreateCmd.Flags().StringP("model", "m", "issue", "Model type: issue (default), project")
	viewCreateCmd.Flags().StringP("team", "t", "", "Team key")
	viewCreateCmd.Flags().Bool("shared", false, "Make the view shared")
	viewCreateCmd.Flags().String("filter-json", "", "Raw JSON filter (IssueFilter or ProjectFilter schema)")
	_ = viewCreateCmd.MarkFlagRequired("name")

	// Update flags
	viewUpdateCmd.Flags().String("name", "", "New name for the view")
	viewUpdateCmd.Flags().StringP("description", "d", "", "New description")
	viewUpdateCmd.Flags().Bool("shared", false, "Set shared status")
	viewUpdateCmd.Flags().String("filter-json", "", "New raw JSON filter")
}
