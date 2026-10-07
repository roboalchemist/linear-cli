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

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage Linear time schedules",
	Long: `Manage Linear time schedules (used by triage responsibilities).

Examples:
  linear-cli schedule list             # List all time schedules
  linear-cli schedule get SCHEDULE-ID  # Get time schedule details`,
}

var scheduleListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List time schedules",
	Long:    `List all time schedules in the workspace.`,
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

		schedules, err := client.GetTimeSchedules(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list time schedules: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(schedules.Nodes)
			return
		}

		if len(schedules.Nodes) == 0 {
			if plaintext {
				fmt.Println("No time schedules found")
			} else {
				fmt.Printf("\n%s No time schedules found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Timezone", "Created"}
		rows := [][]string{}
		for _, s := range schedules.Nodes {
			// The Linear API does not expose a timezone on TimeSchedule, so the
			// column is retained for interface parity but rendered as "-".
			rows = append(rows, []string{
				s.Name,
				"-",
				formatDateShort(s.CreatedAt.Format(time.RFC3339)),
			})
		}

		if plaintext {
			fmt.Println("# Time Schedules")
			fmt.Println(strings.Join(headers, "\t"))
			for _, r := range rows {
				fmt.Println(strings.Join(r, "\t"))
			}
			return
		}

		output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
	},
}

var scheduleGetCmd = &cobra.Command{
	Use:     "get SCHEDULE-ID",
	Aliases: []string{"show"},
	Short:   "Get time schedule details",
	Long:    `Get details for a specific time schedule.`,
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
		schedule, err := client.GetTimeSchedule(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get time schedule: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(schedule)
			return
		}

		if plaintext {
			fmt.Printf("# Time Schedule %s\n", schedule.ID)
			fmt.Printf("Name: %s\n", schedule.Name)
			if schedule.ExternalId != nil && *schedule.ExternalId != "" {
				fmt.Printf("External ID: %s\n", *schedule.ExternalId)
			}
			if schedule.ExternalUrl != nil && *schedule.ExternalUrl != "" {
				fmt.Printf("External URL: %s\n", *schedule.ExternalUrl)
			}
			fmt.Printf("Created: %s\n", formatDateShort(schedule.CreatedAt.Format(time.RFC3339)))
			fmt.Printf("Updated: %s\n", formatDateShort(schedule.UpdatedAt.Format(time.RFC3339)))
			return
		}

		fmt.Printf("\n%s Time Schedule %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("🕒"),
			color.New(color.FgWhite, color.Bold).Sprint(schedule.ID))
		fmt.Printf("   Name: %s\n", schedule.Name)
		if schedule.ExternalId != nil && *schedule.ExternalId != "" {
			fmt.Printf("   External ID: %s\n", *schedule.ExternalId)
		}
		if schedule.ExternalUrl != nil && *schedule.ExternalUrl != "" {
			fmt.Printf("   External URL: %s\n", *schedule.ExternalUrl)
		}
		fmt.Printf("   Created: %s | Updated: %s\n",
			formatDateShort(schedule.CreatedAt.Format(time.RFC3339)),
			formatDateShort(schedule.UpdatedAt.Format(time.RFC3339)))
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(scheduleListCmd)
	scheduleCmd.AddCommand(scheduleGetCmd)

	scheduleListCmd.Flags().IntP("limit", "l", 50, "Maximum number of time schedules to return")
}
