package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var graphqlCmd = &cobra.Command{
	Use:     "graphql QUERY",
	Aliases: []string{"gql", "gl"},
	Short:   "Run an arbitrary GraphQL query against the Linear API",
	Long: `Run an arbitrary GraphQL query or mutation against the Linear API and print the raw JSON response.

Examples:
  linear-cli graphql 'query { viewer { id name email } }'
  linear-cli graphql 'query($id: String!) { issue(id: $id) { title } }' -v '{"id": "abc"}'
  linear-cli gql 'query { teams { nodes { id name } } }'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		query := args[0]

		// Get auth header
		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Create API client
		client := api.NewClient(authHeader)

		// Parse variables if provided
		var variables map[string]interface{}
		varsStr, _ := cmd.Flags().GetString("variables")
		if varsStr != "" {
			if err := json.Unmarshal([]byte(varsStr), &variables); err != nil {
				output.Error(fmt.Sprintf("Failed to parse variables JSON: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		// Execute the raw query
		data, err := client.ExecuteRaw(context.Background(), query, variables)
		if err != nil {
			output.Error(fmt.Sprintf("GraphQL request failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Pretty-print the raw JSON response
		var pretty json.RawMessage
		if err := json.Unmarshal(data, &pretty); err != nil {
			// If we can't re-parse, just print as-is
			fmt.Println(string(data))
			return
		}

		formatted, err := json.MarshalIndent(pretty, "", "  ")
		if err != nil {
			fmt.Println(string(data))
			return
		}
		fmt.Println(string(formatted))
	},
}

func init() {
	rootCmd.AddCommand(graphqlCmd)

	graphqlCmd.Flags().StringP("variables", "v", "", "JSON string of GraphQL variables")
}
