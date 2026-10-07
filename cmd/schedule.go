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

var scheduleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a time schedule",
	Long: `Create a new time schedule.

The Linear API requires schedule entries, supplied as a JSON array via
--entries. Each entry has startsAt/endsAt (ISO 8601) and a userId or userEmail.

Examples:
  linear-cli schedule create --name "Primary on-call" \
    --entries '[{"startsAt":"2026-08-01T00:00:00Z","endsAt":"2026-08-08T00:00:00Z","userId":"USER-ID"}]'`,
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
		teamID, _ := cmd.Flags().GetString("team-id")
		externalID, _ := cmd.Flags().GetString("external-id")
		externalURL, _ := cmd.Flags().GetString("external-url")
		entriesRaw, _ := cmd.Flags().GetString("entries")
		configRaw, _ := cmd.Flags().GetString("config")
		inputJSON, _ := cmd.Flags().GetString("input-json")

		if teamID != "" {
			fmt.Fprintln(os.Stderr, "Warning: the Linear time schedule API does not accept a team; --team-id was ignored.")
		}

		input := map[string]interface{}{"name": name}
		if externalID != "" {
			input["externalId"] = externalID
		}
		if externalURL != "" {
			input["externalUrl"] = externalURL
		}
		if entriesRaw != "" {
			entries, err := parseJSONValue(entriesRaw)
			if err != nil {
				output.Error(fmt.Sprintf("Invalid --entries: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["entries"] = entries
		}
		if configRaw != "" {
			config, err := parseJSONValue(configRaw)
			if err != nil {
				output.Error(fmt.Sprintf("Invalid --config: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["config"] = config
		}
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if _, ok := input["entries"]; !ok {
			output.Error("--entries is required (a JSON array of TimeScheduleEntryInput objects, or provide entries via --input-json).", plaintext, jsonOut)
			os.Exit(1)
		}

		schedule, err := client.CreateTimeSchedule(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create time schedule: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(schedule)
		} else {
			output.Success(fmt.Sprintf("Created time schedule %s",
				color.New(color.FgWhite, color.Bold).Sprint(schedule.Name)), plaintext, jsonOut)
		}
	},
}

var scheduleUpdateCmd = &cobra.Command{
	Use:     "update SCHEDULE-ID",
	Aliases: []string{"edit"},
	Short:   "Update a time schedule",
	Long: `Update a time schedule's name, entries, or external references.

Examples:
  linear-cli schedule update SCHEDULE-ID --name "Renamed schedule"
  linear-cli schedule update SCHEDULE-ID --entries '[{"startsAt":"2026-08-01T00:00:00Z","endsAt":"2026-08-08T00:00:00Z","userId":"USER-ID"}]'`,
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
		if cmd.Flags().Changed("external-id") {
			v, _ := cmd.Flags().GetString("external-id")
			input["externalId"] = v
		}
		if cmd.Flags().Changed("external-url") {
			v, _ := cmd.Flags().GetString("external-url")
			input["externalUrl"] = v
		}
		if cmd.Flags().Changed("entries") {
			raw, _ := cmd.Flags().GetString("entries")
			v, err := parseJSONValue(raw)
			if err != nil {
				output.Error(fmt.Sprintf("Invalid --entries: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["entries"] = v
		}
		if cmd.Flags().Changed("config") {
			raw, _ := cmd.Flags().GetString("config")
			v, err := parseJSONValue(raw)
			if err != nil {
				output.Error(fmt.Sprintf("Invalid --config: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["config"] = v
		}
		inputJSON, _ := cmd.Flags().GetString("input-json")
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --entries, --external-id, --external-url, --config, or --input-json.", plaintext, jsonOut)
			os.Exit(1)
		}

		schedule, err := client.UpdateTimeSchedule(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update time schedule: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(schedule)
		} else {
			output.Success(fmt.Sprintf("Updated time schedule %s",
				color.New(color.FgWhite, color.Bold).Sprint(schedule.Name)), plaintext, jsonOut)
		}
	},
}

var scheduleDeleteCmd = &cobra.Command{
	Use:     "delete SCHEDULE-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a time schedule",
	Long:    `Delete a time schedule.`,
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
		if err := client.DeleteTimeSchedule(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete time schedule: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted time schedule", plaintext, jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.AddCommand(scheduleListCmd)
	scheduleCmd.AddCommand(scheduleGetCmd)
	scheduleCmd.AddCommand(scheduleCreateCmd)
	scheduleCmd.AddCommand(scheduleUpdateCmd)
	scheduleCmd.AddCommand(scheduleDeleteCmd)

	scheduleListCmd.Flags().IntP("limit", "l", 50, "Maximum number of time schedules to return")

	// Create flags
	scheduleCreateCmd.Flags().StringP("name", "n", "", "Time schedule name (required)")
	scheduleCreateCmd.Flags().String("team-id", "", "Team ID (accepted for compatibility; not sent to the API)")
	scheduleCreateCmd.Flags().String("external-id", "", "External schedule identifier")
	scheduleCreateCmd.Flags().String("external-url", "", "URL to the external schedule")
	scheduleCreateCmd.Flags().String("entries", "", "Schedule entries as a JSON array (required)")
	scheduleCreateCmd.Flags().String("config", "", "[ALPHA] Schedule configuration as JSON")
	scheduleCreateCmd.Flags().String("input-json", "", "Additional TimeScheduleCreateInput fields as a JSON object")
	_ = scheduleCreateCmd.MarkFlagRequired("name")

	// Update flags
	scheduleUpdateCmd.Flags().StringP("name", "n", "", "New time schedule name")
	scheduleUpdateCmd.Flags().String("external-id", "", "New external schedule identifier")
	scheduleUpdateCmd.Flags().String("external-url", "", "New external schedule URL")
	scheduleUpdateCmd.Flags().String("entries", "", "New schedule entries as a JSON array")
	scheduleUpdateCmd.Flags().String("config", "", "[ALPHA] New schedule configuration as JSON")
	scheduleUpdateCmd.Flags().String("input-json", "", "Additional TimeScheduleUpdateInput fields as a JSON object")
}
