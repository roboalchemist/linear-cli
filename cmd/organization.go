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

var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Show the current Linear workspace",
	Long: `Inspect the authenticated user's Linear workspace (organization).

Examples:
  linear-cli organization get
  linear-cli organization get --json`,
}

var organizationGetCmd = &cobra.Command{
	Use:     "get",
	Aliases: []string{"show"},
	Short:   "Get workspace details",
	Long:    `Print details about the current workspace.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		org, err := client.GetOrganization(context.Background())
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get organization: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(org)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n", org.Name)
			fmt.Printf("URL Key: %s\n", org.URLKey)
			fmt.Printf("ID: %s\n", org.ID)
			fmt.Printf("Created: %s\n", org.CreatedAt.Format("2006-01-02"))
			fmt.Printf("Users: %d\n", org.UserCount)
			fmt.Printf("Issues: %d\n", org.CreatedIssueCount)
			fmt.Printf("Customers: %d\n", org.CustomerCount)
			fmt.Printf("Release Channel: %s\n", org.ReleaseChannel)
			fmt.Printf("SAML: %t\n", org.SAMLEnabled)
			fmt.Printf("SCIM: %t\n", org.SCIMEnabled)
			if len(org.PreviousURLKeys) > 0 {
				fmt.Printf("Previous URL Keys: %s\n", strings.Join(org.PreviousURLKeys, ", "))
			}
			if org.LogoURL != nil && *org.LogoURL != "" {
				fmt.Printf("Logo: %s\n", *org.LogoURL)
			}
			return
		}

		fmt.Printf("\n%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("🗂  Workspace:"),
			color.New(color.FgWhite, color.Bold).Sprint(org.Name))
		fmt.Println(strings.Repeat("─", 50))
		fmt.Printf("   URL Key: %s\n", color.New(color.FgCyan).Sprint(org.URLKey))
		fmt.Printf("   ID: %s\n", org.ID)
		fmt.Printf("   Created: %s\n", org.CreatedAt.Format("2006-01-02"))
		fmt.Printf("   Users: %d\n", org.UserCount)
		fmt.Printf("   Issues: %d\n", org.CreatedIssueCount)
		fmt.Printf("   Customers: %d\n", org.CustomerCount)
		fmt.Printf("   Release Channel: %s\n", org.ReleaseChannel)
		fmt.Printf("   SAML: %t | SCIM: %t\n", org.SAMLEnabled, org.SCIMEnabled)
		if len(org.PreviousURLKeys) > 0 {
			fmt.Printf("   Previous URL Keys: %s\n", strings.Join(org.PreviousURLKeys, ", "))
		}
		if org.LogoURL != nil && *org.LogoURL != "" {
			fmt.Printf("   Logo: %s\n", color.New(color.FgBlue, color.Underline).Sprint(*org.LogoURL))
		}
	},
}

func init() {
	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(organizationGetCmd)
}
