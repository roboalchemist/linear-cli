package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var issueRelationCmd = &cobra.Command{
	Use:   "relation",
	Short: "Manage issue relationships",
	Long: `Manage relationships between Linear issues.

Relationship types:
  blocks      Issue blocks another issue
  blocked-by  Issue is blocked by another issue
  related     Issues are related
  duplicate   Issue is a duplicate of another
  parent      Set the parent of an issue
  sub-issue   Make another issue a sub-issue of this one

Examples:
  linear-cli issue relation list LIN-123
  linear-cli issue relation add LIN-123 --type blocks --target LIN-456
  linear-cli issue relation add LIN-123 --type parent --target LIN-100
  linear-cli issue relation remove LIN-123 --type blocks --target LIN-456
  linear-cli issue relation update RELATION-ID --type related`,
}

var relationListCmd = &cobra.Command{
	Use:     "list [issue-id]",
	Aliases: []string{"ls"},
	Short:   "List all relationships for an issue",
	Long:    `List all relationships for an issue including relations, parent, and sub-issues.`,
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
		issue, err := client.GetIssue(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		type relationEntry struct {
			Type       string    `json:"type"`
			Identifier string    `json:"identifier"`
			Title      string    `json:"title"`
			State      string    `json:"state"`
			RelationID string    `json:"relationId,omitempty"`
			IssueID    string    `json:"issueId"`
		}

		var entries []relationEntry

		// Parent
		if issue.Parent != nil {
			state := ""
			if issue.Parent.State != nil {
				state = issue.Parent.State.Name
			}
			entries = append(entries, relationEntry{
				Type:       "parent",
				Identifier: issue.Parent.Identifier,
				Title:      issue.Parent.Title,
				State:      state,
				IssueID:    issue.Parent.ID,
			})
		}

		// Children (sub-issues)
		if issue.Children != nil {
			for _, child := range issue.Children.Nodes {
				state := ""
				if child.State != nil {
					state = child.State.Name
				}
				entries = append(entries, relationEntry{
					Type:       "sub-issue",
					Identifier: child.Identifier,
					Title:      child.Title,
					State:      state,
					IssueID:    child.ID,
				})
			}
		}

		// Formal relations
		if issue.Relations != nil {
			for _, rel := range issue.Relations.Nodes {
				if rel.RelatedIssue != nil {
					state := ""
					if rel.RelatedIssue.State != nil {
						state = rel.RelatedIssue.State.Name
					}
					entries = append(entries, relationEntry{
						Type:       rel.Type,
						Identifier: rel.RelatedIssue.Identifier,
						Title:      rel.RelatedIssue.Title,
						State:      state,
						RelationID: rel.ID,
						IssueID:    rel.RelatedIssue.ID,
					})
				}
			}
		}

		if len(entries) == 0 {
			output.Info(fmt.Sprintf("No relationships found for %s", issue.Identifier), plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(entries)
			return
		}

		if plaintext {
			fmt.Printf("# Relationships for %s\n\n", issue.Identifier)
			for _, e := range entries {
				fmt.Printf("- **%s**: %s - %s", formatRelationType(e.Type), e.Identifier, e.Title)
				if e.State != "" {
					fmt.Printf(" [%s]", e.State)
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d relationships\n", len(entries))
			return
		}

		// Table output
		headers := []string{"Type", "Issue", "Title", "State"}
		rows := make([][]string, len(entries))
		for i, e := range entries {
			typeStr := formatRelationType(e.Type)
			switch e.Type {
			case "parent":
				typeStr = color.New(color.FgMagenta).Sprint(typeStr)
			case "sub-issue":
				typeStr = color.New(color.FgCyan).Sprint(typeStr)
			case "blocks":
				typeStr = color.New(color.FgRed).Sprint(typeStr)
			case "duplicate":
				typeStr = color.New(color.FgYellow).Sprint(typeStr)
			default:
				typeStr = color.New(color.FgBlue).Sprint(typeStr)
			}

			rows[i] = []string{
				typeStr,
				e.Identifier,
				truncateString(e.Title, 40),
				e.State,
			}
		}

		tableData := output.TableData{
			Headers: headers,
			Rows:    rows,
		}
		output.Table(tableData, false, false)

		fmt.Printf("\n%s %d relationships for %s\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(entries),
			color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier))
	},
}

var relationAddCmd = &cobra.Command{
	Use:     "add [issue-id]",
	Aliases: []string{"create", "new"},
	Short:   "Add a relationship between issues",
	Long: `Add a relationship between two issues.

Types:
  blocks      This issue blocks the target
  blocked-by  This issue is blocked by the target
  related     Issues are related (bidirectional)
  duplicate   This issue is a duplicate of the target
  parent      Set the target as this issue's parent
  sub-issue   Make the target a sub-issue of this issue

Examples:
  linear-cli issue relation add LIN-123 --type blocks --target LIN-456
  linear-cli issue relation add LIN-123 --type blocked-by --target LIN-456
  linear-cli issue relation add LIN-123 --type related --target LIN-456
  linear-cli issue relation add LIN-123 --type parent --target LIN-100
  linear-cli issue relation add LIN-123 --type sub-issue --target LIN-456`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		relType, _ := cmd.Flags().GetString("type")
		target, _ := cmd.Flags().GetString("target")

		if relType == "" {
			output.Error("--type is required", plaintext, jsonOut)
			os.Exit(1)
		}
		if target == "" {
			output.Error("--target is required", plaintext, jsonOut)
			os.Exit(1)
		}

		validTypes := []string{"blocks", "blocked-by", "related", "duplicate", "parent", "sub-issue"}
		found := false
		for _, vt := range validTypes {
			if relType == vt {
				found = true
				break
			}
		}
		if !found {
			output.Error(fmt.Sprintf("Invalid type '%s'. Valid types: %s", relType, strings.Join(validTypes, ", ")), plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		issueID := args[0]

		switch relType {
		case "parent":
			// Resolve the target to get its ID
			targetIssue, err := client.GetIssue(context.Background(), target)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve target issue '%s': %v", target, err), plaintext, jsonOut)
				os.Exit(1)
			}

			input := map[string]interface{}{
				"parentId": targetIssue.ID,
			}
			issue, err := client.UpdateIssue(context.Background(), issueID, input)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to set parent: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(issue)
			} else if plaintext {
				fmt.Printf("Set %s as parent of %s\n", target, issue.Identifier)
			} else {
				fmt.Printf("%s Set %s as parent of %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(target),
					color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier))
			}

		case "sub-issue":
			// Resolve issueID to get its real ID
			parentIssue, err := client.GetIssue(context.Background(), issueID)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve issue '%s': %v", issueID, err), plaintext, jsonOut)
				os.Exit(1)
			}

			input := map[string]interface{}{
				"parentId": parentIssue.ID,
			}
			childIssue, err := client.UpdateIssue(context.Background(), target, input)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to set sub-issue: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(childIssue)
			} else if plaintext {
				fmt.Printf("Made %s a sub-issue of %s\n", childIssue.Identifier, issueID)
			} else {
				fmt.Printf("%s Made %s a sub-issue of %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(childIssue.Identifier),
					color.New(color.FgCyan, color.Bold).Sprint(issueID))
			}

		case "blocked-by":
			// Swap direction: target blocks issueID
			targetIssue, err := client.GetIssue(context.Background(), target)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve target issue '%s': %v", target, err), plaintext, jsonOut)
				os.Exit(1)
			}
			srcIssue, err := client.GetIssue(context.Background(), issueID)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve issue '%s': %v", issueID, err), plaintext, jsonOut)
				os.Exit(1)
			}

			relation, err := client.CreateIssueRelation(context.Background(), targetIssue.ID, srcIssue.ID, "blocks")
			if err != nil {
				output.Error(fmt.Sprintf("Failed to create relation: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(relation)
			} else if plaintext {
				fmt.Printf("Added blocked-by relation: %s is blocked by %s\n", issueID, target)
			} else {
				fmt.Printf("%s %s is now blocked by %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(issueID),
					color.New(color.FgCyan, color.Bold).Sprint(target))
			}

		default:
			// blocks, related, duplicate — direct API call
			srcIssue, err := client.GetIssue(context.Background(), issueID)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve issue '%s': %v", issueID, err), plaintext, jsonOut)
				os.Exit(1)
			}
			targetIssue, err := client.GetIssue(context.Background(), target)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve target issue '%s': %v", target, err), plaintext, jsonOut)
				os.Exit(1)
			}

			relation, err := client.CreateIssueRelation(context.Background(), srcIssue.ID, targetIssue.ID, relType)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to create relation: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(relation)
			} else if plaintext {
				fmt.Printf("Added %s relation: %s %s %s\n", relType, issueID, relType, target)
			} else {
				fmt.Printf("%s Added %s relation: %s %s %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					formatRelationType(relType),
					color.New(color.FgCyan, color.Bold).Sprint(issueID),
					relType,
					color.New(color.FgCyan, color.Bold).Sprint(target))
			}
		}
	},
}

var relationRemoveCmd = &cobra.Command{
	Use:     "remove [issue-id]",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a relationship between issues",
	Long: `Remove a relationship between two issues.

For parent/sub-issue relationships, this unsets the parentId.
For other relations, it finds and deletes the matching relation.

Examples:
  linear-cli issue relation remove LIN-123 --type parent --target LIN-100
  linear-cli issue relation remove LIN-123 --type blocks --target LIN-456
  linear-cli issue relation remove LIN-123 --type sub-issue --target LIN-789`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		relType, _ := cmd.Flags().GetString("type")
		target, _ := cmd.Flags().GetString("target")

		if relType == "" {
			output.Error("--type is required", plaintext, jsonOut)
			os.Exit(1)
		}
		if target == "" {
			output.Error("--target is required", plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		issueID := args[0]

		switch relType {
		case "parent":
			// Unset parent on this issue
			input := map[string]interface{}{
				"parentId": nil,
			}
			issue, err := client.UpdateIssue(context.Background(), issueID, input)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to remove parent: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(issue)
			} else if plaintext {
				fmt.Printf("Removed parent from %s\n", issue.Identifier)
			} else {
				fmt.Printf("%s Removed parent from %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier))
			}

		case "sub-issue":
			// Unset parent on the target issue
			input := map[string]interface{}{
				"parentId": nil,
			}
			childIssue, err := client.UpdateIssue(context.Background(), target, input)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to remove sub-issue: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(childIssue)
			} else if plaintext {
				fmt.Printf("Removed %s as sub-issue of %s\n", childIssue.Identifier, issueID)
			} else {
				fmt.Printf("%s Removed %s as sub-issue of %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(childIssue.Identifier),
					color.New(color.FgCyan, color.Bold).Sprint(issueID))
			}

		default:
			// For blocks, blocked-by, related, duplicate: find and delete the relation
			issue, err := client.GetIssue(context.Background(), issueID)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			// Resolve target identifier
			targetIssue, err := client.GetIssue(context.Background(), target)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to resolve target issue '%s': %v", target, err), plaintext, jsonOut)
				os.Exit(1)
			}

			// Find matching relation
			var relationID string
			if issue.Relations != nil {
				for _, rel := range issue.Relations.Nodes {
					if rel.RelatedIssue != nil && rel.RelatedIssue.ID == targetIssue.ID {
						// Match type (blocked-by maps to blocks in the API)
						apiType := relType
						if apiType == "blocked-by" {
							apiType = "blocks"
						}
						if rel.Type == apiType || rel.Type == relType {
							relationID = rel.ID
							break
						}
					}
				}
			}

			if relationID == "" {
				output.Error(fmt.Sprintf("No %s relation found between %s and %s", relType, issueID, target), plaintext, jsonOut)
				os.Exit(1)
			}

			err = client.DeleteIssueRelation(context.Background(), relationID)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to delete relation: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}

			if jsonOut {
				output.JSON(map[string]interface{}{"success": true, "relationId": relationID})
			} else if plaintext {
				fmt.Printf("Removed %s relation between %s and %s\n", relType, issueID, target)
			} else {
				fmt.Printf("%s Removed %s relation between %s and %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					formatRelationType(relType),
					color.New(color.FgCyan, color.Bold).Sprint(issueID),
					color.New(color.FgCyan, color.Bold).Sprint(target))
			}
		}
	},
}

var relationUpdateCmd = &cobra.Command{
	Use:     "update [relation-id]",
	Aliases: []string{"edit"},
	Short:   "Update a relation's type",
	Long: `Update the type of an existing issue relation.

Only formal relations (blocks, related, duplicate) can be updated.
Parent/sub-issue relationships should be removed and re-added instead.

The relation ID can be found via 'issue relation list' with --json output.

Examples:
  linear-cli issue relation update RELATION-UUID --type related`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")

		newType, _ := cmd.Flags().GetString("type")
		if newType == "" {
			output.Error("--type is required", plaintext, jsonOut)
			os.Exit(1)
		}

		validTypes := []string{"blocks", "related", "duplicate"}
		found := false
		for _, vt := range validTypes {
			if newType == vt {
				found = true
				break
			}
		}
		if !found {
			output.Error(fmt.Sprintf("Invalid type '%s'. Valid types for update: %s", newType, strings.Join(validTypes, ", ")), plaintext, jsonOut)
			os.Exit(1)
		}

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		input := map[string]interface{}{
			"type": newType,
		}

		relation, err := client.UpdateIssueRelation(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update relation: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(relation)
		} else if plaintext {
			issueIdent := relation.ID
			if relation.Issue != nil {
				issueIdent = relation.Issue.Identifier
			}
			relatedIdent := ""
			if relation.RelatedIssue != nil {
				relatedIdent = relation.RelatedIssue.Identifier
			}
			fmt.Printf("Updated relation to %s: %s -> %s\n", newType, issueIdent, relatedIdent)
		} else {
			issueIdent := relation.ID
			if relation.Issue != nil {
				issueIdent = relation.Issue.Identifier
			}
			relatedIdent := ""
			if relation.RelatedIssue != nil {
				relatedIdent = relation.RelatedIssue.Identifier
			}
			fmt.Printf("%s Updated relation to %s: %s -> %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				formatRelationType(newType),
				color.New(color.FgCyan, color.Bold).Sprint(issueIdent),
				color.New(color.FgCyan, color.Bold).Sprint(relatedIdent))
		}
	},
}

func formatRelationType(t string) string {
	switch t {
	case "blocks":
		return "Blocks"
	case "blocked-by":
		return "Blocked by"
	case "related":
		return "Related to"
	case "duplicate":
		return "Duplicate of"
	case "parent":
		return "Parent"
	case "sub-issue":
		return "Sub-issue"
	default:
		return t
	}
}

func init() {
	issueCmd.AddCommand(issueRelationCmd)
	issueRelationCmd.AddCommand(relationListCmd)
	issueRelationCmd.AddCommand(relationAddCmd)
	issueRelationCmd.AddCommand(relationRemoveCmd)
	issueRelationCmd.AddCommand(relationUpdateCmd)

	// relation add flags
	relationAddCmd.Flags().String("type", "", "Relation type: blocks, blocked-by, related, duplicate, parent, sub-issue")
	relationAddCmd.Flags().String("target", "", "Target issue identifier (e.g., LIN-456)")
	_ = relationAddCmd.MarkFlagRequired("type")
	_ = relationAddCmd.MarkFlagRequired("target")

	// relation remove flags
	relationRemoveCmd.Flags().String("type", "", "Relation type: blocks, blocked-by, related, duplicate, parent, sub-issue")
	relationRemoveCmd.Flags().String("target", "", "Target issue identifier (e.g., LIN-456)")
	_ = relationRemoveCmd.MarkFlagRequired("type")
	_ = relationRemoveCmd.MarkFlagRequired("target")

	// relation update flags
	relationUpdateCmd.Flags().String("type", "", "New relation type: blocks, related, duplicate")
	_ = relationUpdateCmd.MarkFlagRequired("type")
}
