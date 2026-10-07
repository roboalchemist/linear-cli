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

var customerCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a customer",
	Long: `Create a new customer.

Examples:
  linear-cli customer create --name "Acme Inc" --domains acme.com,acme.io
  linear-cli customer create --name "Acme Inc" --revenue 1000000 --size 250 --owner-id USER-ID`,
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
		domains, _ := cmd.Flags().GetStringSlice("domains")
		externalIDs, _ := cmd.Flags().GetStringSlice("external-ids")
		ownerID, _ := cmd.Flags().GetString("owner-id")
		statusID, _ := cmd.Flags().GetString("status-id")
		tierID, _ := cmd.Flags().GetString("tier-id")
		logoURL, _ := cmd.Flags().GetString("logo-url")
		slackChannelID, _ := cmd.Flags().GetString("slack-channel-id")
		inputJSON, _ := cmd.Flags().GetString("input-json")

		input := map[string]interface{}{"name": name}
		if len(domains) > 0 {
			input["domains"] = domains
		}
		if len(externalIDs) > 0 {
			input["externalIds"] = externalIDs
		}
		if cmd.Flags().Changed("revenue") {
			v, _ := cmd.Flags().GetInt("revenue")
			input["revenue"] = v
		}
		if cmd.Flags().Changed("size") {
			v, _ := cmd.Flags().GetInt("size")
			input["size"] = v
		}
		if ownerID != "" {
			input["ownerId"] = ownerID
		}
		if statusID != "" {
			input["statusId"] = statusID
		}
		if tierID != "" {
			input["tierId"] = tierID
		}
		if logoURL != "" {
			input["logoUrl"] = logoURL
		}
		if slackChannelID != "" {
			input["slackChannelId"] = slackChannelID
		}
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		customer, err := client.CreateCustomer(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create customer: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(customer)
		} else {
			output.Success(fmt.Sprintf("Created customer %s",
				color.New(color.FgWhite, color.Bold).Sprint(customer.Name)), plaintext, jsonOut)
		}
	},
}

var customerUpdateCmd = &cobra.Command{
	Use:     "update CUSTOMER-ID",
	Aliases: []string{"edit"},
	Short:   "Update a customer",
	Long: `Update a customer's name, domains, revenue, size, owner, status, or tier.

Examples:
  linear-cli customer update CUSTOMER-ID --name "Acme Corporation"
  linear-cli customer update CUSTOMER-ID --domains acme.com --status-id STATUS-ID`,
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
		if cmd.Flags().Changed("domains") {
			v, _ := cmd.Flags().GetStringSlice("domains")
			input["domains"] = v
		}
		if cmd.Flags().Changed("external-ids") {
			v, _ := cmd.Flags().GetStringSlice("external-ids")
			input["externalIds"] = v
		}
		if cmd.Flags().Changed("revenue") {
			v, _ := cmd.Flags().GetInt("revenue")
			input["revenue"] = v
		}
		if cmd.Flags().Changed("size") {
			v, _ := cmd.Flags().GetInt("size")
			input["size"] = v
		}
		if cmd.Flags().Changed("owner-id") {
			v, _ := cmd.Flags().GetString("owner-id")
			input["ownerId"] = v
		}
		if cmd.Flags().Changed("status-id") {
			v, _ := cmd.Flags().GetString("status-id")
			input["statusId"] = v
		}
		if cmd.Flags().Changed("tier-id") {
			v, _ := cmd.Flags().GetString("tier-id")
			input["tierId"] = v
		}
		if cmd.Flags().Changed("logo-url") {
			v, _ := cmd.Flags().GetString("logo-url")
			input["logoUrl"] = v
		}
		if cmd.Flags().Changed("slack-channel-id") {
			v, _ := cmd.Flags().GetString("slack-channel-id")
			input["slackChannelId"] = v
		}
		inputJSON, _ := cmd.Flags().GetString("input-json")
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}
		if len(input) == 0 {
			output.Error("No fields to update. Use --name, --domains, --external-ids, --revenue, --size, --owner-id, --status-id, --tier-id, --logo-url, --slack-channel-id, or --input-json.", plaintext, jsonOut)
			os.Exit(1)
		}

		customer, err := client.UpdateCustomer(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update customer: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(customer)
		} else {
			output.Success(fmt.Sprintf("Updated customer %s",
				color.New(color.FgWhite, color.Bold).Sprint(customer.Name)), plaintext, jsonOut)
		}
	},
}

var customerDeleteCmd = &cobra.Command{
	Use:     "delete CUSTOMER-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a customer",
	Long:    `Delete a customer.`,
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
		if err := client.DeleteCustomer(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete customer: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted customer", plaintext, jsonOut)
	},
}

var customerMergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge two customers",
	Long: `Merge a source customer into a target customer.

All needs from the source are transferred to the target, and the source
customer is archived.

Examples:
  linear-cli customer merge --source CUSTOMER-ID --target CUSTOMER-ID`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		source, _ := cmd.Flags().GetString("source")
		target, _ := cmd.Flags().GetString("target")

		customer, err := client.MergeCustomer(context.Background(), source, target)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to merge customers: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(customer)
		} else {
			output.Success(fmt.Sprintf("Merged customers into %s",
				color.New(color.FgWhite, color.Bold).Sprint(customer.Name)), plaintext, jsonOut)
		}
	},
}

var customerStatusCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a customer status",
	Long: `Create a new customer status.

Examples:
  linear-cli customer status create --name "Active" --color "#22c55e"`,
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
		statusColor, _ := cmd.Flags().GetString("color")
		description, _ := cmd.Flags().GetString("description")
		displayName, _ := cmd.Flags().GetString("display-name")
		inputJSON, _ := cmd.Flags().GetString("input-json")

		input := map[string]interface{}{
			"name":  name,
			"color": statusColor,
		}
		if description != "" {
			input["description"] = description
		}
		if displayName != "" {
			input["displayName"] = displayName
		}
		if cmd.Flags().Changed("position") {
			v, _ := cmd.Flags().GetFloat64("position")
			input["position"] = v
		}
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		status, err := client.CreateCustomerStatus(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create customer status: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(status)
		} else {
			output.Success(fmt.Sprintf("Created customer status %s",
				color.New(color.FgWhite, color.Bold).Sprint(status.Name)), plaintext, jsonOut)
		}
	},
}

var customerStatusDeleteCmd = &cobra.Command{
	Use:     "delete STATUS-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a customer status",
	Long:    `Delete a customer status. The status must not be in use and cannot be the last remaining status.`,
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
		if err := client.DeleteCustomerStatus(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete customer status: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted customer status", plaintext, jsonOut)
	},
}

var customerTierCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a customer tier",
	Long: `Create a new customer tier.

Examples:
  linear-cli customer tier create --name "Enterprise" --color "#0ea5e9"`,
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
		tierColor, _ := cmd.Flags().GetString("color")
		description, _ := cmd.Flags().GetString("description")
		displayName, _ := cmd.Flags().GetString("display-name")
		inputJSON, _ := cmd.Flags().GetString("input-json")

		input := map[string]interface{}{
			"name":  name,
			"color": tierColor,
		}
		if description != "" {
			input["description"] = description
		}
		if displayName != "" {
			input["displayName"] = displayName
		}
		if cmd.Flags().Changed("position") {
			v, _ := cmd.Flags().GetFloat64("position")
			input["position"] = v
		}
		if err := mergeJSONInput(input, inputJSON); err != nil {
			output.Error(fmt.Sprintf("Invalid --input-json: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		tier, err := client.CreateCustomerTier(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create customer tier: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(tier)
		} else {
			output.Success(fmt.Sprintf("Created customer tier %s",
				color.New(color.FgWhite, color.Bold).Sprint(tier.Name)), plaintext, jsonOut)
		}
	},
}

var customerTierDeleteCmd = &cobra.Command{
	Use:     "delete TIER-ID",
	Aliases: []string{"rm"},
	Short:   "Delete a customer tier",
	Long:    `Delete a customer tier. The tier must not be in use by any customers.`,
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
		if err := client.DeleteCustomerTier(context.Background(), args[0]); err != nil {
			output.Error(fmt.Sprintf("Failed to delete customer tier: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		output.Success("Deleted customer tier", plaintext, jsonOut)
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
	customerCmd.AddCommand(customerCreateCmd)
	customerCmd.AddCommand(customerUpdateCmd)
	customerCmd.AddCommand(customerDeleteCmd)
	customerCmd.AddCommand(customerMergeCmd)
	customerCmd.AddCommand(customerNeedCmd)
	customerCmd.AddCommand(customerStatusCmd)
	customerCmd.AddCommand(customerTierCmd)

	customerNeedCmd.AddCommand(customerNeedListCmd)
	customerStatusCmd.AddCommand(customerStatusListCmd)
	customerStatusCmd.AddCommand(customerStatusCreateCmd)
	customerStatusCmd.AddCommand(customerStatusDeleteCmd)
	customerTierCmd.AddCommand(customerTierListCmd)
	customerTierCmd.AddCommand(customerTierCreateCmd)
	customerTierCmd.AddCommand(customerTierDeleteCmd)

	// List flags
	customerListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customers to return")
	customerListCmd.Flags().Bool("include-archived", false, "Include archived customers")

	// Create flags
	customerCreateCmd.Flags().StringP("name", "n", "", "Customer name (required)")
	customerCreateCmd.Flags().StringSlice("domains", nil, "Email domains (comma-separated)")
	customerCreateCmd.Flags().StringSlice("external-ids", nil, "External system identifiers (comma-separated)")
	customerCreateCmd.Flags().Int("revenue", 0, "Annual revenue in dollars")
	customerCreateCmd.Flags().Int("size", 0, "Organization size (e.g. employees)")
	customerCreateCmd.Flags().String("owner-id", "", "Owner user ID")
	customerCreateCmd.Flags().String("status-id", "", "Customer status ID")
	customerCreateCmd.Flags().String("tier-id", "", "Customer tier ID")
	customerCreateCmd.Flags().String("logo-url", "", "URL of the customer's logo image")
	customerCreateCmd.Flags().String("slack-channel-id", "", "Slack channel ID to link")
	customerCreateCmd.Flags().String("input-json", "", "Additional CustomerCreateInput fields as a JSON object")
	_ = customerCreateCmd.MarkFlagRequired("name")

	// Update flags
	customerUpdateCmd.Flags().StringP("name", "n", "", "New customer name")
	customerUpdateCmd.Flags().StringSlice("domains", nil, "New email domains (comma-separated)")
	customerUpdateCmd.Flags().StringSlice("external-ids", nil, "New external system identifiers (comma-separated)")
	customerUpdateCmd.Flags().Int("revenue", 0, "New annual revenue in dollars")
	customerUpdateCmd.Flags().Int("size", 0, "New organization size")
	customerUpdateCmd.Flags().String("owner-id", "", "New owner user ID")
	customerUpdateCmd.Flags().String("status-id", "", "New customer status ID")
	customerUpdateCmd.Flags().String("tier-id", "", "New customer tier ID")
	customerUpdateCmd.Flags().String("logo-url", "", "New logo image URL")
	customerUpdateCmd.Flags().String("slack-channel-id", "", "New Slack channel ID")
	customerUpdateCmd.Flags().String("input-json", "", "Additional CustomerUpdateInput fields as a JSON object")

	// Merge flags
	customerMergeCmd.Flags().String("source", "", "Source customer ID to merge from (required)")
	customerMergeCmd.Flags().String("target", "", "Target customer ID to merge into (required)")
	_ = customerMergeCmd.MarkFlagRequired("source")
	_ = customerMergeCmd.MarkFlagRequired("target")

	// Need list flags
	customerNeedListCmd.Flags().String("customer", "", "Filter by customer ID")
	customerNeedListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customer needs to return")

	// Status list flags
	customerStatusListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customer statuses to return")

	// Status create flags
	customerStatusCreateCmd.Flags().StringP("name", "n", "", "Status internal name (required)")
	customerStatusCreateCmd.Flags().StringP("color", "c", "", "Status color as a HEX string, e.g. #22c55e (required)")
	customerStatusCreateCmd.Flags().StringP("description", "d", "", "Status description")
	customerStatusCreateCmd.Flags().String("display-name", "", "User-facing status name")
	customerStatusCreateCmd.Flags().Float64("position", 0, "Sort position")
	customerStatusCreateCmd.Flags().String("input-json", "", "Additional CustomerStatusCreateInput fields as a JSON object")
	_ = customerStatusCreateCmd.MarkFlagRequired("name")
	_ = customerStatusCreateCmd.MarkFlagRequired("color")

	// Tier list flags
	customerTierListCmd.Flags().IntP("limit", "l", 25, "Maximum number of customer tiers to return")

	// Tier create flags
	customerTierCreateCmd.Flags().StringP("name", "n", "", "Tier internal name (required)")
	customerTierCreateCmd.Flags().StringP("color", "c", "", "Tier color as a HEX string, e.g. #0ea5e9 (required)")
	customerTierCreateCmd.Flags().StringP("description", "d", "", "Tier description")
	customerTierCreateCmd.Flags().String("display-name", "", "User-facing tier name")
	customerTierCreateCmd.Flags().Float64("position", 0, "Sort position")
	customerTierCreateCmd.Flags().String("input-json", "", "Additional CustomerTierCreateInput fields as a JSON object")
	_ = customerTierCreateCmd.MarkFlagRequired("name")
	_ = customerTierCreateCmd.MarkFlagRequired("color")
}
