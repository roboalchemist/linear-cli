package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/roboalchemist/linear-cli/pkg/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage Linear documents",
	Long: `Create, list, view, search, update, and delete Linear documents.

Examples:
  linear-cli document list                          # List all documents
  linear-cli document list --project PROJECT-ID     # List documents for a project
  linear-cli document get DOC-ID                    # View a document
  linear-cli document search "spec"                 # Search documents
  linear-cli document create --title "My Doc"       # Create a document
  linear-cli document update DOC-ID --title "New"   # Update a document
  linear-cli document delete DOC-ID                 # Delete a document`,
}

var documentListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List documents",
	Long:    `List Linear documents with optional filtering.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		filter := buildDocumentFilter(cmd)

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 50
		}

		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		docs, err := client.GetDocuments(context.Background(), filter, limit, "", orderBy)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch documents: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		renderDocumentCollection(docs, plaintext, jsonOut, "No documents found", "documents", "# Documents")
	},
}

func renderDocumentCollection(docs *api.Documents, plaintext, jsonOut bool, emptyMessage, summaryLabel, plaintextTitle string) {
	if len(docs.Nodes) == 0 {
		output.Info(emptyMessage, plaintext, jsonOut)
		return
	}

	if jsonOut {
		output.JSON(docs.Nodes)
		return
	}

	if plaintext {
		fmt.Println(plaintextTitle)
		for _, doc := range docs.Nodes {
			fmt.Printf("## %s\n", doc.Title)
			fmt.Printf("- **ID**: %s\n", doc.ID)
			if doc.Project != nil {
				fmt.Printf("- **Project**: %s\n", doc.Project.Name)
			}
			if doc.Team != nil {
				fmt.Printf("- **Team**: %s\n", doc.Team.Key)
			}
			if doc.Creator != nil {
				fmt.Printf("- **Creator**: %s\n", doc.Creator.Name)
			}
			fmt.Printf("- **Updated**: %s\n", doc.UpdatedAt.Format("2006-01-02"))
			if doc.URL != "" {
				fmt.Printf("- **URL**: %s\n", doc.URL)
			}
			fmt.Println()
		}
		fmt.Printf("\nTotal: %d %s\n", len(docs.Nodes), summaryLabel)
		return
	}

	headers := []string{"Title", "Project", "Team", "Creator", "Updated", "URL"}
	rows := make([][]string, len(docs.Nodes))

	for i, doc := range docs.Nodes {
		project := ""
		if doc.Project != nil {
			project = truncateString(doc.Project.Name, 25)
		}

		team := ""
		if doc.Team != nil {
			team = doc.Team.Key
		}

		creator := ""
		if doc.Creator != nil {
			creator = doc.Creator.Name
		}

		icon := ""
		if doc.Icon != nil && *doc.Icon != "" {
			icon = *doc.Icon + " "
		}

		rows[i] = []string{
			icon + truncateString(doc.Title, 35),
			project,
			team,
			creator,
			doc.UpdatedAt.Format("2006-01-02"),
			doc.URL,
		}
	}

	tableData := output.TableData{
		Headers: headers,
		Rows:    rows,
	}

	output.Table(tableData, false, false)

	fmt.Printf("\n%s %d %s\n",
		color.New(color.FgGreen).Sprint("✓"),
		len(docs.Nodes),
		summaryLabel)

	if docs.PageInfo.HasNextPage {
		fmt.Printf("%s Use --limit to see more results\n",
			color.New(color.FgYellow).Sprint("ℹ️"))
	}
}

var documentGetCmd = &cobra.Command{
	Use:     "get [document-id]",
	Aliases: []string{"show"},
	Short:   "Get document details",
	Long:    `Get detailed information about a specific document, including its full content.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		doc, err := client.GetDocument(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(doc)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n\n", doc.Title)

			fmt.Printf("## Metadata\n")
			fmt.Printf("- **ID**: %s\n", doc.ID)
			if doc.Creator != nil {
				fmt.Printf("- **Creator**: %s (%s)\n", doc.Creator.Name, doc.Creator.Email)
			}
			if doc.UpdatedBy != nil {
				fmt.Printf("- **Updated by**: %s (%s)\n", doc.UpdatedBy.Name, doc.UpdatedBy.Email)
			}
			fmt.Printf("- **Created**: %s\n", doc.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("- **Updated**: %s\n", doc.UpdatedAt.Format("2006-01-02 15:04:05"))
			if doc.Project != nil {
				fmt.Printf("- **Project**: %s\n", doc.Project.Name)
			}
			if doc.Team != nil {
				fmt.Printf("- **Team**: %s (%s)\n", doc.Team.Name, doc.Team.Key)
			}
			if doc.URL != "" {
				fmt.Printf("- **URL**: %s\n", doc.URL)
			}

			if doc.Content != "" {
				fmt.Printf("\n## Content\n%s\n", doc.Content)
			}
			return
		}

		// Rich display
		icon := ""
		if doc.Icon != nil && *doc.Icon != "" {
			icon = *doc.Icon + " "
		}

		fmt.Printf("%s%s\n",
			icon,
			color.New(color.FgWhite, color.Bold).Sprint(doc.Title))

		fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Details:"))

		if doc.Creator != nil {
			fmt.Printf("Creator: %s\n",
				color.New(color.FgCyan).Sprint(doc.Creator.Name))
		}

		if doc.UpdatedBy != nil {
			fmt.Printf("Updated by: %s\n",
				color.New(color.FgCyan).Sprint(doc.UpdatedBy.Name))
		}

		fmt.Printf("Created: %s\n", doc.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", doc.UpdatedAt.Format("2006-01-02 15:04:05"))

		if doc.Project != nil {
			fmt.Printf("Project: %s\n",
				color.New(color.FgBlue).Sprint(doc.Project.Name))
		}

		if doc.Team != nil {
			fmt.Printf("Team: %s\n",
				color.New(color.FgMagenta).Sprint(doc.Team.Name))
		}

		if doc.URL != "" {
			fmt.Printf("URL: %s\n",
				color.New(color.FgBlue, color.Underline).Sprint(doc.URL))
		}

		if doc.Content != "" {
			fmt.Printf("\n%s\n", color.New(color.FgYellow).Sprint("Content:"))
			fmt.Printf("%s\n", doc.Content)
		}
	},
}

var documentSearchCmd = &cobra.Command{
	Use:     "search [query]",
	Aliases: []string{"find"},
	Short:   "Search documents by keyword",
	Long: `Perform a full-text search across Linear documents.

Examples:
  linear-cli document search "spec"
  linear-cli document search "onboarding" --team ENG
  linear-cli document search "design" --json`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			output.Error("Search query is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		limit, _ := cmd.Flags().GetInt("limit")
		if limit == 0 {
			limit = 50
		}

		sortBy, _ := cmd.Flags().GetString("sort")
		orderBy := ""
		if sortBy != "" {
			switch sortBy {
			case "created", "createdAt":
				orderBy = "createdAt"
			case "updated", "updatedAt":
				orderBy = "updatedAt"
			case "linear":
				orderBy = ""
			default:
				output.Error(fmt.Sprintf("Invalid sort option: %s. Valid options are: linear, created, updated", sortBy), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		teamID, _ := cmd.Flags().GetString("team")
		includeComments, _ := cmd.Flags().GetBool("include-comments")

		docs, err := client.SearchDocuments(context.Background(), query, limit, "", orderBy, teamID, includeComments)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to search documents: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		emptyMsg := fmt.Sprintf("No documents found matching %q", query)
		renderDocumentCollection(docs, plaintext, jsonOut, emptyMsg, "matches", "# Search Results")
	},
}

var documentCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new document",
	Long: `Create a new document in Linear.

The content can be provided inline via --content or read from a markdown file via --content-file.
Use --content-file - to read from stdin.

Examples:
  linear-cli document create --title "My Doc" --content "Some text"
  linear-cli document create --title "My Doc" --content-file document.md
  cat doc.md | linear-cli document create --title "My Doc" --content-file -`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		title, _ := cmd.Flags().GetString("title")
		contentFlag, _ := cmd.Flags().GetString("content")
		filePath, _ := cmd.Flags().GetString("content-file")
		content, err := resolveBodyFromFlags(contentFlag, cmd.Flags().Changed("content"), filePath, "content", "content-file")
		if err != nil {
			output.Error(err.Error(), plaintext, jsonOut)
			os.Exit(1)
		}
		projectID, _ := cmd.Flags().GetString("project")
		issueID, _ := cmd.Flags().GetString("issue")
		teamKey, _ := cmd.Flags().GetString("team")
		icon, _ := cmd.Flags().GetString("icon")
		docColor, _ := cmd.Flags().GetString("color")

		if title == "" {
			output.Error("Title is required (--title)", plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{
			"title": title,
		}

		if content != "" {
			input["content"] = content
		}

		if projectID != "" {
			input["projectId"] = projectID
		}

		// NOTE: issueId is intentionally NOT added to create input.
		// Linear API's documentCreate mutation doesn't support issueId despite
		// the schema advertising it. We work around this by creating the document
		// first, then updating it to link the issue.

		if teamKey != "" {
			// Resolve team key to ID
			team, err := client.GetTeam(context.Background(), teamKey)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to find team '%s': %v", teamKey, err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["teamId"] = team.ID
		}

		if icon != "" {
			input["icon"] = icon
		}

		if docColor != "" {
			input["color"] = docColor
		}

		doc, err := client.CreateDocument(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// If issue ID was provided, update the document to link it
		// (workaround for Linear API limitation where documentCreate doesn't accept issueId)
		if issueID != "" {
			updateInput := map[string]interface{}{
				"issueId": issueID,
			}
			doc, err = client.UpdateDocument(context.Background(), doc.ID, updateInput)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to link document to issue: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
		}

		if jsonOut {
			output.JSON(doc)
		} else if plaintext {
			fmt.Printf("Created document: %s\n", doc.Title)
			fmt.Printf("ID: %s\n", doc.ID)
			if doc.URL != "" {
				fmt.Printf("URL: %s\n", doc.URL)
			}
		} else {
			fmt.Printf("%s Created document: %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Bold).Sprint(doc.Title))
			if doc.URL != "" {
				fmt.Printf("  URL: %s\n",
					color.New(color.FgBlue, color.Underline).Sprint(doc.URL))
			}
		}
	},
}

var documentUpdateCmd = &cobra.Command{
	Use:     "update [document-id]",
	Aliases: []string{"edit"},
	Short:   "Update a document",
	Long: `Update fields of an existing document.

The content can be provided inline via --content or read from a markdown file via --content-file.
Use --content-file - to read from stdin.

Examples:
  linear-cli document update DOC-ID --title "New Title"
  linear-cli document update DOC-ID --content "Updated content"
  linear-cli document update DOC-ID --content-file updated-doc.md
  linear-cli document update DOC-ID --icon "📝" --color "#ff0000"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		input := make(map[string]interface{})

		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			input["title"] = title
		}

		filePath, _ := cmd.Flags().GetString("content-file")
		if cmd.Flags().Changed("content") || filePath != "" {
			contentFlag, _ := cmd.Flags().GetString("content")
			content, err := resolveBodyFromFlags(contentFlag, cmd.Flags().Changed("content"), filePath, "content", "content-file")
			if err != nil {
				output.Error(err.Error(), plaintext, jsonOut)
				os.Exit(1)
			}
			input["content"] = content
		}

		if cmd.Flags().Changed("icon") {
			icon, _ := cmd.Flags().GetString("icon")
			input["icon"] = icon
		}

		if cmd.Flags().Changed("color") {
			docColor, _ := cmd.Flags().GetString("color")
			input["color"] = docColor
		}

		if cmd.Flags().Changed("project") {
			projectID, _ := cmd.Flags().GetString("project")
			input["projectId"] = projectID
		}

		if cmd.Flags().Changed("issue") {
			issueID, _ := cmd.Flags().GetString("issue")
			input["issueId"] = issueID
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use flags to specify what to update.", plaintext, jsonOut)
			os.Exit(1)
		}

		doc, err := client.UpdateDocument(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(doc)
		} else if plaintext {
			fmt.Printf("Updated document: %s\n", doc.Title)
			fmt.Printf("ID: %s\n", doc.ID)
			if doc.URL != "" {
				fmt.Printf("URL: %s\n", doc.URL)
			}
		} else {
			fmt.Printf("%s Updated document: %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Bold).Sprint(doc.Title))
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(doc.ID))
			if doc.URL != "" {
				fmt.Printf("  URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(doc.URL))
			}
		}
	},
}

var documentDeleteCmd = &cobra.Command{
	Use:   "delete [document-id]",
	Short: "Delete a document",
	Long:  `Delete a document from Linear.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		err = client.DeleteDocument(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete document: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "id": args[0], "action": "deleted"})
		} else if plaintext {
			fmt.Printf("Deleted document %s\n", args[0])
		} else {
			fmt.Printf("%s Deleted document %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Faint).Sprint(args[0]))
		}
	},
}

func buildDocumentFilter(cmd *cobra.Command) map[string]interface{} {
	filter := make(map[string]interface{})

	if projectID, _ := cmd.Flags().GetString("project"); projectID != "" {
		filter["project"] = map[string]interface{}{"id": map[string]interface{}{"eq": projectID}}
	}

	if issueID, _ := cmd.Flags().GetString("issue"); issueID != "" {
		filter["issue"] = map[string]interface{}{"id": map[string]interface{}{"eq": issueID}}
	}

	if teamKey, _ := cmd.Flags().GetString("team"); teamKey != "" {
		filter["team"] = map[string]interface{}{"key": map[string]interface{}{"eq": teamKey}}
	}

	newerThan, _ := cmd.Flags().GetString("newer-than")
	createdAt, err := utils.ParseTimeExpression(newerThan)
	if err != nil {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		output.Error(fmt.Sprintf("Invalid newer-than value: %v", err), plaintext, jsonOut)
		os.Exit(1)
	}
	if createdAt != "" {
		filter["createdAt"] = map[string]interface{}{"gte": createdAt}
	}

	return filter
}

func init() {
	rootCmd.AddCommand(documentCmd)
	documentCmd.AddCommand(documentListCmd)
	documentCmd.AddCommand(documentGetCmd)
	documentCmd.AddCommand(documentSearchCmd)
	documentCmd.AddCommand(documentCreateCmd)
	documentCmd.AddCommand(documentUpdateCmd)
	documentCmd.AddCommand(documentDeleteCmd)

	// List command flags
	documentListCmd.Flags().String("project", "", "Filter by project ID")
	documentListCmd.Flags().String("issue", "", "Filter by issue ID")
	documentListCmd.Flags().StringP("team", "t", "", "Filter by team key")
	documentListCmd.Flags().IntP("limit", "l", 50, "Maximum number of documents to return")
	documentListCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	documentListCmd.Flags().StringP("newer-than", "n", "", "Show documents created after this time (default: 6_months_ago, use 'all_time' for no filter)")

	// Search command flags
	documentSearchCmd.Flags().StringP("team", "t", "", "Filter by team ID")
	documentSearchCmd.Flags().IntP("limit", "l", 50, "Maximum number of results to return")
	documentSearchCmd.Flags().StringP("sort", "o", "linear", "Sort order: linear (default), created, updated")
	documentSearchCmd.Flags().Bool("include-comments", false, "Include document comments in search")

	// Create command flags
	documentCreateCmd.Flags().String("title", "", "Document title (required)")
	documentCreateCmd.Flags().String("content", "", "Document content (markdown)")
	documentCreateCmd.Flags().String("content-file", "", "Read content from a markdown file (use - for stdin)")
	documentCreateCmd.Flags().String("project", "", "Project ID to associate with")
	documentCreateCmd.Flags().String("issue", "", "Issue ID to associate with")
	documentCreateCmd.Flags().StringP("team", "t", "", "Team key to associate with")
	documentCreateCmd.Flags().String("icon", "", "Document icon (emoji)")
	documentCreateCmd.Flags().String("color", "", "Document icon color (hex)")
	_ = documentCreateCmd.MarkFlagRequired("title")

	// Update command flags
	documentUpdateCmd.Flags().String("title", "", "New title for the document")
	documentUpdateCmd.Flags().String("content", "", "New content for the document (markdown)")
	documentUpdateCmd.Flags().String("content-file", "", "Read content from a markdown file (use - for stdin)")
	documentUpdateCmd.Flags().String("icon", "", "New icon (emoji)")
	documentUpdateCmd.Flags().String("color", "", "New icon color (hex)")
	documentUpdateCmd.Flags().String("project", "", "Project ID to associate with")
	documentUpdateCmd.Flags().String("issue", "", "Issue ID to associate with")

	// Delete has no extra flags
}
