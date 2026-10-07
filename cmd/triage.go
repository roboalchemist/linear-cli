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

var triageCmd = &cobra.Command{
	Use:   "triage",
	Short: "Manage Linear triage settings",
	Long: `Manage Linear triage responsibilities.

Examples:
  linear-cli triage responsibility list            # List triage responsibilities
  linear-cli triage responsibility get TR-ID       # Get a triage responsibility`,
}

var triageResponsibilityCmd = &cobra.Command{
	Use:   "responsibility",
	Short: "Manage triage responsibilities",
	Long:  `Manage the teams and users responsible for triage.`,
}

var triageResponsibilityListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List triage responsibilities",
	Long:    `List all triage responsibilities in the workspace.`,
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

		responsibilities, err := client.GetTriageResponsibilities(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list triage responsibilities: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(responsibilities.Nodes)
			return
		}

		if len(responsibilities.Nodes) == 0 {
			if plaintext {
				fmt.Println("No triage responsibilities found")
			} else {
				fmt.Printf("\n%s No triage responsibilities found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Team", "Action", "Created"}
		rows := [][]string{}
		for _, r := range responsibilities.Nodes {
			team := ""
			if r.Team != nil {
				team = r.Team.Key
			}
			rows = append(rows, []string{
				team,
				r.Action,
				formatDateShort(r.CreatedAt.Format(time.RFC3339)),
			})
		}

		if plaintext {
			fmt.Println("# Triage Responsibilities")
			fmt.Println(strings.Join(headers, "\t"))
			for _, row := range rows {
				fmt.Println(strings.Join(row, "\t"))
			}
			return
		}

		output.Table(output.TableData{Headers: headers, Rows: rows}, plaintext, jsonOut)
	},
}

var triageResponsibilityGetCmd = &cobra.Command{
	Use:     "get TR-ID",
	Aliases: []string{"show"},
	Short:   "Get a triage responsibility",
	Long:    `Get details for a specific triage responsibility.`,
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
		tr, err := client.GetTriageResponsibility(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get triage responsibility: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(tr)
			return
		}

		team := ""
		if tr.Team != nil {
			team = tr.Team.Key
		}
		currentUser := ""
		if tr.CurrentUser != nil {
			currentUser = tr.CurrentUser.Name
		}
		schedule := ""
		if tr.TimeSchedule != nil {
			schedule = tr.TimeSchedule.Name
		}

		if plaintext {
			fmt.Printf("# Triage Responsibility %s\n", tr.ID)
			fmt.Printf("Team: %s\n", team)
			fmt.Printf("Action: %s\n", tr.Action)
			if currentUser != "" {
				fmt.Printf("Current User: %s\n", currentUser)
			}
			if schedule != "" {
				fmt.Printf("Time Schedule: %s\n", schedule)
			}
			fmt.Printf("Created: %s\n", formatDateShort(tr.CreatedAt.Format(time.RFC3339)))
			return
		}

		fmt.Printf("\n%s Triage Responsibility %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("🎯"),
			color.New(color.FgWhite, color.Bold).Sprint(tr.ID))
		fmt.Printf("   Team: %s\n", team)
		fmt.Printf("   Action: %s\n", tr.Action)
		if currentUser != "" {
			fmt.Printf("   Current User: %s\n", currentUser)
		}
		if schedule != "" {
			fmt.Printf("   Time Schedule: %s\n", schedule)
		}
		fmt.Printf("   Created: %s\n", formatDateShort(tr.CreatedAt.Format(time.RFC3339)))
	},
}

func init() {
	rootCmd.AddCommand(triageCmd)
	triageCmd.AddCommand(triageResponsibilityCmd)
	triageResponsibilityCmd.AddCommand(triageResponsibilityListCmd)
	triageResponsibilityCmd.AddCommand(triageResponsibilityGetCmd)

	triageResponsibilityListCmd.Flags().IntP("limit", "l", 50, "Maximum number of triage responsibilities to return")
}
