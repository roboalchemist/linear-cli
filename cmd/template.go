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

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateGetCmd)
	templateCmd.AddCommand(templateSearchCmd)

	templateListCmd.Flags().IntP("limit", "l", 50, "Maximum number of templates to return")
	templateSearchCmd.Flags().IntP("limit", "l", 50, "Maximum number of templates to return")
}
