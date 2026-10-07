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

// fetchInitiativeRelations returns up to limit initiative relations, optionally
// filtered client-side by initiative ID. Because the Linear API exposes no
// filter argument for initiativeRelations, matching is done locally; results
// are gathered across bounded pages so matches beyond the first page are found.
func fetchInitiativeRelations(ctx context.Context, client *api.Client, limit int, initiativeFilter string) ([]api.InitiativeRelation, error) {
	// Pages are capped at a safe size and walked with a cursor; the Linear API
	// returns large relation pages slowly.
	const maxPageSize = 50

	pageSize := limit
	if initiativeFilter != "" {
		pageSize = maxPageSize
	} else if pageSize < 1 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	results := []api.InitiativeRelation{}
	after := ""
	for pages := 0; pages < 3; pages++ {
		page, err := client.GetInitiativeRelations(ctx, pageSize, after, "createdAt", false)
		if err != nil {
			return nil, err
		}
		for _, r := range page.Nodes {
			if initiativeFilter == "" ||
				(r.Initiative != nil && r.Initiative.ID == initiativeFilter) ||
				(r.RelatedInitiative != nil && r.RelatedInitiative.ID == initiativeFilter) {
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

// initiativeRelationCmd is the parent command: initiative relation
var initiativeRelationCmd = &cobra.Command{
	Use:   "relation",
	Short: "Manage initiative relations",
	Long: `Manage Linear initiative parent-child relations.

Examples:
  linear-cli initiative relation list                            # List all initiative relations
  linear-cli initiative relation list --initiative INITIATIVE-ID # Filter by initiative`,
}

var initiativeRelationListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List initiative relations",
	Long: `List initiative parent-child relations, optionally filtered by initiative.

The Linear API does not support server-side filtering of relations, so
--initiative is applied client-side over the fetched page.`,
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
		initiativeFilter, _ := cmd.Flags().GetString("initiative")

		nodes, err := fetchInitiativeRelations(context.Background(), client, limit, initiativeFilter)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list initiative relations: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(nodes)
			return
		}

		if len(nodes) == 0 {
			if plaintext {
				fmt.Println("No initiative relations found")
			} else {
				fmt.Printf("\n%s No initiative relations found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Initiative Relations")
			fmt.Println("Initiative\tRelated\tType")
			for _, r := range nodes {
				initiative := ""
				if r.Initiative != nil {
					initiative = r.Initiative.Name
				}
				related := ""
				if r.RelatedInitiative != nil {
					related = r.RelatedInitiative.Name
				}
				fmt.Printf("%s\t%s\t%s\n", initiative, related, "parent-child")
			}
		} else {
			headers := []string{"Initiative", "Related", "Type"}
			rows := [][]string{}

			for _, r := range nodes {
				initiative := color.New(color.FgWhite, color.Faint).Sprint("unknown")
				if r.Initiative != nil {
					initiative = r.Initiative.Name
				}
				related := color.New(color.FgWhite, color.Faint).Sprint("unknown")
				if r.RelatedInitiative != nil {
					related = r.RelatedInitiative.Name
				}
				rows = append(rows, []string{
					initiative,
					related,
					color.New(color.FgCyan).Sprint("parent-child"),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d initiative relations\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(nodes))
		}
	},
}

func init() {
	initiativeCmd.AddCommand(initiativeRelationCmd)
	initiativeRelationCmd.AddCommand(initiativeRelationListCmd)

	initiativeRelationListCmd.Flags().IntP("limit", "l", 25, "Maximum number of initiative relations to return")
	initiativeRelationListCmd.Flags().String("initiative", "", "Filter by initiative ID")
}
