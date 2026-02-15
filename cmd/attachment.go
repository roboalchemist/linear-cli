package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/roboalchemist/linear-cli/pkg/api"
	"github.com/roboalchemist/linear-cli/pkg/auth"
	"github.com/roboalchemist/linear-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage issue attachments",
	Long: `Manage attachments (linked resources) on Linear issues.

Examples:
  linear-cli issue attachment list LIN-123
  linear-cli issue attachment create LIN-123 --url "https://example.com" --title "Spec Doc"
  linear-cli issue attachment link LIN-123 --url "https://github.com/org/repo/pull/42"
  linear-cli issue attachment update ATTACHMENT-ID --title "New Title"
  linear-cli issue attachment delete ATTACHMENT-ID`,
}

var attachmentListCmd = &cobra.Command{
	Use:     "list [issue-id]",
	Aliases: []string{"ls"},
	Short:   "List attachments for an issue",
	Long:    `List all attachments linked to a specific issue.`,
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
		limit, _ := cmd.Flags().GetInt("limit")

		attachments, err := client.GetIssueAttachments(context.Background(), args[0], limit, "")
		if err != nil {
			output.Error(fmt.Sprintf("Failed to fetch attachments: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if len(attachments.Nodes) == 0 {
			output.Info("No attachments found", plaintext, jsonOut)
			return
		}

		if jsonOut {
			output.JSON(attachments.Nodes)
			return
		}

		if plaintext {
			fmt.Printf("# Attachments for %s\n", args[0])
			for _, att := range attachments.Nodes {
				fmt.Printf("## %s\n", att.Title)
				fmt.Printf("- **ID**: %s\n", att.ID)
				fmt.Printf("- **URL**: %s\n", att.URL)
				if att.Subtitle != nil && *att.Subtitle != "" {
					fmt.Printf("- **Subtitle**: %s\n", *att.Subtitle)
				}
				if att.SourceType != nil && *att.SourceType != "" {
					fmt.Printf("- **Source Type**: %s\n", *att.SourceType)
				}
				if att.Creator != nil {
					fmt.Printf("- **Creator**: %s\n", att.Creator.Name)
				} else if att.ExternalUserCreator != nil {
					fmt.Printf("- **Creator**: %s (external)\n", att.ExternalUserCreator.Name)
				}
				fmt.Printf("- **Created**: %s\n", att.CreatedAt.Format("2006-01-02 15:04"))
				fmt.Printf("- **Updated**: %s\n", att.UpdatedAt.Format("2006-01-02 15:04"))
				if att.ArchivedAt != nil {
					fmt.Printf("- **Archived**: %s\n", att.ArchivedAt.Format("2006-01-02 15:04"))
				}
				fmt.Println()
			}
			fmt.Printf("\nTotal: %d attachments\n", len(attachments.Nodes))
			return
		}

		headers := []string{"Title", "URL", "Source", "Creator", "Created"}
		rows := make([][]string, len(attachments.Nodes))

		for i, att := range attachments.Nodes {
			creator := ""
			if att.Creator != nil {
				creator = att.Creator.Name
			} else if att.ExternalUserCreator != nil {
				creator = att.ExternalUserCreator.Name + " (ext)"
			}
			sourceType := ""
			if att.SourceType != nil {
				sourceType = *att.SourceType
			}
			rows[i] = []string{
				truncateString(att.Title, 35),
				truncateString(att.URL, 45),
				sourceType,
				creator,
				att.CreatedAt.Format("2006-01-02"),
			}
		}

		output.Table(output.TableData{
			Headers: headers,
			Rows:    rows,
		}, false, false)

		fmt.Printf("\n%s %d attachments\n",
			color.New(color.FgGreen).Sprint("✓"),
			len(attachments.Nodes))
	},
}

var attachmentCreateCmd = &cobra.Command{
	Use:     "create [issue-id]",
	Aliases: []string{"new"},
	Short:   "Create an attachment on an issue",
	Long:    `Attach a URL to an issue with a title and optional subtitle.`,
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

		urlFlag, _ := cmd.Flags().GetString("url")
		title, _ := cmd.Flags().GetString("title")

		if urlFlag == "" {
			output.Error("URL is required (--url)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Resolve issue ID (could be identifier like LIN-123)
		issue, err := client.GetIssue(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to resolve issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		input := map[string]interface{}{
			"issueId": issue.ID,
			"url":     urlFlag,
		}

		if title != "" {
			input["title"] = title
		}
		if cmd.Flags().Changed("subtitle") {
			subtitle, _ := cmd.Flags().GetString("subtitle")
			input["subtitle"] = subtitle
		}
		if cmd.Flags().Changed("icon-url") {
			iconURL, _ := cmd.Flags().GetString("icon-url")
			input["iconUrl"] = iconURL
		}
		if cmd.Flags().Changed("metadata") {
			metadataStr, _ := cmd.Flags().GetString("metadata")
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
				output.Error(fmt.Sprintf("Invalid metadata JSON: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["metadata"] = metadata
		}
		if cmd.Flags().Changed("group-by-source") {
			groupBySource, _ := cmd.Flags().GetBool("group-by-source")
			input["groupBySource"] = groupBySource
		}
		if cmd.Flags().Changed("comment-body") {
			commentBody, _ := cmd.Flags().GetString("comment-body")
			input["commentBody"] = commentBody
		}
		if cmd.Flags().Changed("create-as-user") {
			createAsUser, _ := cmd.Flags().GetString("create-as-user")
			input["createAsUser"] = createAsUser
		}

		attachment, err := client.CreateAttachment(context.Background(), input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to create attachment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(attachment)
		} else if plaintext {
			fmt.Printf("Created attachment: %s (%s)\n", attachment.Title, attachment.ID)
		} else {
			fmt.Printf("%s Created attachment %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(attachment.Title))
			fmt.Printf("  URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(attachment.URL))
		}
	},
}

var attachmentLinkCmd = &cobra.Command{
	Use:   "link [issue-id]",
	Short: "Smart link a URL to an issue",
	Long: `Smart link that auto-detects the URL type (GitHub PR, Slack thread, Notion page, etc.)
and creates the appropriate rich attachment.`,
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

		urlFlag, _ := cmd.Flags().GetString("url")
		title, _ := cmd.Flags().GetString("title")

		if urlFlag == "" {
			output.Error("URL is required (--url)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Resolve issue ID
		issue, err := client.GetIssue(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to resolve issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		attachment, err := client.LinkURL(context.Background(), issue.ID, urlFlag, title)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to link URL: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(attachment)
		} else if plaintext {
			fmt.Printf("Linked: %s (%s)\n", attachment.Title, attachment.URL)
		} else {
			fmt.Printf("%s Linked %s to issue\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(attachment.Title))
			fmt.Printf("  URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(attachment.URL))
		}
	},
}

var attachmentUpdateCmd = &cobra.Command{
	Use:     "update [attachment-id]",
	Aliases: []string{"edit"},
	Short:   "Update an attachment",
	Long:    `Update the title or subtitle of an attachment.`,
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
		input := make(map[string]interface{})

		if cmd.Flags().Changed("title") {
			title, _ := cmd.Flags().GetString("title")
			input["title"] = title
		}
		if cmd.Flags().Changed("subtitle") {
			subtitle, _ := cmd.Flags().GetString("subtitle")
			input["subtitle"] = subtitle
		}
		if cmd.Flags().Changed("icon-url") {
			iconURL, _ := cmd.Flags().GetString("icon-url")
			input["iconUrl"] = iconURL
		}
		if cmd.Flags().Changed("metadata") {
			metadataStr, _ := cmd.Flags().GetString("metadata")
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
				output.Error(fmt.Sprintf("Invalid metadata JSON: %v", err), plaintext, jsonOut)
				os.Exit(1)
			}
			input["metadata"] = metadata
		}

		if len(input) == 0 {
			output.Error("No updates specified. Use --title, --subtitle, --icon-url, or --metadata.", plaintext, jsonOut)
			os.Exit(1)
		}

		attachment, err := client.UpdateAttachment(context.Background(), args[0], input)
		if err != nil {
			output.Error(fmt.Sprintf("Failed to update attachment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(attachment)
		} else if plaintext {
			fmt.Printf("Updated attachment: %s\n", attachment.Title)
		} else {
			fmt.Printf("%s Updated attachment %s\n",
				color.New(color.FgGreen).Sprint("✓"),
				color.New(color.FgCyan, color.Bold).Sprint(attachment.Title))
		}
	},
}

var attachmentDeleteCmd = &cobra.Command{
	Use:     "delete [attachment-id]",
	Aliases: []string{"rm"},
	Short:   "Delete an attachment",
	Long:    `Delete an attachment from an issue.`,
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
		err = client.DeleteAttachment(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to delete attachment: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		if jsonOut {
			output.JSON(map[string]interface{}{
				"status":  "success",
				"message": "Attachment deleted",
			})
		} else if plaintext {
			fmt.Println("Attachment deleted")
		} else {
			fmt.Printf("%s Attachment deleted\n",
				color.New(color.FgGreen).Sprint("✓"))
		}
	},
}

var attachmentUploadCmd = &cobra.Command{
	Use:     "upload [issue-id]",
	Aliases: []string{"up"},
	Short:   "Upload a file as an attachment to an issue",
	Long:    `Upload local files (CSV, images, PDFs, etc.) as attachments to a Linear issue.`,
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

		// Get file paths from flag
		filePaths, _ := cmd.Flags().GetStringArray("file")
		if len(filePaths) == 0 {
			output.Error("At least one file is required (--file)", plaintext, jsonOut)
			os.Exit(1)
		}

		// Resolve issue ID
		issue, err := client.GetIssue(context.Background(), args[0])
		if err != nil {
			output.Error(fmt.Sprintf("Failed to resolve issue: %v", err), plaintext, jsonOut)
			os.Exit(1)
		}

		// Process each file
		for _, filePath := range filePaths {
			// Read file
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to read file %s: %v", filePath, err), plaintext, jsonOut)
				continue
			}

			// Get file info
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to stat file %s: %v", filePath, err), plaintext, jsonOut)
				continue
			}

			filename := fileInfo.Name()
			size := int(fileInfo.Size())

			// Auto-detect content type
			contentType := "application/octet-stream"
			if detected := detectContentType(filename); detected != "" {
				contentType = detected
			}

			// Get custom title or use filename
			title := filename
			if cmd.Flags().Changed("title") {
				title, _ = cmd.Flags().GetString("title")
			}

			if !plaintext && !jsonOut {
				fmt.Printf("Uploading %s (%d bytes)...\n", filename, size)
			}

			// Step 1: Request presigned upload URL
			uploadFile, err := client.FileUpload(context.Background(), filename, contentType, size, false)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to get upload URL: %v", err), plaintext, jsonOut)
				continue
			}

			// Step 2: Upload file to presigned URL
			err = client.UploadFileToURL(context.Background(), uploadFile.UploadURL, uploadFile.Headers, fileData, contentType)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to upload file: %v", err), plaintext, jsonOut)
				continue
			}

			// Step 3: Create attachment with the asset URL
			input := map[string]interface{}{
				"issueId": issue.ID,
				"url":     uploadFile.AssetURL,
				"title":   title,
			}

			attachment, err := client.CreateAttachment(context.Background(), input)
			if err != nil {
				output.Error(fmt.Sprintf("Failed to create attachment: %v", err), plaintext, jsonOut)
				continue
			}

			if jsonOut {
				output.JSON(attachment)
			} else if plaintext {
				fmt.Printf("Uploaded: %s (%s)\n", attachment.Title, attachment.URL)
			} else {
				fmt.Printf("%s Uploaded %s\n",
					color.New(color.FgGreen).Sprint("✓"),
					color.New(color.FgCyan, color.Bold).Sprint(attachment.Title))
				fmt.Printf("  URL: %s\n", color.New(color.FgBlue, color.Underline).Sprint(attachment.URL))
			}
		}
	},
}

func init() {
	issueCmd.AddCommand(attachmentCmd)
	attachmentCmd.AddCommand(attachmentListCmd)
	attachmentCmd.AddCommand(attachmentCreateCmd)
	attachmentCmd.AddCommand(attachmentLinkCmd)
	attachmentCmd.AddCommand(attachmentUploadCmd)
	attachmentCmd.AddCommand(attachmentUpdateCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)

	// List flags
	attachmentListCmd.Flags().IntP("limit", "l", 50, "Maximum number of attachments to fetch")

	// Create flags
	attachmentCreateCmd.Flags().String("url", "", "URL to attach (required)")
	attachmentCreateCmd.Flags().String("title", "", "Attachment title")
	attachmentCreateCmd.Flags().String("subtitle", "", "Attachment subtitle")
	attachmentCreateCmd.Flags().String("icon-url", "", "Custom icon URL for the attachment")
	attachmentCreateCmd.Flags().String("metadata", "", "Metadata as JSON object")
	attachmentCreateCmd.Flags().Bool("group-by-source", false, "Group by source in Linear UI")
	attachmentCreateCmd.Flags().String("comment-body", "", "Create a comment with the attachment")
	attachmentCreateCmd.Flags().String("create-as-user", "", "Create as a specific user (user ID)")
	_ = attachmentCreateCmd.MarkFlagRequired("url")

	// Link flags
	attachmentLinkCmd.Flags().String("url", "", "URL to link (required)")
	attachmentLinkCmd.Flags().String("title", "", "Optional title override")
	_ = attachmentLinkCmd.MarkFlagRequired("url")

	// Upload flags
	attachmentUploadCmd.Flags().StringArray("file", []string{}, "File path to upload (can be specified multiple times)")
	attachmentUploadCmd.Flags().String("title", "", "Custom title for the attachment (defaults to filename)")
	_ = attachmentUploadCmd.MarkFlagRequired("file")

	// Update flags
	attachmentUpdateCmd.Flags().String("title", "", "New title")
	attachmentUpdateCmd.Flags().String("subtitle", "", "New subtitle")
	attachmentUpdateCmd.Flags().String("icon-url", "", "New icon URL")
	attachmentUpdateCmd.Flags().String("metadata", "", "New metadata as JSON object")
}

// detectContentType returns the MIME type for a file based on its extension
func detectContentType(filename string) string {
	ext := ""
	for i := len(filename) - 1; i >= 0 && filename[i] != '.'; i-- {
		if filename[i] == '/' || filename[i] == '\\' {
			return ""
		}
	}
	if idx := len(filename) - 1; idx >= 0 && filename[idx] != '/' && filename[idx] != '\\' {
		for i := idx; i >= 0; i-- {
			if filename[i] == '.' {
				ext = filename[i:]
				break
			}
			if filename[i] == '/' || filename[i] == '\\' {
				break
			}
		}
	}

	// Common MIME types
	mimeTypes := map[string]string{
		".csv":  "text/csv",
		".txt":  "text/plain",
		".pdf":  "application/pdf",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".avi":  "video/x-msvideo",
		".webm": "video/webm",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".zip":  "application/zip",
		".tar":  "application/x-tar",
		".gz":   "application/gzip",
		".json": "application/json",
		".xml":  "application/xml",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".md":   "text/markdown",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}

	if contentType, ok := mimeTypes[ext]; ok {
		return contentType
	}

	return ""
}
