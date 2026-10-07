package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var roadmapCmd = &cobra.Command{
	Use:   "roadmap",
	Short: "Manage Linear roadmaps",
	Long: `Manage Linear roadmaps (deprecated in favor of initiatives, but still readable).

Examples:
  linear-cli roadmap list              # List roadmaps
  linear-cli roadmap list -l 10        # Limit results
  linear-cli roadmap get ROADMAP-ID    # Get roadmap details with linked projects`,
}

var roadmapListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List roadmaps",
	Long:    `List all roadmaps in the workspace.`,
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

		roadmaps, err := client.GetRoadmaps(context.Background(), limit, "", "createdAt", false)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list roadmaps: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(roadmaps.Nodes)
			return
		}

		if len(roadmaps.Nodes) == 0 {
			if plaintext {
				fmt.Println("No roadmaps found")
			} else {
				fmt.Printf("\n%s No roadmaps found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Roadmaps")
			fmt.Println("Name\tOwner\tSlug\tCreated")
			for _, r := range roadmaps.Nodes {
				owner := ""
				if r.Owner != nil {
					owner = r.Owner.Name
				}
				fmt.Printf("%s\t%s\t%s\t%s\n", r.Name, owner, r.SlugId, r.CreatedAt.Format("2006-01-02"))
			}
		} else {
			headers := []string{"Name", "Owner", "Slug", "Created"}
			rows := [][]string{}

			for _, r := range roadmaps.Nodes {
				owner := color.New(color.FgYellow).Sprint("Unassigned")
				if r.Owner != nil {
					owner = r.Owner.Name
				}
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(truncateString(r.Name, 30)),
					owner,
					r.SlugId,
					r.CreatedAt.Format("2006-01-02"),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)

			fmt.Printf("\n%s %d roadmaps\n",
				color.New(color.FgGreen).Sprint("✓"),
				len(roadmaps.Nodes))
		}
	},
}

var roadmapGetCmd = &cobra.Command{
	Use:     "get ROADMAP-ID",
	Aliases: []string{"show"},
	Short:   "Get roadmap details",
	Long:    `Get details for a specific roadmap, including its linked projects.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		plaintext := viper.GetBool("plaintext")
		jsonOut := viper.GetBool("json")
		roadmapID := args[0]

		authHeader, err := auth.GetAuthHeader()
		if err != nil {
			output.Error(fmt.Sprintf("Authentication failed: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		client := api.NewClient(authHeader)

		roadmap, err := client.GetRoadmap(context.Background(), roadmapID)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get roadmap: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(roadmap)
			return
		}

		owner := "Unassigned"
		if roadmap.Owner != nil {
			owner = roadmap.Owner.Name
		}

		if plaintext {
			fmt.Printf("# %s\n\n", roadmap.Name)
			if roadmap.Description != "" {
				fmt.Printf("## Description\n%s\n\n", roadmap.Description)
			}
			fmt.Printf("## Core Details\n")
			fmt.Printf("- **ID**: %s\n", roadmap.ID)
			fmt.Printf("- **Slug**: %s\n", roadmap.SlugId)
			fmt.Printf("- **Owner**: %s\n", owner)
			if roadmap.Color != "" {
				fmt.Printf("- **Color**: %s\n", roadmap.Color)
			}
			fmt.Printf("- **Created**: %s\n", roadmap.CreatedAt.Format("2006-01-02"))
			fmt.Printf("- **Updated**: %s\n", roadmap.UpdatedAt.Format("2006-01-02"))
			if roadmap.ArchivedAt != nil {
				fmt.Printf("- **Archived**: %s\n", roadmap.ArchivedAt.Format("2006-01-02"))
			}
			if roadmap.URL != "" {
				fmt.Printf("- **URL**: %s\n", roadmap.URL)
			}
			if roadmap.Projects != nil && len(roadmap.Projects.Nodes) > 0 {
				fmt.Printf("\n## Linked Projects\n")
				for _, p := range roadmap.Projects.Nodes {
					lead := ""
					if p.Lead != nil {
						lead = p.Lead.Name
					}
					fmt.Printf("- %s [%s] %.0f%%", p.Name, p.State, p.Progress*100)
					if lead != "" {
						fmt.Printf(" (lead: %s)", lead)
					}
					fmt.Println()
				}
			}
			return
		}

		fmt.Printf("\n%s %s\n",
			color.New(color.FgCyan, color.Bold).Sprint("🗺️ Roadmap:"),
			color.New(color.FgWhite, color.Bold).Sprint(roadmap.Name))
		if roadmap.Description != "" {
			fmt.Printf("   %s\n", roadmap.Description)
		}
		fmt.Printf("   Owner: %s\n", owner)
		fmt.Printf("   Slug: %s\n", color.New(color.FgCyan).Sprint(roadmap.SlugId))
		fmt.Printf("   Created: %s | Updated: %s\n",
			roadmap.CreatedAt.Format("2006-01-02"),
			roadmap.UpdatedAt.Format("2006-01-02"))
		if roadmap.URL != "" {
			fmt.Printf("   URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(roadmap.URL))
		}

		if roadmap.Projects != nil && len(roadmap.Projects.Nodes) > 0 {
			fmt.Printf("\n   %s Linked Projects:\n", color.New(color.FgCyan, color.Bold).Sprint("📁"))
			headers := []string{"Name", "State", "Progress", "Lead"}
			rows := [][]string{}
			for _, p := range roadmap.Projects.Nodes {
				lead := "Unassigned"
				if p.Lead != nil {
					lead = p.Lead.Name
				}
				rows = append(rows, []string{
					p.Name,
					p.State,
					fmt.Sprintf("%.0f%%", p.Progress*100),
					lead,
				})
			}
			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(roadmapCmd)
	roadmapCmd.AddCommand(roadmapListCmd)
	roadmapCmd.AddCommand(roadmapGetCmd)

	roadmapListCmd.Flags().IntP("limit", "l", 25, "Maximum number of roadmaps to return")
}
