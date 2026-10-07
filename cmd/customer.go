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

var customerCmd = &cobra.Command{
	Use:   "customer",
	Short: "Manage Linear customers",
	Long: `Manage Linear customers, their needs, statuses, and tiers.

Examples:
  linear-cli customer list                         # List customers
  linear-cli customer list --include-archived      # Include archived customers
  linear-cli customer get CUSTOMER-ID-OR-SLUG      # Get customer details
  linear-cli customer need list --customer ID      # List needs for a customer
  linear-cli customer status list                  # List customer statuses
  linear-cli customer tier list                    # List customer tiers`,
}

var customerListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List customers",
	Long:    `List customers, optionally including archived ones.`,
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
		includeArchived, _ := cmd.Flags().GetBool("include-archived")

		filter := map[string]interface{}{}
		if includeArchived {
			filter["includeArchived"] = true
		}

		customers, err := client.GetCustomers(context.Background(), filter, limit, "", "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list customers: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(customers.Nodes)
			return
		}

		if len(customers.Nodes) == 0 {
			if plaintext {
				fmt.Println("No customers found")
			} else {
				fmt.Printf("\n%s No customers found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Slug", "Status", "Tier", "Owner", "Domains", "Needs", "URL"}
		rows := [][]string{}

		for _, c := range customers.Nodes {
			status := ""
			if c.Status != nil {
				status = c.Status.Name
			}
			tier := ""
			if c.Tier != nil {
				tier = c.Tier.Name
			}
			owner := ""
			if c.Owner != nil {
				owner = c.Owner.Name
			}

			name := c.Name
			if !plaintext {
				name = color.New(color.FgWhite, color.Bold).Sprint(c.Name)
			}

			rows = append(rows, []string{
				name,
				c.SlugID,
				status,
				tier,
				owner,
				strings.Join(c.Domains, ","),
				fmt.Sprintf("%.0f", c.ApproximateNeedCount),
				c.URL,
			})
		}

		if plaintext {
			fmt.Println("# Customers")
			fmt.Println(strings.Join(headers, "\t"))
			for _, row := range rows {
				fmt.Println(strings.Join(row, "\t"))
			}
			return
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d customers\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(customers.Nodes))
	},
}

var customerGetCmd = &cobra.Command{
	Use:     "get CUSTOMER-ID",
	Aliases: []string{"show"},
	Short:   "Get customer details",
	Long:    `Get details for a specific customer, looked up by ID or slug.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		customerID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		customer, err := client.GetCustomer(context.Background(), customerID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get customer: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(customer)
			return
		}

		status := ""
		if customer.Status != nil {
			status = customer.Status.Name
		}
		tier := ""
		if customer.Tier != nil {
			tier = customer.Tier.Name
		}
		owner := ""
		if customer.Owner != nil {
			owner = customer.Owner.Name
			if owner == "" {
				owner = customer.Owner.DisplayName
			}
		}
		domains := strings.Join(customer.Domains, ", ")
		externalIDs := strings.Join(customer.ExternalIDs, ", ")

		if plaintext {
			fmt.Printf("# Customer: %s\n", customer.Name)
			fmt.Printf("ID: %s\n", customer.ID)
			fmt.Printf("Slug: %s\n", customer.SlugID)
			fmt.Printf("Status: %s\n", status)
			fmt.Printf("Tier: %s\n", tier)
			fmt.Printf("Owner: %s\n", owner)
			fmt.Printf("Domains: %s\n", domains)
			if externalIDs != "" {
				fmt.Printf("External IDs: %s\n", externalIDs)
			}
			if customer.Revenue != nil {
				fmt.Printf("Revenue: %d\n", *customer.Revenue)
			}
			if customer.Size != nil {
				fmt.Printf("Size: %.0f\n", *customer.Size)
			}
			fmt.Printf("Needs: %.0f\n", customer.ApproximateNeedCount)
			fmt.Printf("URL: %s\n", customer.URL)
			fmt.Printf("Created: %s\n", customer.CreatedAt.Format("2006-01-02"))
			fmt.Printf("Updated: %s\n", customer.UpdatedAt.Format("2006-01-02"))
			if customer.ArchivedAt != nil {
				fmt.Printf("Archived: %s\n", customer.ArchivedAt.Format("2006-01-02"))
			}
			if len(customer.Needs) > 0 {
				fmt.Println("\nNeeds:")
				for _, n := range customer.Needs {
					issue := ""
					if n.Issue != nil {
						issue = n.Issue.Identifier
					}
					project := ""
					if n.Project != nil {
						project = n.Project.Name
					}
					body := customerNeedBody(n)
					fmt.Printf("  %s\t%s\t%s\t%s\n", truncateString(body, 60), issue, project, n.CreatedAt.Format("2006-01-02"))
				}
			}
		} else {
			fmt.Printf("\n%s Customer: %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("🏢"),
				color.New(color.FgWhite, color.Bold).Sprint(customer.Name))
			fmt.Printf("   Slug: %s\n", color.New(color.FgCyan).Sprint(customer.SlugID))
			fmt.Printf("   Status: %s\n", status)
			fmt.Printf("   Tier: %s\n", tier)
			fmt.Printf("   Owner: %s\n", owner)
			fmt.Printf("   Domains: %s\n", domains)
			if externalIDs != "" {
				fmt.Printf("   External IDs: %s\n", externalIDs)
			}
			if customer.Revenue != nil {
				fmt.Printf("   Revenue: %d\n", *customer.Revenue)
			}
			if customer.Size != nil {
				fmt.Printf("   Size: %.0f\n", *customer.Size)
			}
			fmt.Printf("   Needs: %s\n", color.New(color.FgGreen).Sprintf("%.0f", customer.ApproximateNeedCount))
			fmt.Printf("   URL: %s\n", customer.URL)
			fmt.Printf("   Created: %s | Updated: %s\n",
				customer.CreatedAt.Format("2006-01-02"),
				customer.UpdatedAt.Format("2006-01-02"))
			if customer.ArchivedAt != nil {
				fmt.Printf("   Archived: %s\n", customer.ArchivedAt.Format("2006-01-02"))
			}

			if len(customer.Needs) > 0 {
				fmt.Printf("\n   %s Needs:\n\n", color.New(color.FgCyan, color.Bold).Sprint("📋"))
				needHeaders := []string{"Body", "Issue", "Project", "Created"}
				needRows := [][]string{}
				for _, n := range customer.Needs {
					issue := ""
					if n.Issue != nil {
						issue = n.Issue.Identifier
					}
					project := ""
					if n.Project != nil {
						project = n.Project.Name
					}
					needRows = append(needRows, []string{
						truncateString(customerNeedBody(n), 60),
						issue,
						project,
						n.CreatedAt.Format("2006-01-02"),
					})
				}
				output.Table(output.TableData{
					Headers: needHeaders,
					Rows:    needRows,
				}, plaintext, jsonOut)
			}
		}
	},
}

var customerNeedCmd = &cobra.Command{
	Use:   "need",
	Short: "Manage customer needs",
	Long:  `Manage customer needs (product requests and feedback).`,
}

var customerNeedListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List customer needs",
	Long:    `List customer needs, optionally filtered to a single customer.`,
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
		customerID, _ := cmd.Flags().GetString("customer")

		needs, err := client.GetCustomerNeeds(context.Background(), customerID, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list customer needs: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(needs.Nodes)
			return
		}

		if len(needs.Nodes) == 0 {
			if plaintext {
				fmt.Println("No customer needs found")
			} else {
				fmt.Printf("\n%s No customer needs found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Body", "Issue", "Project", "Created"}
		rows := [][]string{}

		for _, n := range needs.Nodes {
			issue := ""
			if n.Issue != nil {
				issue = n.Issue.Identifier
			}
			project := ""
			if n.Project != nil {
				project = n.Project.Name
			}
			rows = append(rows, []string{
				truncateString(customerNeedBody(n), 60),
				issue,
				project,
				n.CreatedAt.Format("2006-01-02"),
			})
		}

		if plaintext {
			fmt.Println("# Customer Needs")
			fmt.Println(strings.Join(headers, "\t"))
			for _, row := range rows {
				fmt.Println(strings.Join(row, "\t"))
			}
			return
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d customer needs\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(needs.Nodes))
	},
}

var customerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Manage customer statuses",
	Long:  `Manage customer lifecycle statuses.`,
}

var customerStatusListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List customer statuses",
	Long:    `List all customer statuses defined in the workspace.`,
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

		statuses, err := client.GetCustomerStatuses(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list customer statuses: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(statuses.Nodes)
			return
		}

		if len(statuses.Nodes) == 0 {
			if plaintext {
				fmt.Println("No customer statuses found")
			} else {
				fmt.Printf("\n%s No customer statuses found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Color", "Type", "Position", "Description"}
		rows := [][]string{}

		for _, s := range statuses.Nodes {
			name := s.Name
			if s.DisplayName != "" {
				name = s.DisplayName
			}
			statusType := ""
			if s.Type != nil {
				statusType = *s.Type
			}
			desc := ""
			if s.Description != nil {
				desc = *s.Description
			}

			if !plaintext {
				name = color.New(color.FgWhite, color.Bold).Sprint(name)
			}

			rows = append(rows, []string{
				name,
				s.Color,
				statusType,
				fmt.Sprintf("%.0f", s.Position),
				desc,
			})
		}

		if plaintext {
			fmt.Println("# Customer Statuses")
			fmt.Println(strings.Join(headers, "\t"))
			for _, row := range rows {
				fmt.Println(strings.Join(row, "\t"))
			}
			return
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d customer statuses\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(statuses.Nodes))
	},
}

var customerTierCmd = &cobra.Command{
	Use:   "tier",
	Short: "Manage customer tiers",
	Long:  `Manage customer tiers (segments).`,
}

var customerTierListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List customer tiers",
	Long:    `List all customer tiers defined in the workspace.`,
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

		tiers, err := client.GetCustomerTiers(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list customer tiers: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(tiers.Nodes)
			return
		}

		if len(tiers.Nodes) == 0 {
			if plaintext {
				fmt.Println("No customer tiers found")
			} else {
				fmt.Printf("\n%s No customer tiers found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		headers := []string{"Name", "Color", "Position"}
		rows := [][]string{}

		for _, t := range tiers.Nodes {
			name := t.Name
			if t.DisplayName != "" {
				name = t.DisplayName
			}
			if !plaintext {
				name = color.New(color.FgWhite, color.Bold).Sprint(name)
			}
			rows = append(rows, []string{
				name,
				t.Color,
				fmt.Sprintf("%.0f", t.Position),
			})
		}

		if plaintext {
			fmt.Println("# Customer Tiers")
			fmt.Println(strings.Join(headers, "\t"))
			for _, row := range rows {
				fmt.Println(strings.Join(row, "\t"))
			}
			return
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, plaintext, jsonOut)

		fmt.Printf("\n%s %d customer tiers\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(tiers.Nodes))
	},
}

// customerNeedBody returns the best available body text for a customer need,
// preferring the effective content and falling back to the raw body.
func customerNeedBody(n api.CustomerNeed) string {
	if n.Content != nil && *n.Content != "" {
		return strings.ReplaceAll(*n.Content, "\n", " ")
	}
	if n.Body != nil {
		return strings.ReplaceAll(*n.Body, "\n", " ")
	}
	return ""
}

func init() {
	rootCmd.AddCommand(customerCmd)
	customerCmd.AddCommand(customerListCmd)
	customerCmd.AddCommand(customerGetCmd)
	customerCmd.AddCommand(customerNeedCmd)
	customerCmd.AddCommand(customerStatusCmd)
	customerCmd.AddCommand(customerTierCmd)

	customerNeedCmd.AddCommand(customerNeedListCmd)
	customerStatusCmd.AddCommand(customerStatusListCmd)
	customerTierCmd.AddCommand(customerTierListCmd)

	// List flags
	customerListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customers to return")
	customerListCmd.Flags().Bool("include-archived", false, "Include archived customers")

	// Need list flags
	customerNeedListCmd.Flags().String("customer", "", "Filter by customer ID")
	customerNeedListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customer needs to return")

	// Status list flags
	customerStatusListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customer statuses to return")

	// Tier list flags
	customerTierListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customer tiers to return")
}
