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

// treeNode represents one node in the dependency tree (for JSON output)
type treeNode struct {
	Type       string     `json:"type"`
	ID         string     `json:"id"`
	Identifier string     `json:"identifier"`
	Title      string     `json:"title"`
	State      string     `json:"state,omitempty"`
	Children   []treeNode `json:"children,omitempty"`
}

// treeWalker holds state during recursive tree traversal
type treeWalker struct {
	client   *api.Client
	ctx      context.Context
	visited  map[string]bool
	maxDepth int
}

var issueTreeCmd = &cobra.Command{
	Use:     "tree [issue-id]",
	Aliases: []string{"deps"},
	Short:   "Display issue dependency tree",
	Long: `Display the full dependency tree of an issue, showing parent, children,
blocking/blocked-by, and related issues in a hierarchical view.

Recursively walks all relations up to the specified depth, handling
circular references gracefully.

Examples:
  linear-cli issue tree LIN-123
  linear-cli issue tree LIN-123 --depth 2
  linear-cli issue tree LIN-123 --json
  linear-cli issue tree LIN-123 --depth 5 --json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut, _ := cmd.Flags().GetBool("json")
		if !jsonOut {
			jsonOut = viper.GetBool("json")
		}
		depth, _ := cmd.Flags().GetInt("depth")

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error("Not authenticated. Run 'linear-cli auth' first.", plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)
		ctx := context.Background()

		issue, err := client.GetIssue(ctx, args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		walker := &treeWalker{
			client:   client,
			ctx:      ctx,
			visited:  make(map[string]bool),
			maxDepth: depth,
		}

		if jsonOut {
			node := walker.buildTreeNode(issue, "root", 0)
			output.JSON(node)
			return
		}

		// Text output
		walker.printTree(issue, plaintext)
	},
}

// issueStateName returns the state name for an issue, or empty string
func issueStateName(issue *api.Issue) string {
	if issue == nil || issue.State == nil {
		return ""
	}
	return issue.State.Name
}

// formatIssueOneLiner formats an issue as "IDENT (State) - Title"
func formatIssueOneLiner(issue *api.Issue, plaintext bool) string {
	state := issueStateName(issue)
	if plaintext {
		if state != "" {
			return fmt.Sprintf("%s (%s) - %s", issue.Identifier, state, issue.Title)
		}
		return fmt.Sprintf("%s - %s", issue.Identifier, issue.Title)
	}
	ident := color.New(color.FgCyan, color.Bold).Sprint(issue.Identifier)
	stateStr := ""
	if state != "" {
		stateStr = fmt.Sprintf(" (%s)", colorizeState(issue.State))
	}
	return fmt.Sprintf("%s%s - %s", ident, stateStr, issue.Title)
}

// colorizeState returns a colorized state name
func colorizeState(state *api.State) string {
	if state == nil {
		return ""
	}
	var stateColor *color.Color
	switch state.Type {
	case "triage":
		stateColor = color.New(color.FgMagenta)
	case "backlog":
		stateColor = color.New(color.FgWhite, color.Faint)
	case "unstarted":
		stateColor = color.New(color.FgWhite)
	case "started":
		stateColor = color.New(color.FgYellow)
	case "completed":
		stateColor = color.New(color.FgGreen)
	case "canceled":
		stateColor = color.New(color.FgRed)
	default:
		stateColor = color.New(color.FgWhite)
	}
	return stateColor.Sprint(state.Name)
}

// printTree prints the text representation of the issue tree
func (w *treeWalker) printTree(issue *api.Issue, plaintext bool) {
	// Mark root as visited
	w.visited[issue.ID] = true

	// Print root issue
	fmt.Println(formatIssueOneLiner(issue, plaintext))

	// Collect the sections we need to print
	type section struct {
		label string
		items []*api.Issue
	}

	var sections []section

	// Parent
	if issue.Parent != nil {
		sections = append(sections, section{label: "parent", items: []*api.Issue{issue.Parent}})
	}

	// Sub-issues (children)
	if issue.Children != nil && len(issue.Children.Nodes) > 0 {
		children := make([]*api.Issue, len(issue.Children.Nodes))
		for i := range issue.Children.Nodes {
			children[i] = &issue.Children.Nodes[i]
		}
		sections = append(sections, section{label: "sub-issues", items: children})
	}

	// Relations grouped by type
	blocks, blockedBy, related, duplicates := groupRelations(issue)
	if len(blocks) > 0 {
		sections = append(sections, section{label: "blocks", items: blocks})
	}
	if len(blockedBy) > 0 {
		sections = append(sections, section{label: "blocked-by", items: blockedBy})
	}
	if len(related) > 0 {
		sections = append(sections, section{label: "related", items: related})
	}
	if len(duplicates) > 0 {
		sections = append(sections, section{label: "duplicates", items: duplicates})
	}

	// Print each section
	for i, sec := range sections {
		isLastSection := i == len(sections)-1
		w.printSection(sec.label, sec.items, "", isLastSection, plaintext, 1)
	}
}

// printSection prints a labeled section of the tree
func (w *treeWalker) printSection(label string, items []*api.Issue, prefix string, isLast bool, plaintext bool, depth int) {
	// Determine connectors
	connector := "\u251c\u2500\u2500" // ├──
	if isLast {
		connector = "\u2514\u2500\u2500" // └──
	}

	childPrefix := prefix + "\u2502   " // │   (for continuing lines)
	if isLast {
		childPrefix = prefix + "    "
	}

	if len(items) == 1 && (label == "parent") {
		// Single-item sections print inline
		item := items[0]
		circularMark := ""
		if w.visited[item.ID] {
			circularMark = " [circular]"
		}
		if plaintext {
			fmt.Printf("%s%s %s: %s%s\n", prefix, connector, label, formatIssueOneLiner(item, plaintext), circularMark)
		} else {
			labelStr := color.New(color.FgMagenta).Sprint(label)
			fmt.Printf("%s%s %s: %s%s\n", prefix, connector, labelStr, formatIssueOneLiner(item, plaintext), circularMark)
		}

		// Recurse into parent if not visited and within depth
		if !w.visited[item.ID] && depth < w.maxDepth {
			w.visited[item.ID] = true
			fullIssue, err := w.client.GetIssue(w.ctx, item.Identifier)
			if err == nil {
				w.printSubTree(fullIssue, childPrefix, plaintext, depth+1)
			}
		}
		return
	}

	// Multi-item section: print label, then items underneath
	if plaintext {
		fmt.Printf("%s%s %s:\n", prefix, connector, label)
	} else {
		labelStr := color.New(color.FgMagenta).Sprint(label + ":")
		fmt.Printf("%s%s %s\n", prefix, connector, labelStr)
	}

	for j, item := range items {
		isLastItem := j == len(items)-1
		itemConnector := "\u251c\u2500\u2500" // ├──
		if isLastItem {
			itemConnector = "\u2514\u2500\u2500" // └──
		}

		circularMark := ""
		if w.visited[item.ID] {
			circularMark = " [circular]"
		}

		fmt.Printf("%s%s %s%s\n", childPrefix, itemConnector, formatIssueOneLiner(item, plaintext), circularMark)

		// Recurse if not visited and within depth
		if !w.visited[item.ID] && depth < w.maxDepth {
			w.visited[item.ID] = true
			fullIssue, err := w.client.GetIssue(w.ctx, item.Identifier)
			if err == nil {
				subPrefix := childPrefix + "\u2502   "
				if isLastItem {
					subPrefix = childPrefix + "    "
				}
				w.printSubTree(fullIssue, subPrefix, plaintext, depth+1)
			}
		}
	}
}

// printSubTree prints the sub-tree of relations for a recursively fetched issue
func (w *treeWalker) printSubTree(issue *api.Issue, prefix string, plaintext bool, depth int) {
	type section struct {
		label string
		items []*api.Issue
	}

	var sections []section

	// Parent (only if not already visited)
	if issue.Parent != nil && !w.visited[issue.Parent.ID] {
		sections = append(sections, section{label: "parent", items: []*api.Issue{issue.Parent}})
	}

	// Sub-issues
	if issue.Children != nil && len(issue.Children.Nodes) > 0 {
		children := make([]*api.Issue, 0)
		for i := range issue.Children.Nodes {
			children = append(children, &issue.Children.Nodes[i])
		}
		sections = append(sections, section{label: "sub-issues", items: children})
	}

	// Relations
	blocks, blockedBy, related, duplicates := groupRelations(issue)
	if len(blocks) > 0 {
		sections = append(sections, section{label: "blocks", items: blocks})
	}
	if len(blockedBy) > 0 {
		sections = append(sections, section{label: "blocked-by", items: blockedBy})
	}
	if len(related) > 0 {
		sections = append(sections, section{label: "related", items: related})
	}
	if len(duplicates) > 0 {
		sections = append(sections, section{label: "duplicates", items: duplicates})
	}

	for i, sec := range sections {
		isLastSection := i == len(sections)-1
		w.printSection(sec.label, sec.items, prefix, isLastSection, plaintext, depth)
	}
}

// buildTreeNode recursively builds a JSON tree node
func (w *treeWalker) buildTreeNode(issue *api.Issue, relType string, depth int) treeNode {
	w.visited[issue.ID] = true

	node := treeNode{
		Type:       relType,
		ID:         issue.ID,
		Identifier: issue.Identifier,
		Title:      issue.Title,
		State:      issueStateName(issue),
	}

	if depth >= w.maxDepth {
		return node
	}

	// Parent
	if issue.Parent != nil {
		if w.visited[issue.Parent.ID] {
			node.Children = append(node.Children, treeNode{
				Type:       "parent",
				ID:         issue.Parent.ID,
				Identifier: issue.Parent.Identifier,
				Title:      issue.Parent.Title,
				State:      issueStateName(issue.Parent),
			})
		} else {
			fullParent, err := w.client.GetIssue(w.ctx, issue.Parent.Identifier)
			if err == nil {
				node.Children = append(node.Children, w.buildTreeNode(fullParent, "parent", depth+1))
			} else {
				node.Children = append(node.Children, treeNode{
					Type:       "parent",
					ID:         issue.Parent.ID,
					Identifier: issue.Parent.Identifier,
					Title:      issue.Parent.Title,
					State:      issueStateName(issue.Parent),
				})
			}
		}
	}

	// Sub-issues
	if issue.Children != nil {
		for i := range issue.Children.Nodes {
			child := &issue.Children.Nodes[i]
			if w.visited[child.ID] {
				node.Children = append(node.Children, treeNode{
					Type:       "sub-issue",
					ID:         child.ID,
					Identifier: child.Identifier,
					Title:      child.Title,
					State:      issueStateName(child),
				})
			} else {
				fullChild, err := w.client.GetIssue(w.ctx, child.Identifier)
				if err == nil {
					node.Children = append(node.Children, w.buildTreeNode(fullChild, "sub-issue", depth+1))
				} else {
					node.Children = append(node.Children, treeNode{
						Type:       "sub-issue",
						ID:         child.ID,
						Identifier: child.Identifier,
						Title:      child.Title,
						State:      issueStateName(child),
					})
				}
			}
		}
	}

	// Relations
	if issue.Relations != nil {
		for _, rel := range issue.Relations.Nodes {
			if rel.RelatedIssue == nil {
				continue
			}
			relName := normalizeRelationType(rel.Type)
			ri := rel.RelatedIssue
			if w.visited[ri.ID] {
				node.Children = append(node.Children, treeNode{
					Type:       relName,
					ID:         ri.ID,
					Identifier: ri.Identifier,
					Title:      ri.Title,
					State:      issueStateName(ri),
				})
			} else {
				fullRel, err := w.client.GetIssue(w.ctx, ri.Identifier)
				if err == nil {
					node.Children = append(node.Children, w.buildTreeNode(fullRel, relName, depth+1))
				} else {
					node.Children = append(node.Children, treeNode{
						Type:       relName,
						ID:         ri.ID,
						Identifier: ri.Identifier,
						Title:      ri.Title,
						State:      issueStateName(ri),
					})
				}
			}
		}
	}

	return node
}

// groupRelations splits issue relations into typed groups
func groupRelations(issue *api.Issue) (blocks, blockedBy, related, duplicates []*api.Issue) {
	if issue.Relations == nil {
		return
	}
	for i := range issue.Relations.Nodes {
		rel := &issue.Relations.Nodes[i]
		if rel.RelatedIssue == nil {
			continue
		}
		switch strings.ToLower(rel.Type) {
		case "blocks":
			blocks = append(blocks, rel.RelatedIssue)
		case "blocked-by", "blocked":
			blockedBy = append(blockedBy, rel.RelatedIssue)
		case "related":
			related = append(related, rel.RelatedIssue)
		case "duplicate":
			duplicates = append(duplicates, rel.RelatedIssue)
		default:
			related = append(related, rel.RelatedIssue)
		}
	}
	return
}

// normalizeRelationType maps API relation types to display names
func normalizeRelationType(t string) string {
	switch strings.ToLower(t) {
	case "blocks":
		return "blocks"
	case "blocked-by", "blocked":
		return "blocked-by"
	case "related":
		return "related"
	case "duplicate":
		return "duplicate"
	default:
		return t
	}
}

func init() {
	issueCmd.AddCommand(issueTreeCmd)

	issueTreeCmd.Flags().Int("depth", 3, "Maximum recursion depth")
	issueTreeCmd.Flags().Bool("json", false, "Output as JSON tree")
}
