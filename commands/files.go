package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

// filesCmd represents the files command group
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Manage Canvas files",
	Long: `Manage Canvas files including listing, uploading, downloading, and deleting files.

Examples:
  canvas files list --course-id 123
  canvas files get 456
  canvas files upload --course-id 123 document.pdf
  canvas files download 456 --destination ./downloaded.pdf
  canvas files delete 456`,
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(newFilesListCmd())
	filesCmd.AddCommand(newFilesGetCmd())
	filesCmd.AddCommand(newFilesUploadCmd())
	filesCmd.AddCommand(newFilesDownloadCmd())
	filesCmd.AddCommand(newFilesDeleteCmd())
	filesCmd.AddCommand(newFilesQuotaCmd())
}

func newFilesListCmd() *cobra.Command {
	opts := &options.FilesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files",
		Long: `List files in a course, folder, or user's files.

You must specify one of --course-id, --folder-id, or --user-id.

Examples:
  canvas files list --course-id 123
  canvas files list --folder-id 456
  canvas files list --user-id 789
  canvas files list --course-id 123 --search "assignment"
  canvas files list --course-id 123 --sort name --order asc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().StringSliceVar(&opts.ContentTypes, "content-types", []string{}, "Filter by MIME types (comma-separated)")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by file name")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by (name, size, created_at, updated_at, content_type)")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Order direction (asc, desc)")

	return cmd
}

func newFilesGetCmd() *cobra.Command {
	opts := &options.FilesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <file-id>",
		Short: "Get file details",
		Long: `Get details of a specific file by ID.

Examples:
  canvas files get 456
  canvas files get 456 --include user`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (comma-separated)")

	return cmd
}

func newFilesUploadCmd() *cobra.Command {
	opts := &options.FilesUploadOptions{}

	cmd := &cobra.Command{
		Use:   "upload <file-path>",
		Short: "Upload a file",
		Long: `Upload a file to a course, folder, or user's files.

You must specify one of --course-id, --folder-id, or --user-id.

Examples:
  canvas files upload document.pdf --course-id 123
  canvas files upload image.png --folder-id 456
  canvas files upload data.csv --user-id 789
  canvas files upload file.pdf --course-id 123 --on-duplicate overwrite`,
		Args: ExactArgsWithUsage(1, "file-path"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FilePath = args[0]

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesUpload(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.FolderID, "folder-id", 0, "Folder ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")
	cmd.Flags().StringVar(&opts.OnDuplicate, "on-duplicate", "rename", "How to handle duplicates (overwrite, rename)")
	cmd.Flags().Int64Var(&opts.ParentFolderID, "parent-folder", 0, "Parent folder ID")
	cmd.Flags().BoolVar(&opts.Hidden, "hidden", false, "Hide from students")
	cmd.Flags().BoolVar(&opts.Locked, "locked", false, "Lock the file")

	return cmd
}

func newFilesDownloadCmd() *cobra.Command {
	opts := &options.FilesDownloadOptions{}

	cmd := &cobra.Command{
		Use:   "download <file-id>",
		Short: "Download a file",
		Long: `Download a file from Canvas to your local system.

Examples:
  canvas files download 456
  canvas files download 456 --destination ./my-file.pdf`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesDownload(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Destination, "destination", "", "Destination file path (default: current directory with original filename)")

	return cmd
}

func newFilesDeleteCmd() *cobra.Command {
	opts := &options.FilesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <file-id>",
		Short: "Delete a file",
		Long: `Delete a file from Canvas.

Examples:
  canvas files delete 456`,
		Args: ExactArgsWithUsage(1, "file-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid file ID: %s", args[0])
			}
			opts.FileID = fileID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")

	return cmd
}

func newFilesQuotaCmd() *cobra.Command {
	opts := &options.FilesQuotaOptions{}

	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Get storage quota information",
		Long: `Get storage quota information for a course or user.

You must specify either --course-id or --user-id.

Examples:
  canvas files quota --course-id 123
  canvas files quota --user-id 789`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runFilesQuota(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.UserID, "user-id", 0, "User ID")

	return cmd
}

func runFilesList(ctx context.Context, client *api.Client, opts *options.FilesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"folder_id": opts.FolderID,
		"user_id":   opts.UserID,
	})

	filesService := api.NewFilesService(client)

	apiOpts := &api.ListFilesOptions{
		ContentTypes: opts.ContentTypes,
		SearchTerm:   opts.SearchTerm,
		Include:      opts.Include,
		Sort:         opts.Sort,
		Order:        opts.Order,
	}

	var files []api.Attachment
	var err error

	if opts.CourseID > 0 {
		files, err = filesService.ListCourseFiles(ctx, opts.CourseID, apiOpts)
	} else if opts.FolderID > 0 {
		files, err = filesService.ListFolderFiles(ctx, opts.FolderID, apiOpts)
	} else {
		files, err = filesService.ListUserFiles(ctx, opts.UserID, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "files.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"folder_id": opts.FolderID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files found")
		logger.LogCommandComplete(ctx, "files.list", 0)
		return nil
	}

	printVerbose("Found %d files:\n\n", len(files))

	logger.LogCommandComplete(ctx, "files.list", len(files))
	return formatOutput(files, nil)
}

func runFilesGet(ctx context.Context, client *api.Client, opts *options.FilesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.get", map[string]interface{}{
		"file_id": opts.FileID,
	})

	filesService := api.NewFilesService(client)

	file, err := filesService.Get(ctx, opts.FileID, opts.Include)
	if err != nil {
		logger.LogCommandError(ctx, "files.get", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to get file: %w", err)
	}

	logger.LogCommandComplete(ctx, "files.get", 1)
	return formatOutput(file, nil)
}

func runFilesUpload(ctx context.Context, client *api.Client, opts *options.FilesUploadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.upload", map[string]interface{}{
		"file_path": opts.FilePath,
		"course_id": opts.CourseID,
		"folder_id": opts.FolderID,
		"user_id":   opts.UserID,
	})

	// Check if file exists
	if _, err := os.Stat(opts.FilePath); os.IsNotExist(err) {
		logger.LogCommandError(ctx, "files.upload", err, map[string]interface{}{
			"file_path": opts.FilePath,
		})
		return fmt.Errorf("file does not exist: %s", opts.FilePath)
	}

	filesService := api.NewFilesService(client)

	params := &api.UploadParams{
		OnDuplicate:    opts.OnDuplicate,
		ParentFolderID: opts.ParentFolderID,
		Hidden:         opts.Hidden,
		Locked:         opts.Locked,
	}

	fmt.Printf("Uploading %s...\n", filepath.Base(opts.FilePath))

	var uploadedFile *api.Attachment
	var err error

	if opts.CourseID > 0 {
		uploadedFile, err = filesService.UploadToCourse(ctx, opts.CourseID, opts.FilePath, params)
	} else if opts.FolderID > 0 {
		uploadedFile, err = filesService.UploadToFolder(ctx, opts.FolderID, opts.FilePath, params)
	} else {
		uploadedFile, err = filesService.UploadToUser(ctx, opts.UserID, opts.FilePath, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "files.upload", err, map[string]interface{}{
			"file_path": opts.FilePath,
			"course_id": opts.CourseID,
			"folder_id": opts.FolderID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to upload file: %w", err)
	}

	printInfo("✅ File uploaded successfully\n\n")
	fmt.Printf("   ID: %d\n", uploadedFile.ID)
	fmt.Printf("   Name: %s\n", uploadedFile.DisplayName)
	fmt.Printf("   Size: %s\n", formatFileSize(uploadedFile.Size))
	fmt.Printf("   URL: %s\n", uploadedFile.URL)

	logger.LogCommandComplete(ctx, "files.upload", 1)
	return nil
}

func runFilesDownload(ctx context.Context, client *api.Client, opts *options.FilesDownloadOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.download", map[string]interface{}{
		"file_id":     opts.FileID,
		"destination": opts.Destination,
	})

	filesService := api.NewFilesService(client)

	// Get file info first to get the filename
	file, err := filesService.Get(ctx, opts.FileID, nil)
	if err != nil {
		logger.LogCommandError(ctx, "files.download", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine destination path
	destPath := opts.Destination
	if destPath == "" {
		destPath = file.Filename
	}

	fmt.Printf("Downloading %s...\n", file.DisplayName)

	// Download file
	if err := filesService.Download(ctx, opts.FileID, destPath); err != nil {
		logger.LogCommandError(ctx, "files.download", err, map[string]interface{}{
			"file_id":     opts.FileID,
			"destination": destPath,
		})
		return fmt.Errorf("failed to download file: %w", err)
	}

	fmt.Printf("✅ File downloaded to %s\n", destPath)

	logger.LogCommandComplete(ctx, "files.download", 1)
	return nil
}

func runFilesDelete(ctx context.Context, client *api.Client, opts *options.FilesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.delete", map[string]interface{}{
		"file_id": opts.FileID,
		"force":   opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("file", opts.FileID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "files.delete", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		logger.LogCommandComplete(ctx, "files.delete", 0)
		return nil
	}

	filesService := api.NewFilesService(client)

	if err := filesService.Delete(ctx, opts.FileID); err != nil {
		logger.LogCommandError(ctx, "files.delete", err, map[string]interface{}{
			"file_id": opts.FileID,
		})
		return fmt.Errorf("failed to delete file: %w", err)
	}

	printInfo("✅ File %d deleted successfully\n", opts.FileID)

	logger.LogCommandComplete(ctx, "files.delete", 1)
	return nil
}

func runFilesQuota(ctx context.Context, client *api.Client, opts *options.FilesQuotaOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "files.quota", map[string]interface{}{
		"course_id": opts.CourseID,
		"user_id":   opts.UserID,
	})

	filesService := api.NewFilesService(client)

	var quota *api.QuotaInfo
	var err error

	if opts.CourseID > 0 {
		quota, err = filesService.GetCourseQuota(ctx, opts.CourseID)
	} else {
		quota, err = filesService.GetUserQuota(ctx, opts.UserID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "files.quota", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"user_id":   opts.UserID,
		})
		return fmt.Errorf("failed to get quota: %w", err)
	}

	fmt.Println("Storage Quota:")
	fmt.Printf("   Used: %s\n", formatFileSize(quota.QuotaUsed))
	fmt.Printf("   Total: %s\n", formatFileSize(quota.Quota))

	if quota.Quota > 0 {
		percentUsed := float64(quota.QuotaUsed) / float64(quota.Quota) * 100
		fmt.Printf("   Usage: %.1f%%\n", percentUsed)
	}

	logger.LogCommandComplete(ctx, "files.quota", 1)
	return nil
}

// formatFileSize formats a file size in bytes to a human-readable string
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
