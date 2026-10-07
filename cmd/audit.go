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

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Access Linear audit logs",
	Long: `Access Linear audit logs.

Note: reading audit entries requires workspace admin privileges.

Examples:
  linear-cli audit log list                     # List recent audit entries
  linear-cli audit log list --actor USER-ID     # Filter by actor
  linear-cli audit log list -l 10 -j            # JSON output`,
}

var auditLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage audit log entries",
	Long:  `Manage audit log entries.`,
}

var auditLogListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List audit log entries",
	Long:    `List audit log entries for the workspace (requires admin).`,
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
		actor, _ := cmd.Flags().GetString("actor")

		filter := map[string]interface{}{}
		if actor != "" {
			filter["actor"] = map[string]interface{}{
				"id": map[string]interface{}{"eq": actor},
			}
		}
		if len(filter) == 0 {
			filter = nil
		}

		entries, err := client.GetAuditEntries(context.Background(), filter, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list audit entries: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(entries.Nodes)
			return
		}

		if len(entries.Nodes) == 0 {
			if plaintext {
				fmt.Println("No audit entries found")
			} else {
				fmt.Printf("\n%s No audit entries found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Type", "Actor", "Created", "IP"}
		rows := [][]string{}
		for _, e := range entries.Nodes {
			actorName := ""
			if e.Actor != nil {
				actorName = e.Actor.Name
			} else if e.ActorId != nil {
				actorName = *e.ActorId
			}
			ip := ""
			if e.IP != nil {
				ip = *e.IP
			}
			rows = append(rows, []string{
				e.Type,
				actorName,
				formatDateShort(e.CreatedAt.Format(time.RFC3339)),
				ip,
			})
		}

		if plaintext {
			fmt.Println("# Audit Log")
			fmt.Println(strings.Join(headers, "\t"))
			for _, row := range rows {
				fmt.Println(strings.Join(row, "\t"))
			}
			return
		}

		output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(auditLogCmd)
	auditLogCmd.AddCommand(auditLogListCmd)

	auditLogListCmd.Flags().IntP("limit", "l", 50, "Maximum number of audit entries to return")
	auditLogListCmd.Flags().String("actor", "", "Filter by actor user ID")
}
