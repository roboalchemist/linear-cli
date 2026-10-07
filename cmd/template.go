package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage Linear templates",
	Long: `Manage Linear templates (issue, project, document, and more).

Examples:
  linear-cli template list                # List all templates
  linear-cli template get TEMPLATE-ID     # Get template details
  linear-cli template search "bug"        # Search templates by name`,
}

// templateRows builds the shared Name/Type/Team/Created table rows.
func templateRows(templates []api.AdminTemplate) [][]string {
	rows := [][]string{}
	for _, t := range templates {
		team := ""
		if t.Team != nil {
			team = t.Team.Key
		}
		rows = append(rows, []string{
			t.Name,
			t.Type,
			team,
			formatDateShort(t.CreatedAt.Format(time.RFC3339)),
		})
	}
	return rows
}

// renderTemplateList renders a list of templates in the requested format.
func renderTemplateList(headers []string, rows [][]string, plaintext, jsonOut bool) {
	if plaintext {
		fmt.Println(strings.Join(headers, "\t"))
		for _, row := range rows {
			fmt.Println(strings.Join(row, "\t"))
		}
		return
	}
	output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
}

var templateListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List templates",
	Long:    `List all templates in the workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")

		templates, err := client.GetTemplates(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list templates: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if limit > 0 && len(templates) > limit {
			templates = templates[:limit]
		}

		if jsonOut {
			output.JSON(templates)
			return
		}

		if len(templates) == 0 {
			if plaintext {
				fmt.Println("No templates found")
			} else {
				fmt.Printf("\n%s No templates found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Type", "Team", "Created"}
		if plaintext {
			fmt.Println("# Templates")
		}
		renderTemplateList(headers, templateRows(templates), plaintext, jsonOut)
	},
}

var templateGetCmd = &cobra.Command{
	Use:     "get TEMPLATE-ID",
	Aliases: []string{"show"},
	Short:   "Get template details",
	Long:    `Get details for a specific template.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		template, err := client.GetTemplate(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get template: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(template)
			return
		}

		team := ""
		if template.Team != nil {
			team = template.Team.Key
		}
		description := ""
		if template.Description != nil {
			description = *template.Description
		}

		if plaintext {
			fmt.Printf("# Template %s\n", template.ID)
			fmt.Printf("Name: %s\n", template.Name)
			fmt.Printf("Type: %s\n", template.Type)
			fmt.Printf("Team: %s\n", team)
			if description != "" {
				fmt.Printf("Description: %s\n", description)
			}
			fmt.Printf("Created: %s\n", formatDateShort(template.CreatedAt.Format(time.RFC3339)))
			fmt.Printf("Updated: %s\n", formatDateShort(template.UpdatedAt.Format(time.RFC3339)))
			return
		}

		fmt.Printf("\n%s Template %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("📄"),
			color.New(color.FgWhite, color.Bold).Sprint(template.ID))
		fmt.Printf("   Name: %s\n", template.Name)
		fmt.Printf("   Type: %s\n", template.Type)
		fmt.Printf("   Team: %s\n", team)
		if description != "" {
			fmt.Printf("   Description: %s\n", description)
		}
		fmt.Printf("   Created: %s | Updated: %s\n",
			formatDateShort(template.CreatedAt.Format(time.RFC3339)),
			formatDateShort(template.UpdatedAt.Format(time.RFC3339)))
	},
}

var templateSearchCmd = &cobra.Command{
	Use:     "search TERM",
	Aliases: []string{"find"},
	Short:   "Search templates by name",
	Long:    `Search templates whose name matches TERM (case-insensitive).`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		limit, _ := cmd.Flags().GetInt("limit")

		templates, err := client.SearchTemplates(context.Background(), args[0], limit)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to search templates: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(templates)
			return
		}

		if len(templates) == 0 {
			if plaintext {
				fmt.Println("No templates found")
			} else {
				fmt.Printf("\n%s No templates found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Type", "Team", "Created"}
		if plaintext {
			fmt.Println("# Templates")
		}
		renderTemplateList(headers, templateRows(templates), plaintext, jsonOut)
	},
}

var templateCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a template",
	Long: `Create a new template.

--type must be one of: issue, project. --template-data is a JSON object of
pre-filled attributes for the target entity type and defaults to {}.

Examples:
  linear-cli template create --name "Bug report" --type issue --team-id TEAM-ID
  linear-cli template create --name "Feature" --type project \
    --template-data '{"description":"Describe the feature"}'`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		name, _ := cmd.Flags().GetString("name")
		templateType, _ := cmd.Flags().GetString("type")
		teamID, _ := cmd.Flags().GetString("team-id")
		description, _ := cmd.Flags().GetString("description")
		templateDataRaw, _ := cmd.Flags().GetString("template-data")
		inputJSON, _ := cmd.Flags().GetString("input-json")

		if templateType != "issue" && templateType != "project" {
			output.Error("--type must be one of: issue, project.", plaintext, jsonOut)
			os.Exit(1)
		}

		var templateData interface{} = map[string]interface{}{}
		if templateDataRaw != "" {
			templateData, err = parseJSONValue(templateDataRaw)
			if err != nil {
				output.Error(fmt.Sprintf("Invalid --template-data: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		input := map[string]interface{}{
			"name":         name,
			"type":         templateType,
			"templateData": templateData,
		}
		if teamID != "" {
			input["teamId"] = teamID
		}
		if description != "" {
			input["description"] = description
		}
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		template, err := client.CreateTemplate(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create template: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(template)
		} else {
			output.Success(fmt.Sprintf("Created template %s (%s)",
				color.New(color.FgWhite, color.Bold).Sprint(template.Name), template.Type), plaintext, jsonOut)
		}
	},
}

var templateUpdateCmd = &cobra.Command{
	Use:     "update TEMPLATE-ID",
	Aliases: []string{"edit"},
	Short:   "Update a template",
	Long: `Update a template's name, description, team, or template data.

Examples:
  linear-cli template update TEMPLATE-ID --name "Renamed template"
  linear-cli template update TEMPLATE-ID --template-data '{"title":"Updated"}'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		input := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			input["name"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			input["description"] = v
		}
		if cmd.Flags().Changed("team-id") {
			v, _ := cmd.Flags().GetString("team-id")
			input["teamId"] = v
		}
		if cmd.Flags().Changed("color") {
			v, _ := cmd.Flags().GetString("color")
			input["color"] = v
		}
		if cmd.Flags().Changed("icon") {
			v, _ := cmd.Flags().GetString("icon")
			input["icon"] = v
		}
		if cmd.Flags().Changed("sort-order") {
			v, _ := cmd.Flags().GetFloat64("sort-order")
			input["sortOrder"] = v
		}
		if cmd.Flags().Changed("template-data") {
			raw, _ := cmd.Flags().GetString("template-data")
			v, err := parseJSONValue(raw)
			if err != nil {
				output.Error(fmt.Sprintf("Invalid --template-data: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["templateData"] = v
		}
		inputJSON, _ := cmd.Flags().GetString("input-json")
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --description, --team-id, --template-data, --color, --icon, --sort-order, or --input-json.", plaintext, jsonOut)
			os.Exit(1)
		}

		template, err := client.UpdateTemplate(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update template: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(template)
		} else {
			output.Success(fmt.Sprintf("Updated template %s",
				color.New(color.FgWhite, color.Bold).Sprint(template.Name)), plaintext, jsonOut)
		}
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:     "delete TEMPLATE-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a template",
	Long:    `Delete a template.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		if err := client.DeleteTemplate(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete template: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted template", plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateGetCmd)
	templateCmd.AddCommand(templateSearchCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateUpdateCmd)
	templateCmd.AddCommand(templateDeleteCmd)

	templateListCmd.Flags().IntP("limit", "l", 50, "Maximum number of templates to return")
	templateSearchCmd.Flags().IntP("limit", "l", 50, "Maximum number of templates to return")

	// Create flags
	templateCreateCmd.Flags().StringP("name", "n", "", "Template name (required)")
	templateCreateCmd.Flags().String("type", "", "Template type: issue or project (required)")
	templateCreateCmd.Flags().String("team-id", "", "Team ID (omit for a workspace-wide template)")
	templateCreateCmd.Flags().StringP("description", "d", "", "Template description")
	templateCreateCmd.Flags().String("template-data", "", "Template data as JSON (defaults to {})")
	templateCreateCmd.Flags().String("input-json", "", "Additional TemplateCreateInput fields as a JSON object")
	_ = templateCreateCmd.MarkFlagRequired("name")
	_ = templateCreateCmd.MarkFlagRequired("type")

	// Update flags
	templateUpdateCmd.Flags().StringP("name", "n", "", "New template name")
	templateUpdateCmd.Flags().StringP("description", "d", "", "New description")
	templateUpdateCmd.Flags().String("team-id", "", "New team ID")
	templateUpdateCmd.Flags().String("template-data", "", "New template data as JSON")
	templateUpdateCmd.Flags().String("color", "", "New template icon color")
	templateUpdateCmd.Flags().String("icon", "", "New template icon")
	templateUpdateCmd.Flags().Float64("sort-order", 0, "New sort order")
	templateUpdateCmd.Flags().String("input-json", "", "Additional TemplateUpdateInput fields as a JSON object")
}
