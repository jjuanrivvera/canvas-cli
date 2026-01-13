package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	filesCourseID     int64
	filesUserID       int64
	filesFolderID     int64
	filesContentTypes []string
	filesSearchTerm   string
	filesInclude      []string
	filesSort         string
	filesOrder        string
	filesOnDuplicate  string
	filesParentFolder int64
	filesDestination  string
	filesHidden       bool
	filesLocked       bool
	filesForce        bool
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

// filesListCmd represents the files list command
var filesListCmd = &cobra.Command{
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
	RunE: runFilesList,
}

// filesGetCmd represents the files get command
var filesGetCmd = &cobra.Command{
	Use:   "get <file-id>",
	Short: "Get file details",
	Long: `Get details of a specific file by ID.

Examples:
  canvas files get 456
  canvas files get 456 --include user`,
	Args: cobra.ExactArgs(1),
	RunE: runFilesGet,
}

// filesUploadCmd represents the files upload command
var filesUploadCmd = &cobra.Command{
	Use:   "upload <file-path>",
	Short: "Upload a file",
	Long: `Upload a file to a course, folder, or user's files.

You must specify one of --course-id, --folder-id, or --user-id.

Examples:
  canvas files upload document.pdf --course-id 123
  canvas files upload image.png --folder-id 456
  canvas files upload data.csv --user-id 789
  canvas files upload file.pdf --course-id 123 --on-duplicate overwrite`,
	Args: cobra.ExactArgs(1),
	RunE: runFilesUpload,
}

// filesDownloadCmd represents the files download command
var filesDownloadCmd = &cobra.Command{
	Use:   "download <file-id>",
	Short: "Download a file",
	Long: `Download a file from Canvas to your local system.

Examples:
  canvas files download 456
  canvas files download 456 --destination ./my-file.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: runFilesDownload,
}

// filesDeleteCmd represents the files delete command
var filesDeleteCmd = &cobra.Command{
	Use:   "delete <file-id>",
	Short: "Delete a file",
	Long: `Delete a file from Canvas.

Examples:
  canvas files delete 456`,
	Args: cobra.ExactArgs(1),
	RunE: runFilesDelete,
}

// filesQuotaCmd represents the files quota command
var filesQuotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Get storage quota information",
	Long: `Get storage quota information for a course or user.

You must specify either --course-id or --user-id.

Examples:
  canvas files quota --course-id 123
  canvas files quota --user-id 789`,
	RunE: runFilesQuota,
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.AddCommand(filesListCmd)
	filesCmd.AddCommand(filesGetCmd)
	filesCmd.AddCommand(filesUploadCmd)
	filesCmd.AddCommand(filesDownloadCmd)
	filesCmd.AddCommand(filesDeleteCmd)
	filesCmd.AddCommand(filesQuotaCmd)

	// List flags
	filesListCmd.Flags().Int64Var(&filesCourseID, "course-id", 0, "Course ID")
	filesListCmd.Flags().Int64Var(&filesFolderID, "folder-id", 0, "Folder ID")
	filesListCmd.Flags().Int64Var(&filesUserID, "user-id", 0, "User ID")
	filesListCmd.Flags().StringSliceVar(&filesContentTypes, "content-types", []string{}, "Filter by MIME types (comma-separated)")
	filesListCmd.Flags().StringVar(&filesSearchTerm, "search", "", "Search by file name")
	filesListCmd.Flags().StringSliceVar(&filesInclude, "include", []string{}, "Additional data to include (comma-separated)")
	filesListCmd.Flags().StringVar(&filesSort, "sort", "", "Sort by (name, size, created_at, updated_at, content_type)")
	filesListCmd.Flags().StringVar(&filesOrder, "order", "", "Order direction (asc, desc)")

	// Get flags
	filesGetCmd.Flags().StringSliceVar(&filesInclude, "include", []string{}, "Additional data to include (comma-separated)")

	// Upload flags
	filesUploadCmd.Flags().Int64Var(&filesCourseID, "course-id", 0, "Course ID")
	filesUploadCmd.Flags().Int64Var(&filesFolderID, "folder-id", 0, "Folder ID")
	filesUploadCmd.Flags().Int64Var(&filesUserID, "user-id", 0, "User ID")
	filesUploadCmd.Flags().StringVar(&filesOnDuplicate, "on-duplicate", "rename", "How to handle duplicates (overwrite, rename)")
	filesUploadCmd.Flags().Int64Var(&filesParentFolder, "parent-folder", 0, "Parent folder ID")
	filesUploadCmd.Flags().BoolVar(&filesHidden, "hidden", false, "Hide from students")
	filesUploadCmd.Flags().BoolVar(&filesLocked, "locked", false, "Lock the file")

	// Download flags
	filesDownloadCmd.Flags().StringVar(&filesDestination, "destination", "", "Destination file path (default: current directory with original filename)")

	// Delete flags
	filesDeleteCmd.Flags().BoolVarP(&filesForce, "force", "f", false, "Skip confirmation prompt")

	// Quota flags
	filesQuotaCmd.Flags().Int64Var(&filesCourseID, "course-id", 0, "Course ID")
	filesQuotaCmd.Flags().Int64Var(&filesUserID, "user-id", 0, "User ID")
}

func runFilesList(cmd *cobra.Command, args []string) error {
	// Validate that exactly one context is specified
	contextsSpecified := 0
	if filesCourseID > 0 {
		contextsSpecified++
	}
	if filesFolderID > 0 {
		contextsSpecified++
	}
	if filesUserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of --course-id, --folder-id, or --user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of --course-id, --folder-id, or --user-id")
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create files service
	filesService := api.NewFilesService(client)

	// Build options
	opts := &api.ListFilesOptions{
		ContentTypes: filesContentTypes,
		SearchTerm:   filesSearchTerm,
		Include:      filesInclude,
		Sort:         filesSort,
		Order:        filesOrder,
	}

	// List files based on context
	ctx := context.Background()
	var files []api.Attachment

	if filesCourseID > 0 {
		files, err = filesService.ListCourseFiles(ctx, filesCourseID, opts)
	} else if filesFolderID > 0 {
		files, err = filesService.ListFolderFiles(ctx, filesFolderID, opts)
	} else {
		files, err = filesService.ListUserFiles(ctx, filesUserID, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files found")
		return nil
	}

	// Display files
	fmt.Printf("Found %d files:\n\n", len(files))

	for _, file := range files {
		fmt.Printf("📄 %s\n", file.DisplayName)
		fmt.Printf("   ID: %d\n", file.ID)
		fmt.Printf("   Filename: %s\n", file.Filename)
		if file.Size > 0 {
			fmt.Printf("   Size: %s\n", formatFileSize(file.Size))
		}
		if file.ContentType != "" {
			fmt.Printf("   Type: %s\n", file.ContentType)
		}
		if !file.CreatedAt.IsZero() {
			fmt.Printf("   Created: %s\n", file.CreatedAt.Format("2006-01-02 15:04"))
		}
		if file.Locked {
			fmt.Printf("   🔒 Locked\n")
		}
		if file.Hidden {
			fmt.Printf("   👁️  Hidden\n")
		}
		fmt.Println()
	}

	return nil
}

func runFilesGet(cmd *cobra.Command, args []string) error {
	// Parse file ID
	fileID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create files service
	filesService := api.NewFilesService(client)

	// Get file
	ctx := context.Background()
	file, err := filesService.Get(ctx, fileID, filesInclude)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	// Display file details
	fmt.Printf("📄 %s\n", file.DisplayName)
	fmt.Printf("   ID: %d\n", file.ID)
	fmt.Printf("   Filename: %s\n", file.Filename)
	if file.UUID != "" {
		fmt.Printf("   UUID: %s\n", file.UUID)
	}
	if file.Size > 0 {
		fmt.Printf("   Size: %s\n", formatFileSize(file.Size))
	}
	if file.ContentType != "" {
		fmt.Printf("   Content Type: %s\n", file.ContentType)
	}
	if file.URL != "" {
		fmt.Printf("   URL: %s\n", file.URL)
	}
	if !file.CreatedAt.IsZero() {
		fmt.Printf("   Created: %s\n", file.CreatedAt.Format("2006-01-02 15:04"))
	}
	if !file.UpdatedAt.IsZero() {
		fmt.Printf("   Updated: %s\n", file.UpdatedAt.Format("2006-01-02 15:04"))
	}
	if file.Locked {
		fmt.Printf("   🔒 Locked\n")
	}
	if file.Hidden {
		fmt.Printf("   👁️  Hidden\n")
	}

	return nil
}

func runFilesUpload(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Validate that exactly one context is specified
	contextsSpecified := 0
	if filesCourseID > 0 {
		contextsSpecified++
	}
	if filesFolderID > 0 {
		contextsSpecified++
	}
	if filesUserID > 0 {
		contextsSpecified++
	}

	if contextsSpecified == 0 {
		return fmt.Errorf("must specify one of --course-id, --folder-id, or --user-id")
	}
	if contextsSpecified > 1 {
		return fmt.Errorf("can only specify one of --course-id, --folder-id, or --user-id")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create files service
	filesService := api.NewFilesService(client)

	// Build upload parameters
	params := &api.UploadParams{
		OnDuplicate:    filesOnDuplicate,
		ParentFolderID: filesParentFolder,
		Hidden:         filesHidden,
		Locked:         filesLocked,
	}

	// Upload file based on context
	ctx := context.Background()
	var uploadedFile *api.Attachment

	fmt.Printf("Uploading %s...\n", filepath.Base(filePath))

	if filesCourseID > 0 {
		uploadedFile, err = filesService.UploadToCourse(ctx, filesCourseID, filePath, params)
	} else if filesFolderID > 0 {
		uploadedFile, err = filesService.UploadToFolder(ctx, filesFolderID, filePath, params)
	} else {
		uploadedFile, err = filesService.UploadToUser(ctx, filesUserID, filePath, params)
	}

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	// Display success message
	fmt.Printf("✅ File uploaded successfully\n\n")
	fmt.Printf("   ID: %d\n", uploadedFile.ID)
	fmt.Printf("   Name: %s\n", uploadedFile.DisplayName)
	fmt.Printf("   Size: %s\n", formatFileSize(uploadedFile.Size))
	fmt.Printf("   URL: %s\n", uploadedFile.URL)

	return nil
}

func runFilesDownload(cmd *cobra.Command, args []string) error {
	// Parse file ID
	fileID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file ID: %w", err)
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create files service
	filesService := api.NewFilesService(client)
	ctx := context.Background()

	// Get file info first to get the filename
	file, err := filesService.Get(ctx, fileID, nil)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Determine destination path
	destPath := filesDestination
	if destPath == "" {
		destPath = file.Filename
	}

	fmt.Printf("Downloading %s...\n", file.DisplayName)

	// Download file
	if err := filesService.Download(ctx, fileID, destPath); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	fmt.Printf("✅ File downloaded to %s\n", destPath)

	return nil
}

func runFilesDelete(cmd *cobra.Command, args []string) error {
	// Parse file ID
	fileID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file ID: %w", err)
	}

	// Confirm deletion
	confirmed, err := confirmDelete("file", fileID, filesForce)
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Delete cancelled")
		return nil
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create files service
	filesService := api.NewFilesService(client)

	// Delete file
	ctx := context.Background()
	if err := filesService.Delete(ctx, fileID); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	fmt.Printf("✅ File %d deleted successfully\n", fileID)

	return nil
}

func runFilesQuota(cmd *cobra.Command, args []string) error {
	// Validate that exactly one context is specified
	if filesCourseID == 0 && filesUserID == 0 {
		return fmt.Errorf("must specify either --course-id or --user-id")
	}
	if filesCourseID > 0 && filesUserID > 0 {
		return fmt.Errorf("can only specify one of --course-id or --user-id")
	}

	// Get API client
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	// Create files service
	filesService := api.NewFilesService(client)
	ctx := context.Background()

	// Get quota based on context
	var quota *api.QuotaInfo
	if filesCourseID > 0 {
		quota, err = filesService.GetCourseQuota(ctx, filesCourseID)
	} else {
		quota, err = filesService.GetUserQuota(ctx, filesUserID)
	}

	if err != nil {
		return fmt.Errorf("failed to get quota: %w", err)
	}

	// Display quota information
	fmt.Println("Storage Quota:")
	fmt.Printf("   Used: %s\n", formatFileSize(quota.QuotaUsed))
	fmt.Printf("   Total: %s\n", formatFileSize(quota.Quota))

	if quota.Quota > 0 {
		percentUsed := float64(quota.QuotaUsed) / float64(quota.Quota) * 100
		fmt.Printf("   Usage: %.1f%%\n", percentUsed)
	}

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
