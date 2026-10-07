package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fetchProjectRelations returns up to limit project relations, optionally
// filtered client-side by project ID. Because the Linear API exposes no filter
// argument for projectRelations, matching is done locally; results are gathered
// across bounded pages so matches beyond the first page are still found.
func fetchProjectRelations(ctx context.Context, client *api.Client, limit int, projectFilter string) ([]api.ProjectRelation, error) {
	// Linear computes projectRelations slowly for large pages; first=100 times
	// out, so pages are capped at a safe size and walked with a cursor.
	const maxPageSize = 50

	pageSize := limit
	if projectFilter != "" {
		pageSize = maxPageSize
	} else if pageSize < 1 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	results := []api.ProjectRelation{}
	after := ""
	for pages := 0; pages < 3; pages++ {
		page, err := client.GetProjectRelations(ctx, pageSize, after, "createdAt", false)
		if err != nil {
			return nil, err
		}
		for _, r := range page.Nodes {
			if projectFilter == "" ||
				(r.Project != nil && r.Project.ID == projectFilter) ||
				(r.RelatedProject != nil && r.RelatedProject.ID == projectFilter) {
				results = append(results, r)
				if limit > 0 && len(results) >= limit {
					return results, nil
				}
			}
		}
		if !page.PageInfo.HasNextPage || page.PageInfo.EndCursor == "" {
			break
		}
		after = page.PageInfo.EndCursor
	}

	return results, nil
}

// projectRelationCmd is the parent command: project relation
var projectRelationCmd = &cobra.Command{
	Use:   "relation",
	Short: "Manage project relations",
	Long: `Manage Linear project dependency relations.

Examples:
  linear-cli project relation list                        # List all project relations
  linear-cli project relation list --project PROJECT-ID   # Filter by project`,
}

var projectRelationListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List project relations",
	Long: `List project dependency relations, optionally filtered by project.

The Linear API does not support server-side filtering of relations, so
--project is applied client-side over the fetched page.`,
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
		projectFilter, _ := cmd.Flags().GetString("project")

		nodes, err := fetchProjectRelations(context.Background(), client, limit, projectFilter)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list project relations: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(nodes)
			return
		}

		if len(nodes) == 0 {
			if plaintext {
				fmt.Println("No project relations found")
			} else {
				fmt.Printf("\n%s No project relations found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Project Relations")
			fmt.Println("Project\tRelated Project\tType")
			for _, r := range nodes {
				project := ""
				if r.Project != nil {
					project = r.Project.Name
				}
				related := ""
				if r.RelatedProject != nil {
					related = r.RelatedProject.Name
				}
				fmt.Printf("%s\t%s\t%s\n", project, related, r.Type)
			}
		} else {
			headers := []string{"Project", "Related Project", "Type"}
			rows := [][]string{}

			for _, r := range nodes {
				project := color.New(color.FgWhite, color.Faint).Sprint("unknown")
				if r.Project != nil {
					project = r.Project.Name
				}
				related := color.New(color.FgWhite, color.Faint).Sprint("unknown")
				if r.RelatedProject != nil {
					related = r.RelatedProject.Name
				}
				rows = append(rows, []string{
					project,
					related,
					color.New(color.FgCyan).Sprint(r.Type),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d project relations\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(nodes))
		}
	},
}

func init() {
	projectCmd.AddCommand(projectRelationCmd)
	projectRelationCmd.AddCommand(projectRelationListCmd)

	projectRelationListCmd.Flags().IntP("limit", "l", 25, "Maximum number of project relations to return")
	projectRelationListCmd.Flags().String("project", "", "Filter by project ID")
}
