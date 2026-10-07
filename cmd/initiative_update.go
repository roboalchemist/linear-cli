package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initiativeUpdateListCmd lists initiative status updates. It hangs off the
// existing `initiative update` command group, so it is invoked as
// `initiative update list`.
var initiativeUpdateListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List initiative status updates",
	Long: `List initiative status updates, optionally scoped to a single initiative.

Examples:
  linear-cli initiative update list
  linear-cli initiative update list --initiative INITIATIVE-ID
  linear-cli initiative update list --limit 10`,
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
		initiativeID, _ := cmd.Flags().GetString("initiative")

		var filter map[string]interface{}
		if initiativeID != "" {
			filter = map[string]interface{}{
				"initiative": map[string]interface{}{
					"id": map[string]interface{}{"eq": initiativeID},
				},
			}
		}

		updates, err := client.GetInitiativeUpdates(context.Background(), filter, limit, "", "createdAt")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list initiative updates: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(updates.Nodes)
			return
		}

		if len(updates.Nodes) == 0 {
			output.Info("No initiative updates found", plaintext, jsonOut)
			return
		}

		headers := []string{"Initiative", "Health", "Body", "Created"}
		rows := make([][]string, len(updates.Nodes))
		for i, u := range updates.Nodes {
			initiative := "-"
			if u.Initiative != nil {
				initiative = u.Initiative.Name
			}

			health := u.Health
			switch u.Health {
			case "onTrack":
				health = color.New(color.FgGreen).Sprint(u.Health)
			case "atRisk":
				health = color.New(color.FgYellow).Sprint(u.Health)
			case "offTrack":
				health = color.New(color.FgRed).Sprint(u.Health)
			}

			body := strings.ReplaceAll(u.Body, "\n", " ")
			if body == "" {
				body = "-"
			}

			rows[i] = []string{
				truncateString(initiative, 30),
				health,
				truncateString(body, 60),
				u.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		if updates.PageInfo.HasNextPage && !plaintext {
			fmt.Printf("\n%s Use --limit to see more results\n",
				color.New(color.FgYellow).Sprint("ℹ️"))
		}
	},
}

func init() {
	initiativeUpdateCmd.AddCommand(initiativeUpdateListCmd)

	initiativeUpdateListCmd.Flags().String("initiative", "", "Filter by initiative ID")
	initiativeUpdateListCmd.Flags().IntP("limit", "l", 25, "Maximum number of updates to return")
}
