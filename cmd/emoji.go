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

var emojiCmd = &cobra.Command{
	Use:   "emoji",
	Short: "Manage Linear custom emojis",
	Long: `Manage custom emojis in the workspace.

Examples:
  linear-cli emoji list          # List all custom emojis`,
}

var emojiListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List custom emojis",
	Long:    `List all custom emojis in the workspace.`,
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

		emojis, err := client.GetEmojis(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list emojis: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(emojis.Nodes)
			return
		}

		if len(emojis.Nodes) == 0 {
			if plaintext {
				fmt.Println("No emojis found")
			} else {
				fmt.Printf("\n%s No emojis found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Created"}
		rows := [][]string{}
		for _, e := range emojis.Nodes {
			rows = append(rows, []string{
				e.Name,
				formatDateShort(e.CreatedAt.Format(time.RFC3339)),
			})
		}

		if plaintext {
			fmt.Println("# Emojis")
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
	rootCmd.AddCommand(emojiCmd)
	emojiCmd.AddCommand(emojiListCmd)

	emojiListCmd.Flags().IntP("limit", "l", 50, "Maximum number of emojis to return")
}
