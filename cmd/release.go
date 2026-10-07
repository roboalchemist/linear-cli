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

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage Linear releases",
	Long: `Manage Linear releases, release notes, pipelines, and stages.

Examples:
  linear-cli release list                          # List releases
  linear-cli release list --include-archived       # Include archived releases
  linear-cli release get RELEASE-ID                # Get release details
  linear-cli release note list                     # List release notes
  linear-cli release pipeline list                 # List release pipelines
  linear-cli release stage list                    # List release stages`,
}

var releaseListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List releases",
	Long:    `List Linear releases, optionally including archived releases.`,
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

		releases, err := client.GetReleases(context.Background(), nil, limit, "", includeArchived)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list releases: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(releases.Nodes)
			return
		}

		if len(releases.Nodes) == 0 {
			if plaintext {
				fmt.Println("No releases found")
			} else {
				fmt.Printf("\n%s No releases found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Releases")
			fmt.Println("Name\tVersion\tStage\tPipeline\tTarget\tProgress\tStart\tRelease")
			for _, r := range releases.Nodes {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					r.Name, releaseVersion(r), releaseStageName(r), releasePipelineName(r),
					releaseDateStr(r.TargetDate), formatReleaseProgress(r.CurrentProgress),
					releaseDateStr(r.StartDate), releaseCompletedShort(r))
			}
		} else {
			headers := []string{"Name", "Version", "Stage", "Pipeline", "Target", "Progress", "Start", "Release"}
			rows := [][]string{}

			for _, r := range releases.Nodes {
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(r.Name),
					releaseVersion(r),
					releaseStageName(r),
					releasePipelineName(r),
					releaseDateStr(r.TargetDate),
					formatReleaseProgress(r.CurrentProgress),
					releaseDateStr(r.StartDate),
					releaseCompletedShort(r),
				})
			}

			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
	},
}

var releaseGetCmd = &cobra.Command{
	Use:     "get RELEASE-ID",
	Aliases: []string{"show"},
	Short:   "Get release details",
	Long:    `Get details for a specific release by ID or slug.`,
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
		release, err := client.GetRelease(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to get release: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(release)
			return
		}

		if plaintext {
			fmt.Printf("# Release: %s\n", release.Name)
			fmt.Printf("ID: %s\n", release.ID)
			fmt.Printf("Version: %s\n", releaseVersion(*release))
			fmt.Printf("Stage: %s\n", releaseStageName(*release))
			fmt.Printf("Pipeline: %s\n", releasePipelineName(*release))
			if release.Description != nil && *release.Description != "" {
				fmt.Printf("Description: %s\n", *release.Description)
			}
			fmt.Printf("Start: %s\n", releaseDateStr(release.StartDate))
			fmt.Printf("Target: %s\n", releaseDateStr(release.TargetDate))
			if release.CompletedAt != nil {
				fmt.Printf("Released: %s\n", release.CompletedAt.Format("2006-01-02"))
			}
			if release.StartedAt != nil {
				fmt.Printf("Started: %s\n", release.StartedAt.Format("2006-01-02"))
			}
			if release.CanceledAt != nil {
				fmt.Printf("Canceled: %s\n", release.CanceledAt.Format("2006-01-02"))
			}
			fmt.Printf("Progress: %s\n", formatReleaseProgress(release.CurrentProgress))
			fmt.Printf("URL: %s\n", release.URL)
			fmt.Printf("Created: %s\n", release.CreatedAt.Format("2006-01-02"))
			fmt.Printf("Updated: %s\n", release.UpdatedAt.Format("2006-01-02"))
			if release.ArchivedAt != nil {
				fmt.Printf("Archived: %s\n", release.ArchivedAt.Format("2006-01-02"))
			}
		} else {
			fmt.Printf("\n%s Release: %s\n",
				color.New(color.FgCyan, color.Bold).Sprint("🚀"),
				color.New(color.FgWhite, color.Bold).Sprint(release.Name))
			fmt.Printf("   ID: %s\n", color.New(color.FgWhite, color.Faint).Sprint(release.ID))
			fmt.Printf("   Version: %s\n", releaseVersion(*release))
			fmt.Printf("   Stage: %s\n", color.New(color.FgCyan).Sprint(releaseStageName(*release)))
			fmt.Printf("   Pipeline: %s\n", releasePipelineName(*release))
			if release.Description != nil && *release.Description != "" {
				fmt.Printf("   %s\n", *release.Description)
			}
			fmt.Printf("   Start: %s | Target: %s\n",
				releaseDateStr(release.StartDate), releaseDateStr(release.TargetDate))
			if release.CompletedAt != nil {
				fmt.Printf("   Released: %s\n", color.New(color.FgGreen).Sprint(release.CompletedAt.Format("2006-01-02")))
			}
			fmt.Printf("   Progress: %s\n", color.New(color.FgGreen).Sprint(formatReleaseProgress(release.CurrentProgress)))
			fmt.Printf("   URL: %s\n", release.URL)
			fmt.Printf("   Created: %s | Updated: %s\n",
				release.CreatedAt.Format("2006-01-02"),
				release.UpdatedAt.Format("2006-01-02"))
			if release.ArchivedAt != nil {
				fmt.Printf("   Archived: %s\n", release.ArchivedAt.Format("2006-01-02"))
			}
			if len(release.ReleaseNotes) > 0 {
				fmt.Printf("\n   %s Release notes:\n", color.New(color.FgCyan, color.Bold).Sprint("📝"))
				for _, n := range release.ReleaseNotes {
					title := "(untitled)"
					if n.Title != nil && *n.Title != "" {
						title = *n.Title
					}
					fmt.Printf("   - %s (%s)\n", title, n.CreatedAt.Format("2006-01-02"))
				}
			}
		}
	},
}

var releaseNoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage release notes",
	Long:  `Manage Linear release notes.`,
}

var releaseNoteListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List release notes",
	Long:    `List release notes in the workspace.`,
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

		notes, err := client.GetReleaseNotes(context.Background(), nil, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list release notes: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(notes.Nodes)
			return
		}

		if len(notes.Nodes) == 0 {
			if plaintext {
				fmt.Println("No release notes found")
			} else {
				fmt.Printf("\n%s No release notes found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Release Notes")
			fmt.Println("Title\tRelease\tCreated")
			for _, n := range notes.Nodes {
				fmt.Printf("%s\t%s\t%s\n", releaseNoteTitle(n), releaseNoteRelease(n), n.CreatedAt.Format("2006-01-02"))
			}
		} else {
			headers := []string{"Title", "Release", "Created"}
			rows := [][]string{}
			for _, n := range notes.Nodes {
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(releaseNoteTitle(n)),
					releaseNoteRelease(n),
					n.CreatedAt.Format("2006-01-02"),
				})
			}
			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
	},
}

var releasePipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage release pipelines",
	Long:  `Manage Linear release pipelines.`,
}

var releasePipelineListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List release pipelines",
	Long:    `List release pipelines in the workspace.`,
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

		pipelines, err := client.GetReleasePipelines(context.Background(), nil, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list release pipelines: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(pipelines.Nodes)
			return
		}

		if len(pipelines.Nodes) == 0 {
			if plaintext {
				fmt.Println("No release pipelines found")
			} else {
				fmt.Printf("\n%s No release pipelines found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Release Pipelines")
			fmt.Println("Name\tType\tCreated")
			for _, p := range pipelines.Nodes {
				fmt.Printf("%s\t%s\t%s\n", p.Name, p.Type, p.CreatedAt.Format("2006-01-02"))
			}
		} else {
			headers := []string{"Name", "Type", "Created"}
			rows := [][]string{}
			for _, p := range pipelines.Nodes {
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(p.Name),
					p.Type,
					p.CreatedAt.Format("2006-01-02"),
				})
			}
			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
	},
}

var releaseStageCmd = &cobra.Command{
	Use:   "stage",
	Short: "Manage release stages",
	Long:  `Manage Linear release stages.`,
}

var releaseStageListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List release stages",
	Long:    `List release stages in the workspace.`,
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

		stages, err := client.GetReleaseStages(context.Background(), nil, limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to list release stages: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(stages.Nodes)
			return
		}

		if len(stages.Nodes) == 0 {
			if plaintext {
				fmt.Println("No release stages found")
			} else {
				fmt.Printf("\n%s No release stages found\n", color.New(color.FgYellow).Sprint("ℹ️"))
			}
			return
		}

		if plaintext {
			fmt.Println("# Release Stages")
			fmt.Println("Name\tType\tPosition\tColor")
			for _, s := range stages.Nodes {
				fmt.Printf("%s\t%s\t%.0f\t%s\n", s.Name, s.Type, s.Position, s.Color)
			}
		} else {
			headers := []string{"Name", "Type", "Position", "Color"}
			rows := [][]string{}
			for _, s := range stages.Nodes {
				rows = append(rows, []string{
					color.New(color.FgWhite, color.Bold).Sprint(s.Name),
					s.Type,
					fmt.Sprintf("%.0f", s.Position),
					s.Color,
				})
			}
			output.Table(output.TableData{
				Headers: headers,
				Rows:    rows,
			}, plaintext, jsonOut)
		}
	},
}

// releaseVersion returns the release version or "—" when unset.
func releaseVersion(r api.ReleaseDetail) string {
	if r.Version != nil && *r.Version != "" {
		return *r.Version
	}
	return "—"
}

// releaseStageName returns the release's stage name, falling back to its type.
func releaseStageName(r api.ReleaseDetail) string {
	if r.Stage != nil {
		if r.Stage.Name != "" {
			return r.Stage.Name
		}
		return r.Stage.Type
	}
	return ""
}

// releasePipelineName returns the release's pipeline name.
func releasePipelineName(r api.ReleaseDetail) string {
	if r.Pipeline != nil {
		return r.Pipeline.Name
	}
	return ""
}

// releaseDateStr renders a TimelessDate value, returning "—" when unset.
func releaseDateStr(d *string) string {
	if d != nil && *d != "" {
		return formatDateShort(*d)
	}
	return "—"
}

// releaseCompletedShort renders the completed date, returning "—" when unset.
func releaseCompletedShort(r api.ReleaseDetail) string {
	if r.CompletedAt != nil {
		return r.CompletedAt.Format("2006-01-02")
	}
	return "—"
}

// releaseNoteTitle returns the note title or "(untitled)".
func releaseNoteTitle(n api.ReleaseNote) string {
	if n.Title != nil && *n.Title != "" {
		return *n.Title
	}
	return "(untitled)"
}

// releaseNoteRelease returns a human-readable label for the releases a note covers.
func releaseNoteRelease(n api.ReleaseNote) string {
	r := n.LastRelease
	if r == nil {
		r = n.FirstRelease
	}
	if r == nil {
		return ""
	}
	if r.Version != nil && *r.Version != "" {
		return *r.Version
	}
	return r.Name
}

// formatReleaseProgress renders a progress summary from a release's
// currentProgress JSON object (e.g. "1/22 (5%)").
func formatReleaseProgress(progress map[string]interface{}) string {
	if progress == nil {
		return "—"
	}
	total, ok := jsonNumber(progress["scopeCount"])
	if !ok || total == 0 {
		return "—"
	}
	completed, _ := jsonNumber(progress["completedIssueCount"])
	pct := completed / total * 100
	return fmt.Sprintf("%.0f/%.0f (%.0f%%)", completed, total, pct)
}

// jsonNumber coerces a decoded JSON numeric value into a float64.
func jsonNumber(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.AddCommand(releaseListCmd)
	releaseCmd.AddCommand(releaseGetCmd)
	releaseCmd.AddCommand(releaseNoteCmd)
	releaseNoteCmd.AddCommand(releaseNoteListCmd)
	releaseCmd.AddCommand(releasePipelineCmd)
	releasePipelineCmd.AddCommand(releasePipelineListCmd)
	releaseCmd.AddCommand(releaseStageCmd)
	releaseStageCmd.AddCommand(releaseStageListCmd)

	releaseListCmd.Flags().IntP("limit", "l", 25, "Maximum number of releases to return")
	releaseListCmd.Flags().Bool("include-archived", false, "Include archived releases")

	releaseNoteListCmd.Flags().IntP("limit", "l", 25, "Maximum number of release notes to return")
	releasePipelineListCmd.Flags().IntP("limit", "l", 25, "Maximum number of release pipelines to return")
	releaseStageListCmd.Flags().IntP("limit", "l", 25, "Maximum number of release stages to return")
}
