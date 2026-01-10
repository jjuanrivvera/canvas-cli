package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	pagesCourseID      int64
	pagesSort          string
	pagesOrder         string
	pagesSearchTerm    string
	pagesPublished     string
	pagesInclude       []string
	pagesTitle         string
	pagesBody          string
	pagesEditingRoles  string
	pagesNotifyUpdate  bool
	pagesMakePublished bool
	pagesFrontPage     bool
	pagesPublishAt     string
)

// pagesCmd represents the pages command group
var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Manage Canvas wiki pages",
	Long: `Manage Canvas wiki pages including listing, viewing, creating, and updating pages.

Pages are rich content associated with Courses in Canvas. They can be used
for course information, resources, or any other content.

Examples:
  canvas pages list --course-id 123
  canvas pages get --course-id 123 my-page-url
  canvas pages create --course-id 123 --title "Welcome" --body "<p>Hello!</p>"
  canvas pages front --course-id 123`,
}

// pagesListCmd represents the pages list command
var pagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pages in a course",
	Long: `List all wiki pages in a Canvas course.

Examples:
  canvas pages list --course-id 123
  canvas pages list --course-id 123 --sort title --order asc
  canvas pages list --course-id 123 --search "intro"
  canvas pages list --course-id 123 --published true`,
	RunE: runPagesList,
}

// pagesGetCmd represents the pages get command
var pagesGetCmd = &cobra.Command{
	Use:   "get <url-or-id>",
	Short: "Get a specific page",
	Long: `Get details of a specific wiki page by URL or ID.

Examples:
  canvas pages get --course-id 123 my-page-title
  canvas pages get --course-id 123 page_id:456`,
	Args: cobra.ExactArgs(1),
	RunE: runPagesGet,
}

// pagesFrontCmd represents the pages front command
var pagesFrontCmd = &cobra.Command{
	Use:   "front",
	Short: "Get the front page",
	Long: `Get the front page for a course.

Examples:
  canvas pages front --course-id 123`,
	RunE: runPagesFront,
}

// pagesCreateCmd represents the pages create command
var pagesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new page",
	Long: `Create a new wiki page in a course.

Examples:
  canvas pages create --course-id 123 --title "Welcome"
  canvas pages create --course-id 123 --title "Syllabus" --body "<p>Course syllabus</p>"
  canvas pages create --course-id 123 --title "Home" --front-page --published`,
	RunE: runPagesCreate,
}

// pagesUpdateCmd represents the pages update command
var pagesUpdateCmd = &cobra.Command{
	Use:   "update <url-or-id>",
	Short: "Update an existing page",
	Long: `Update an existing wiki page.

Examples:
  canvas pages update --course-id 123 my-page --title "New Title"
  canvas pages update --course-id 123 my-page --body "<p>Updated content</p>"
  canvas pages update --course-id 123 my-page --published`,
	Args: cobra.ExactArgs(1),
	RunE: runPagesUpdate,
}

// pagesDeleteCmd represents the pages delete command
var pagesDeleteCmd = &cobra.Command{
	Use:   "delete <url-or-id>",
	Short: "Delete a page",
	Long: `Delete a wiki page from a course.

Examples:
  canvas pages delete --course-id 123 my-page-url`,
	Args: cobra.ExactArgs(1),
	RunE: runPagesDelete,
}

// pagesDuplicateCmd represents the pages duplicate command
var pagesDuplicateCmd = &cobra.Command{
	Use:   "duplicate <url-or-id>",
	Short: "Duplicate a page",
	Long: `Duplicate a wiki page.

Examples:
  canvas pages duplicate --course-id 123 my-page-url`,
	Args: cobra.ExactArgs(1),
	RunE: runPagesDuplicate,
}

// pagesRevisionsCmd represents the pages revisions command
var pagesRevisionsCmd = &cobra.Command{
	Use:   "revisions <url-or-id>",
	Short: "List page revisions",
	Long: `List all revisions for a wiki page.

Examples:
  canvas pages revisions --course-id 123 my-page-url`,
	Args: cobra.ExactArgs(1),
	RunE: runPagesRevisions,
}

// pagesRevertCmd represents the pages revert command
var pagesRevertCmd = &cobra.Command{
	Use:   "revert <url-or-id> <revision-id>",
	Short: "Revert to a specific revision",
	Long: `Revert a wiki page to a specific revision.

Examples:
  canvas pages revert --course-id 123 my-page-url 5`,
	Args: cobra.ExactArgs(2),
	RunE: runPagesRevert,
}

func init() {
	rootCmd.AddCommand(pagesCmd)
	pagesCmd.AddCommand(pagesListCmd)
	pagesCmd.AddCommand(pagesGetCmd)
	pagesCmd.AddCommand(pagesFrontCmd)
	pagesCmd.AddCommand(pagesCreateCmd)
	pagesCmd.AddCommand(pagesUpdateCmd)
	pagesCmd.AddCommand(pagesDeleteCmd)
	pagesCmd.AddCommand(pagesDuplicateCmd)
	pagesCmd.AddCommand(pagesRevisionsCmd)
	pagesCmd.AddCommand(pagesRevertCmd)

	// List flags
	pagesListCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesListCmd.Flags().StringVar(&pagesSort, "sort", "", "Sort by: title, created_at, updated_at")
	pagesListCmd.Flags().StringVar(&pagesOrder, "order", "", "Sort order: asc, desc")
	pagesListCmd.Flags().StringVar(&pagesSearchTerm, "search", "", "Search by page title")
	pagesListCmd.Flags().StringVar(&pagesPublished, "published", "", "Filter by published status: true, false")
	pagesListCmd.Flags().StringSliceVar(&pagesInclude, "include", []string{}, "Additional data to include (body)")
	pagesListCmd.MarkFlagRequired("course-id")

	// Get flags
	pagesGetCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesGetCmd.MarkFlagRequired("course-id")

	// Front flags
	pagesFrontCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesFrontCmd.MarkFlagRequired("course-id")

	// Create flags
	pagesCreateCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesCreateCmd.Flags().StringVar(&pagesTitle, "title", "", "Page title (required)")
	pagesCreateCmd.Flags().StringVar(&pagesBody, "body", "", "Page body (HTML)")
	pagesCreateCmd.Flags().StringVar(&pagesEditingRoles, "editing-roles", "", "Roles that can edit: teachers,students,members,public")
	pagesCreateCmd.Flags().BoolVar(&pagesNotifyUpdate, "notify", false, "Notify participants of update")
	pagesCreateCmd.Flags().BoolVar(&pagesMakePublished, "published", false, "Publish the page")
	pagesCreateCmd.Flags().BoolVar(&pagesFrontPage, "front-page", false, "Set as front page")
	pagesCreateCmd.Flags().StringVar(&pagesPublishAt, "publish-at", "", "Schedule publication date (ISO 8601)")
	pagesCreateCmd.MarkFlagRequired("course-id")
	pagesCreateCmd.MarkFlagRequired("title")

	// Update flags
	pagesUpdateCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesUpdateCmd.Flags().StringVar(&pagesTitle, "title", "", "New page title")
	pagesUpdateCmd.Flags().StringVar(&pagesBody, "body", "", "New page body (HTML)")
	pagesUpdateCmd.Flags().StringVar(&pagesEditingRoles, "editing-roles", "", "Roles that can edit: teachers,students,members,public")
	pagesUpdateCmd.Flags().BoolVar(&pagesNotifyUpdate, "notify", false, "Notify participants of update")
	pagesUpdateCmd.Flags().BoolVar(&pagesMakePublished, "published", false, "Publish the page")
	pagesUpdateCmd.Flags().BoolVar(&pagesFrontPage, "front-page", false, "Set as front page")
	pagesUpdateCmd.Flags().StringVar(&pagesPublishAt, "publish-at", "", "Schedule publication date (ISO 8601)")
	pagesUpdateCmd.MarkFlagRequired("course-id")

	// Delete flags
	pagesDeleteCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesDeleteCmd.MarkFlagRequired("course-id")

	// Duplicate flags
	pagesDuplicateCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesDuplicateCmd.MarkFlagRequired("course-id")

	// Revisions flags
	pagesRevisionsCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesRevisionsCmd.MarkFlagRequired("course-id")

	// Revert flags
	pagesRevertCmd.Flags().Int64Var(&pagesCourseID, "course-id", 0, "Course ID (required)")
	pagesRevertCmd.MarkFlagRequired("course-id")
}

func runPagesList(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	opts := &api.ListPagesOptions{
		Sort:       pagesSort,
		Order:      pagesOrder,
		SearchTerm: pagesSearchTerm,
		Include:    pagesInclude,
	}

	if pagesPublished != "" {
		pub := pagesPublished == "true"
		opts.Published = &pub
	}

	ctx := context.Background()
	pages, err := pagesService.List(ctx, pagesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to list pages: %w", err)
	}

	if len(pages) == 0 {
		fmt.Println("No pages found")
		return nil
	}

	fmt.Printf("Found %d pages:\n\n", len(pages))

	for _, page := range pages {
		displayPage(&page)
	}

	return nil
}

func runPagesGet(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	ctx := context.Background()
	page, err := pagesService.Get(ctx, pagesCourseID, args[0])
	if err != nil {
		return fmt.Errorf("failed to get page: %w", err)
	}

	displayPageFull(page)

	return nil
}

func runPagesFront(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	ctx := context.Background()
	page, err := pagesService.GetFrontPage(ctx, pagesCourseID)
	if err != nil {
		return fmt.Errorf("failed to get front page: %w", err)
	}

	displayPageFull(page)

	return nil
}

func runPagesCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	params := &api.CreatePageParams{
		Title:          pagesTitle,
		Body:           pagesBody,
		EditingRoles:   pagesEditingRoles,
		NotifyOfUpdate: pagesNotifyUpdate,
		Published:      pagesMakePublished,
		FrontPage:      pagesFrontPage,
		PublishAt:      pagesPublishAt,
	}

	ctx := context.Background()
	page, err := pagesService.Create(ctx, pagesCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	fmt.Println("Page created successfully!")
	displayPage(page)

	return nil
}

func runPagesUpdate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	params := &api.UpdatePageParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &pagesTitle
	}
	if cmd.Flags().Changed("body") {
		params.Body = &pagesBody
	}
	if cmd.Flags().Changed("editing-roles") {
		params.EditingRoles = &pagesEditingRoles
	}
	if cmd.Flags().Changed("notify") {
		params.NotifyOfUpdate = &pagesNotifyUpdate
	}
	if cmd.Flags().Changed("published") {
		params.Published = &pagesMakePublished
	}
	if cmd.Flags().Changed("front-page") {
		params.FrontPage = &pagesFrontPage
	}
	if cmd.Flags().Changed("publish-at") {
		params.PublishAt = &pagesPublishAt
	}

	ctx := context.Background()
	page, err := pagesService.Update(ctx, pagesCourseID, args[0], params)
	if err != nil {
		return fmt.Errorf("failed to update page: %w", err)
	}

	fmt.Println("Page updated successfully!")
	displayPage(page)

	return nil
}

func runPagesDelete(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	ctx := context.Background()
	if err := pagesService.Delete(ctx, pagesCourseID, args[0]); err != nil {
		return fmt.Errorf("failed to delete page: %w", err)
	}

	fmt.Printf("Page '%s' deleted successfully\n", args[0])
	return nil
}

func runPagesDuplicate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	ctx := context.Background()
	page, err := pagesService.Duplicate(ctx, pagesCourseID, args[0])
	if err != nil {
		return fmt.Errorf("failed to duplicate page: %w", err)
	}

	fmt.Println("Page duplicated successfully!")
	displayPage(page)

	return nil
}

func runPagesRevisions(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	ctx := context.Background()
	revisions, err := pagesService.ListRevisions(ctx, pagesCourseID, args[0])
	if err != nil {
		return fmt.Errorf("failed to list revisions: %w", err)
	}

	if len(revisions) == 0 {
		fmt.Println("No revisions found")
		return nil
	}

	fmt.Printf("Found %d revisions:\n\n", len(revisions))

	for _, rev := range revisions {
		displayRevision(&rev)
	}

	return nil
}

func runPagesRevert(cmd *cobra.Command, args []string) error {
	revisionID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid revision ID: %s", args[1])
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	pagesService := api.NewPagesService(client)

	ctx := context.Background()
	revision, err := pagesService.RevertToRevision(ctx, pagesCourseID, args[0], revisionID)
	if err != nil {
		return fmt.Errorf("failed to revert to revision: %w", err)
	}

	fmt.Println("Page reverted successfully!")
	displayRevision(revision)

	return nil
}

func displayPage(page *api.Page) {
	stateIcon := "📄"
	if page.FrontPage {
		stateIcon = "🏠"
	} else if page.Published {
		stateIcon = "📝"
	}

	fmt.Printf("%s %s\n", stateIcon, page.Title)
	fmt.Printf("   URL: %s\n", page.URL)

	if page.Published {
		fmt.Printf("   Published: Yes\n")
	} else {
		fmt.Printf("   Published: No (Draft)\n")
	}

	if page.FrontPage {
		fmt.Printf("   Front Page: Yes\n")
	}

	fmt.Printf("   Updated: %s\n", page.UpdatedAt.Format("2006-01-02 15:04"))

	fmt.Println()
}

func displayPageFull(page *api.Page) {
	displayPage(page)

	if page.EditingRoles != "" {
		fmt.Printf("   Editing Roles: %s\n", page.EditingRoles)
	}

	if page.LastEditedBy != nil {
		fmt.Printf("   Last Edited By: %s\n", page.LastEditedBy.Name)
	}

	if page.Body != "" {
		fmt.Printf("\nContent:\n")
		// Truncate body for display
		body := page.Body
		if len(body) > 500 {
			body = body[:500] + "..."
		}
		// Strip HTML tags for display
		body = stripHTMLTags(body)
		fmt.Println(body)
	}

	fmt.Println()
}

func displayRevision(rev *api.PageRevision) {
	latestMark := ""
	if rev.Latest {
		latestMark = " [Latest]"
	}

	fmt.Printf("Revision %d%s\n", rev.RevisionID, latestMark)
	fmt.Printf("   Updated: %s\n", rev.UpdatedAt.Format("2006-01-02 15:04"))

	if rev.EditedBy != nil {
		fmt.Printf("   Edited By: %s\n", rev.EditedBy.Name)
	}

	if rev.Title != "" {
		fmt.Printf("   Title: %s\n", rev.Title)
	}

	fmt.Println()
}

func stripHTMLTags(s string) string {
	// Simple HTML tag stripper
	result := s
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return result
}
