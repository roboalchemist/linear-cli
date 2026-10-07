package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var documentUnarchiveCmd = &cobra.Command{
	Use:   "unarchive DOC-ID",
	Short: "Restore an archived document",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.UnarchiveDocument(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to unarchive document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Unarchived document %s", args[0]), plaintext, jsonOut)
	},
}

// ---------------------------------------------------------------------------
// Cycles
// ---------------------------------------------------------------------------

var cycleStartTodayCmd = &cobra.Command{
	Use:   "start-today CYCLE-ID",
	Short: "Start the upcoming cycle as of today",
	Long:  `Start the upcoming cycle immediately, as of midnight today.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.StartUpcomingCycleToday(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to start cycle: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Started cycle %s today", args[0]), plaintext, jsonOut)
	},
}

var cycleShiftAllCmd = &cobra.Command{
	Use:   "shift-all CYCLE-ID --days N",
	Short: "Shift a cycle and all following cycles",
	Long:  `Shift the given cycle and all cycles after it by a number of days.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		days, _ := cmd.Flags().GetFloat64("days")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		client := api.NewClient(authHeader)
		if err := client.ShiftAllCycles(context.Background(), args[0], days); err != nil {
			output.Error(fmt.Sprintf("Failed to shift cycles: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		output.Success(fmt.Sprintf("Shifted %s and following cycles by %.0f day(s)", args[0], days), plaintext, jsonOut)
	},
}

func init() {
	documentCmd.AddCommand(documentUnarchiveCmd)

	cycleCmd.AddCommand(cycleStartTodayCmd)
	cycleCmd.AddCommand(cycleShiftAllCmd)
	cycleShiftAllCmd.Flags().Float64("days", 0, "Number of days to shift (required)")
	_ = cycleShiftAllCmd.MarkFlagRequired("days")
}
