package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "Manage Linear favorites (sidebar shortcuts)",
	Long: `Manage your favorites (sidebar shortcuts) in Linear.

Favorites appear in the Linear sidebar and provide quick access to issues,
projects, views, cycles, documents, and more. They can be organized into folders.

Examples:
  linear-cli favorite list                          # List all favorites
  linear-cli favorite list --flat                   # List without folder grouping
  linear-cli favorite add --issue ROB-123           # Add issue to favorites
  linear-cli favorite add --project PROJECT-ID      # Add project to favorites
  linear-cli favorite add --folder "My Folder"      # Create a folder
  linear-cli favorite update FAV-ID --sort-order 5  # Reorder a favorite
  linear-cli favorite remove FAV-ID                 # Remove a favorite`,
}

var favoriteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List favorites",
	Long:    `List all favorites, optionally grouped by folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		limit, _ := cmd.Flags().GetInt("limit")
		flat, _ := cmd.Flags().GetBool("flat")

		favorites, err := client.GetFavorites(context.Background(), limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list favorites: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(favorites.Nodes) == 0 {
			output.Info("No favorites found", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(favorites.Nodes)
			return
		}

		if flat || plaintext {
			renderFavoritesFlat(favorites.Nodes, plaintext)
		} else {
			renderFavoritesGrouped(favorites.Nodes)
		}
	},
}

func renderFavoritesFlat(favorites []api.Favorite, plaintext bool) {
	if plaintext {
		fmt.Println("# Favorites")
		fmt.Println("Type\tTitle\tDetail\tFolder\tID")
		for _, f := range favorites {
			folder := ""
			if f.Parent != nil {
				folder = f.Parent.Title
			} else if f.FolderName != "" {
				folder = f.FolderName
			}
			detail := getFavoriteDetail(&f)
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", f.Type, f.Title, detail, folder, f.ID)
		}
		fmt.Printf("\nTotal: %d favorites\n", len(favorites))
		return
	}

	headers := []string{"Type", "Title", "Detail", "Folder", "ID"}
	rows := make([][]string, len(favorites))

	for i, f := range favorites {
		folder := ""
		if f.Parent != nil {
			folder = f.Parent.Title
		} else if f.FolderName != "" {
			folder = f.FolderName
		}

		typeColor := getTypeColor(f.Type)
		detail := getFavoriteDetail(&f)

		rows[i] = []string{
			typeColor.Sprint(f.Type),
			truncateString(f.Title, 30),
			truncateString(detail, 25),
			folder,
			color.New(color.FgWhite, color.Faint).Sprint(f.ID),
		}
	}

	output.Table(output.TableData{
		Headers: headers,
		Rows:    rows,
	}, false, false)

	fmt.Printf("\n%s %d favorites\n",
		color.New(color.FgGreen).Sprint("✓"),
		len(favorites))
}

func renderFavoritesGrouped(favorites []api.Favorite) {
	// Group by folder
	folders := make(map[string][]api.Favorite)
	rootFavorites := []api.Favorite{}

	for _, f := range favorites {
		if f.Type == "folder" {
			// Folders themselves go to root
			rootFavorites = append(rootFavorites, f)
		} else if f.Parent != nil {
			folders[f.Parent.Title] = append(folders[f.Parent.Title], f)
		} else if f.FolderName != "" {
			folders[f.FolderName] = append(folders[f.FolderName], f)
		} else {
			rootFavorites = append(rootFavorites, f)
		}
	}

	// Sort root favorites by sortOrder
	sort.Slice(rootFavorites, func(i, j int) bool {
		return rootFavorites[i].SortOrder < rootFavorites[j].SortOrder
	})

	// Print root-level favorites first
	if len(rootFavorites) > 0 {
		for _, f := range rootFavorites {
			if f.Type == "folder" {
				// Print folder and its contents
				fmt.Printf("\n%s %s\n",
					color.New(color.FgYellow).Sprint("📁"),
					color.New(color.FgYellow, color.Bold).Sprint(f.Title))

				if folderContents, ok := folders[f.Title]; ok {
					sort.Slice(folderContents, func(i, j int) bool {
						return folderContents[i].SortOrder < folderContents[j].SortOrder
					})
					for _, fc := range folderContents {
						printFavoriteItem(fc, "   ")
					}
				}
			} else {
				printFavoriteItem(f, "")
			}
		}
	}

	// Print any orphaned folder contents (folders we didn't see as items)
	for folderName, contents := range folders {
		// Check if we already printed this folder
		found := false
		for _, f := range rootFavorites {
			if f.Type == "folder" && f.Title == folderName {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("\n%s %s\n",
				color.New(color.FgYellow).Sprint("📁"),
				color.New(color.FgYellow, color.Bold).Sprint(folderName))
			sort.Slice(contents, func(i, j int) bool {
				return contents[i].SortOrder < contents[j].SortOrder
			})
			for _, f := range contents {
				printFavoriteItem(f, "   ")
			}
		}
	}

	fmt.Printf("\n%s %d favorites\n",
		color.New(color.FgGreen).Sprint("✓"),
		len(favorites))
}

func printFavoriteItem(f api.Favorite, indent string) {
	typeColor := getTypeColor(f.Type)
	icon := getTypeIcon(f.Type)
	detail := getFavoriteDetail(&f)

	title := f.Title
	if detail != "" {
		title = fmt.Sprintf("%s (%s)", f.Title, detail)
	}

	fmt.Printf("%s%s %s %s\n",
		indent,
		icon,
		typeColor.Sprint(truncateString(title, 50)),
		color.New(color.FgWhite, color.Faint).Sprint(f.ID))
}

func getTypeColor(favType string) *color.Color {
	switch favType {
	case "issue":
		return color.New(color.FgCyan)
	case "project":
		return color.New(color.FgMagenta)
	case "cycle":
		return color.New(color.FgBlue)
	case "customView", "predefinedViewIssues", "predefinedViewProjects":
		return color.New(color.FgGreen)
	case "document":
		return color.New(color.FgWhite)
	case "initiative":
		return color.New(color.FgRed)
	case "label", "projectLabel":
		return color.New(color.FgYellow)
	case "user":
		return color.New(color.FgCyan)
	case "folder":
		return color.New(color.FgYellow)
	case "pullRequest":
		return color.New(color.FgMagenta)
	default:
		return color.New(color.FgWhite)
	}
}

func getTypeIcon(favType string) string {
	switch favType {
	case "issue":
		return "📋"
	case "project":
		return "📊"
	case "cycle":
		return "🔄"
	case "customView", "predefinedViewIssues", "predefinedViewProjects":
		return "👁"
	case "document":
		return "📄"
	case "initiative":
		return "🎯"
	case "label", "projectLabel":
		return "🏷"
	case "user":
		return "👤"
	case "folder":
		return "📁"
	case "pullRequest":
		return "🔀"
	default:
		return "⭐"
	}
}

func getFavoriteDetail(f *api.Favorite) string {
	switch {
	case f.Issue != nil:
		return f.Issue.Identifier
	case f.Project != nil:
		detail := f.Project.State
		if f.ProjectTab != "" {
			detail = fmt.Sprintf("%s [%s]", detail, f.ProjectTab)
		}
		return detail
	case f.Cycle != nil:
		return fmt.Sprintf("#%d", f.Cycle.Number)
	case f.CustomView != nil:
		return f.CustomView.ModelName
	case f.PredefinedViewType != "":
		if f.PredefinedViewTeam != nil {
			return fmt.Sprintf("%s (%s)", f.PredefinedViewType, f.PredefinedViewTeam.Key)
		}
		return f.PredefinedViewType
	case f.Document != nil:
		return ""
	case f.Initiative != nil:
		if f.InitiativeTab != "" {
			return fmt.Sprintf("[%s]", f.InitiativeTab)
		}
		return ""
	case f.Label != nil:
		return ""
	case f.ProjectLabel != nil:
		return ""
	case f.PullRequest != nil:
		return fmt.Sprintf("#%d", f.PullRequest.Number)
	case f.User != nil && f.User.Email != "" && f.User.Email != f.Title:
		return f.User.Email
	default:
		return f.Detail
	}
}

var favoriteAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add a favorite",
	Long: `Add an item to your favorites. Specify exactly one entity type.

Examples:
  linear-cli favorite add --issue ROB-123
  linear-cli favorite add --project PROJECT-ID
  linear-cli favorite add --project PROJECT-ID --project-tab documents
  linear-cli favorite add --view VIEW-ID
  linear-cli favorite add --cycle CYCLE-ID
  linear-cli favorite add --document DOC-ID
  linear-cli favorite add --initiative INIT-ID
  linear-cli favorite add --initiative INIT-ID --initiative-tab projects
  linear-cli favorite add --label LABEL-ID
  linear-cli favorite add --project-label PROJ-LABEL-ID
  linear-cli favorite add --user USER-ID
  linear-cli favorite add --predefined-view-type myIssues
  linear-cli favorite add --predefined-view-type activeIssues --predefined-view-team TEAM-ID
  linear-cli favorite add --folder "My Folder"
  linear-cli favorite add --issue ROB-123 --parent FOLDER-ID`,
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		// Get all the entity type flags
		issueID, _ := cmd.Flags().GetString("issue")
		projectID, _ := cmd.Flags().GetString("project")
		viewID, _ := cmd.Flags().GetString("view")
		cycleID, _ := cmd.Flags().GetString("cycle")
		documentID, _ := cmd.Flags().GetString("document")
		initiativeID, _ := cmd.Flags().GetString("initiative")
		labelID, _ := cmd.Flags().GetString("label")
		projectLabelID, _ := cmd.Flags().GetString("project-label")
		userID, _ := cmd.Flags().GetString("user")
		folderName, _ := cmd.Flags().GetString("folder")
		predefinedViewType, _ := cmd.Flags().GetString("predefined-view-type")
		predefinedViewTeamID, _ := cmd.Flags().GetString("predefined-view-team")
		projectTab, _ := cmd.Flags().GetString("project-tab")
		initiativeTab, _ := cmd.Flags().GetString("initiative-tab")
		parentID, _ := cmd.Flags().GetString("parent")
		sortOrder, _ := cmd.Flags().GetFloat64("sort-order")

		// Count how many entity types were specified
		entityCount := 0
		if issueID != "" {
			entityCount++
		}
		if projectID != "" {
			entityCount++
		}
		if viewID != "" {
			entityCount++
		}
		if cycleID != "" {
			entityCount++
		}
		if documentID != "" {
			entityCount++
		}
		if initiativeID != "" {
			entityCount++
		}
		if labelID != "" {
			entityCount++
		}
		if projectLabelID != "" {
			entityCount++
		}
		if userID != "" {
			entityCount++
		}
		if folderName != "" {
			entityCount++
		}
		if predefinedViewType != "" {
			entityCount++
		}

		if entityCount == 0 {
			output.Error("Must specify exactly one entity type: --issue, --project, --view, --cycle, --document, --initiative, --label, --project-label, --user, --predefined-view-type, or --folder", plaintext, jsonOut)
			os.Exit(1)
		}
		if entityCount > 1 {
			output.Error("Specify only one entity type at a time", plaintext, jsonOut)
			os.Exit(1)
		}

		input := make(map[string]interface{})

		// Resolve issue identifier to ID if needed
		if issueID != "" {
			// Check if it looks like an identifier (contains hyphen, not a UUID)
			if strings.Contains(issueID, "-") && !isUUID(issueID) {
				issue, err := client.GetIssue(context.Background(), issueID)
				if err != nil {
					output.Error(fmt.Sprintf("Failed to find issue %s: %v", issueID, err), plaintext, jsonOut)
					os.Exit(1)
				}
				input["issueId"] = issue.ID
			} else {
				input["issueId"] = issueID
			}
		}

		if projectID != "" {
			input["projectId"] = projectID
		}
		if viewID != "" {
			input["customViewId"] = viewID
		}
		if cycleID != "" {
			input["cycleId"] = cycleID
		}
		if documentID != "" {
			input["documentId"] = documentID
		}
		if initiativeID != "" {
			input["initiativeId"] = initiativeID
		}
		if labelID != "" {
			input["labelId"] = labelID
		}
		if projectLabelID != "" {
			input["projectLabelId"] = projectLabelID
		}
		if userID != "" {
			input["userId"] = userID
		}
		if folderName != "" {
			input["folderName"] = folderName
		}
		if predefinedViewType != "" {
			input["predefinedViewType"] = predefinedViewType
		}
		if predefinedViewTeamID != "" {
			input["predefinedViewTeamId"] = predefinedViewTeamID
		}
		if projectTab != "" {
			input["projectTab"] = projectTab
		}
		if initiativeTab != "" {
			input["initiativeTab"] = initiativeTab
		}
		if parentID != "" {
			input["parentId"] = parentID
		}
		if cmd.Flags().Changed("sort-order") {
			input["sortOrder"] = sortOrder
		}

		favorite, err := client.CreateFavorite(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create favorite: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(favorite)
			return
		}

		if plaintext {
			fmt.Printf("Created favorite: %s (%s)\n", favorite.Title, favorite.ID)
			fmt.Printf("Type: %s\n", favorite.Type)
		} else {
			fmt.Printf("%s Created favorite %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(favorite.Title))
			fmt.Printf("  Type: %s\n", favorite.Type)
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(favorite.ID))
		}
	},
}

// isUUID checks if a string looks like a UUID
func isUUID(s string) bool {
	// Simple check: UUIDs are 36 chars with 4 hyphens in specific positions
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
		} else if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

var favoriteUpdateCmd = &cobra.Command{
	Use:     "update FAVORITE-ID",
	Aliases: []string{"edit"},
	Short:   "Update a favorite",
	Long: `Update a favorite's sort order, parent folder, or folder name.

Examples:
  linear-cli favorite update FAV-ID --sort-order 5
  linear-cli favorite update FAV-ID --parent FOLDER-ID
  linear-cli favorite update FAV-ID --folder-name "Renamed Folder"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		favoriteID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		input := make(map[string]interface{})

		if cmd.Flags().Changed("sort-order") {
			sortOrder, _ := cmd.Flags().GetFloat64("sort-order")
			input["sortOrder"] = sortOrder
		}
		if cmd.Flags().Changed("parent") {
			parentID, _ := cmd.Flags().GetString("parent")
			input["parentId"] = parentID
		}
		if cmd.Flags().Changed("folder-name") {
			folderName, _ := cmd.Flags().GetString("folder-name")
			input["folderName"] = folderName
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use --sort-order, --parent, or --folder-name.", plaintext, jsonOut)
			os.Exit(1)
		}

		favorite, err := client.UpdateFavorite(context.Background(), favoriteID, input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update favorite: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(favorite)
			return
		}

		if plaintext {
			fmt.Printf("Updated favorite: %s\n", favorite.Title)
			fmt.Printf("ID: %s\n", favorite.ID)
		} else {
			fmt.Printf("%s Updated favorite %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(favorite.Title))
			fmt.Printf("  ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(favorite.ID))
			if favorite.SortOrder != 0 {
				fmt.Printf("  Sort Order: %.0f\n", favorite.SortOrder)
			}
			if favorite.Parent != nil {
				fmt.Printf("  Folder: %s\n", favorite.Parent.Title)
			}
		}
	},
}

var favoriteRemoveCmd = &cobra.Command{
	Use:     "remove FAVORITE-ID",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a favorite",
	Long:    `Remove a favorite from your sidebar.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		favoriteID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		err = client.DeleteFavorite(context.Background(), favoriteID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to remove favorite: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{"success": true, "id": favoriteID, "action": "deleted"})
		} else if plaintext {
			fmt.Printf("Removed favorite %s\n", favoriteID)
		} else {
			fmt.Printf("%s Removed favorite %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgWhite, color.Faint).Sprint(favoriteID))
		}
	},
}

var favoriteGetCmd = &cobra.Command{
	Use:     "get FAVORITE-ID",
	Aliases: []string{"show"},
	Short:   "Get favorite details",
	Long:    `Get detailed information about a specific favorite.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		favoriteID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		favorite, err := client.GetFavorite(context.Background(), favoriteID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get favorite: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(favorite)
			return
		}

		if plaintext {
			fmt.Printf("# %s\n", favorite.Title)
			fmt.Printf("- **ID**: %s\n", favorite.ID)
			fmt.Printf("- **Type**: %s\n", favorite.Type)
			if favorite.Detail != "" {
				fmt.Printf("- **Detail**: %s\n", favorite.Detail)
			}
			if favorite.FolderName != "" {
				fmt.Printf("- **Folder**: %s\n", favorite.FolderName)
			}
			if favorite.ProjectTab != "" {
				fmt.Printf("- **Project Tab**: %s\n", favorite.ProjectTab)
			}
			if favorite.InitiativeTab != "" {
				fmt.Printf("- **Initiative Tab**: %s\n", favorite.InitiativeTab)
			}
			if favorite.PredefinedViewType != "" {
				fmt.Printf("- **Predefined View**: %s\n", favorite.PredefinedViewType)
			}
			if favorite.PredefinedViewTeam != nil {
				fmt.Printf("- **Predefined View Team**: %s (%s)\n", favorite.PredefinedViewTeam.Name, favorite.PredefinedViewTeam.Key)
			}
			fmt.Printf("- **Sort Order**: %.0f\n", favorite.SortOrder)
			if favorite.URL != "" {
				fmt.Printf("- **URL**: %s\n", favorite.URL)
			}
			fmt.Printf("- **Created**: %s\n", favorite.CreatedAt.Format("2006-01-02 15:04:05"))
			if favorite.Parent != nil {
				fmt.Printf("- **Parent**: %s (%s)\n", favorite.Parent.Title, favorite.Parent.ID)
			}
			if favorite.Children != nil && len(favorite.Children.Nodes) > 0 {
				fmt.Printf("- **Children**: %d items\n", len(favorite.Children.Nodes))
			}
			// Show referenced entity details
			if favorite.Issue != nil {
				fmt.Printf("- **Issue**: %s - %s\n", favorite.Issue.Identifier, favorite.Issue.Title)
			}
			if favorite.Project != nil {
				fmt.Printf("- **Project**: %s (%s)\n", favorite.Project.Name, favorite.Project.State)
			}
			if favorite.ProjectTeam != nil {
				fmt.Printf("- **Project Team**: %s (%s)\n", favorite.ProjectTeam.Name, favorite.ProjectTeam.Key)
			}
			if favorite.Cycle != nil {
				fmt.Printf("- **Cycle**: %s (#%d)\n", favorite.Cycle.Name, favorite.Cycle.Number)
			}
			if favorite.CustomView != nil {
				fmt.Printf("- **Custom View**: %s (%s)\n", favorite.CustomView.Name, favorite.CustomView.ModelName)
			}
			if favorite.Document != nil {
				fmt.Printf("- **Document**: %s\n", favorite.Document.Title)
			}
			if favorite.Initiative != nil {
				fmt.Printf("- **Initiative**: %s\n", favorite.Initiative.Name)
			}
			if favorite.Label != nil {
				fmt.Printf("- **Label**: %s\n", favorite.Label.Name)
			}
			if favorite.ProjectLabel != nil {
				fmt.Printf("- **Project Label**: %s\n", favorite.ProjectLabel.Name)
			}
			if favorite.User != nil {
				fmt.Printf("- **User**: %s (%s)\n", favorite.User.Name, favorite.User.Email)
			}
			if favorite.PullRequest != nil {
				fmt.Printf("- **Pull Request**: %s (#%d)\n", favorite.PullRequest.Title, favorite.PullRequest.Number)
			}
			return
		}

		// Rich display
		fmt.Println()
		icon := getTypeIcon(favorite.Type)
		fmt.Printf("%s %s %s\n",
			icon,
			color.New(color.FgCyan, color.Bold).Sprint("Favorite:"),
			color.New(color.FgWhite, color.Bold).Sprint(favorite.Title))
		fmt.Println(strings.Repeat("-", 50))

		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("ID:"), favorite.ID)
		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Type:"), favorite.Type)

		if favorite.Detail != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Detail:"), favorite.Detail)
		}
		if favorite.FolderName != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Folder:"), favorite.FolderName)
		}
		if favorite.ProjectTab != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Project Tab:"), favorite.ProjectTab)
		}
		if favorite.InitiativeTab != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Initiative Tab:"), favorite.InitiativeTab)
		}
		if favorite.PredefinedViewType != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Predefined View:"), favorite.PredefinedViewType)
		}
		if favorite.PredefinedViewTeam != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("Predefined View Team:"),
				favorite.PredefinedViewTeam.Name,
				color.New(color.FgWhite, color.Faint).Sprint(favorite.PredefinedViewTeam.Key))
		}
		fmt.Printf("%s %.0f\n", color.New(color.Bold).Sprint("Sort Order:"), favorite.SortOrder)

		if favorite.URL != "" {
			fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("URL:"), favorite.URL)
		}

		fmt.Printf("%s %s\n", color.New(color.Bold).Sprint("Created:"), favorite.CreatedAt.Format("2006-01-02 15:04:05"))

		if favorite.Parent != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("Parent:"),
				favorite.Parent.Title,
				color.New(color.FgWhite, color.Faint).Sprint(favorite.Parent.ID))
		}

		// Show referenced entity details
		if favorite.Issue != nil {
			fmt.Printf("%s %s - %s\n",
				color.New(color.Bold).Sprint("Issue:"),
				color.New(color.FgCyan).Sprint(favorite.Issue.Identifier),
				favorite.Issue.Title)
		}
		if favorite.Project != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("Project:"),
				favorite.Project.Name,
				favorite.Project.State)
		}
		if favorite.ProjectTeam != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("Project Team:"),
				favorite.ProjectTeam.Name,
				color.New(color.FgWhite, color.Faint).Sprint(favorite.ProjectTeam.Key))
		}
		if favorite.Cycle != nil {
			fmt.Printf("%s %s (#%d)\n",
				color.New(color.Bold).Sprint("Cycle:"),
				favorite.Cycle.Name,
				favorite.Cycle.Number)
		}
		if favorite.CustomView != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("Custom View:"),
				favorite.CustomView.Name,
				favorite.CustomView.ModelName)
		}
		if favorite.Document != nil {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("Document:"),
				favorite.Document.Title)
		}
		if favorite.Initiative != nil {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("Initiative:"),
				favorite.Initiative.Name)
		}
		if favorite.Label != nil {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("Label:"),
				favorite.Label.Name)
		}
		if favorite.ProjectLabel != nil {
			fmt.Printf("%s %s\n",
				color.New(color.Bold).Sprint("Project Label:"),
				favorite.ProjectLabel.Name)
		}
		if favorite.User != nil {
			fmt.Printf("%s %s (%s)\n",
				color.New(color.Bold).Sprint("User:"),
				favorite.User.Name,
				favorite.User.Email)
		}
		if favorite.PullRequest != nil {
			fmt.Printf("%s %s (#%d)\n",
				color.New(color.Bold).Sprint("Pull Request:"),
				favorite.PullRequest.Title,
				favorite.PullRequest.Number)
		}

		if favorite.Children != nil && len(favorite.Children.Nodes) > 0 {
			fmt.Printf("\n%s\n", color.New(color.Bold).Sprint("Children:"))
			for _, child := range favorite.Children.Nodes {
				childIcon := getTypeIcon(child.Type)
				fmt.Printf("  %s %s (%s)\n", childIcon, child.Title, child.ID)
			}
		}

		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(favoriteCmd)
	favoriteCmd.AddCommand(favoriteListCmd)
	favoriteCmd.AddCommand(favoriteGetCmd)
	favoriteCmd.AddCommand(favoriteAddCmd)
	favoriteCmd.AddCommand(favoriteUpdateCmd)
	favoriteCmd.AddCommand(favoriteRemoveCmd)

	// List flags
	favoriteListCmd.Flags().IntP("limit", "l", 100, "Maximum number of favorites to return")
	favoriteListCmd.Flags().Bool("flat", false, "Show flat list without folder grouping")

	// Add flags - entity types
	favoriteAddCmd.Flags().String("issue", "", "Issue ID or identifier (e.g., ROB-123)")
	favoriteAddCmd.Flags().String("project", "", "Project ID")
	favoriteAddCmd.Flags().String("view", "", "Custom view ID")
	favoriteAddCmd.Flags().String("cycle", "", "Cycle ID")
	favoriteAddCmd.Flags().String("document", "", "Document ID")
	favoriteAddCmd.Flags().String("initiative", "", "Initiative ID")
	favoriteAddCmd.Flags().String("label", "", "Issue label ID")
	favoriteAddCmd.Flags().String("project-label", "", "Project label ID")
	favoriteAddCmd.Flags().String("user", "", "User ID")
	favoriteAddCmd.Flags().String("folder", "", "Create a folder with this name")
	favoriteAddCmd.Flags().String("predefined-view-type", "", "Predefined view type (e.g., myIssues, activeIssues)")
	favoriteAddCmd.Flags().String("predefined-view-team", "", "Team ID for predefined view")
	// Add flags - modifiers
	favoriteAddCmd.Flags().String("project-tab", "", "Tab for project favorites: issues, documents, updates, customers")
	favoriteAddCmd.Flags().String("initiative-tab", "", "Tab for initiative favorites: overview, projects, updates")
	favoriteAddCmd.Flags().String("parent", "", "Parent folder ID")
	favoriteAddCmd.Flags().Float64("sort-order", 0, "Sort order (lower = higher in list)")

	// Update flags
	favoriteUpdateCmd.Flags().Float64("sort-order", 0, "New sort order")
	favoriteUpdateCmd.Flags().String("parent", "", "Move to parent folder ID")
	favoriteUpdateCmd.Flags().String("folder-name", "", "Rename folder")
}
