package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/commands/internal/logging"
	"github.com/jjuanrivvera/canvas-cli/commands/internal/options"
	"github.com/jjuanrivvera/canvas-cli/internal/api"
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

func init() {
	rootCmd.AddCommand(pagesCmd)
	pagesCmd.AddCommand(newPagesListCmd())
	pagesCmd.AddCommand(newPagesGetCmd())
	pagesCmd.AddCommand(newPagesFrontCmd())
	pagesCmd.AddCommand(newPagesCreateCmd())
	pagesCmd.AddCommand(newPagesUpdateCmd())
	pagesCmd.AddCommand(newPagesDeleteCmd())
	pagesCmd.AddCommand(newPagesDuplicateCmd())
	pagesCmd.AddCommand(newPagesRevisionsCmd())
	pagesCmd.AddCommand(newPagesRevertCmd())
}

func newPagesListCmd() *cobra.Command {
	opts := &options.PagesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pages in a course",
		Long: `List all wiki pages in a Canvas course.

Examples:
  canvas pages list --course-id 123
  canvas pages list --course-id 123 --sort title --order asc
  canvas pages list --course-id 123 --search "intro"
  canvas pages list --course-id 123 --published true`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Sort, "sort", "", "Sort by: title, created_at, updated_at")
	cmd.Flags().StringVar(&opts.Order, "order", "", "Sort order: asc, desc")
	cmd.Flags().StringVar(&opts.SearchTerm, "search", "", "Search by page title")
	cmd.Flags().StringVar(&opts.Published, "published", "", "Filter by published status: true, false")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Additional data to include (body)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesGetCmd() *cobra.Command {
	opts := &options.PagesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <url-or-id>",
		Short: "Get a specific page",
		Long: `Get details of a specific wiki page by URL or ID.

Examples:
  canvas pages get --course-id 123 my-page-title
  canvas pages get --course-id 123 page_id:456`,
		Args: ExactArgsWithUsage(1, "url-or-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URLOrID = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesFrontCmd() *cobra.Command {
	opts := &options.PagesFrontOptions{}

	cmd := &cobra.Command{
		Use:   "front",
		Short: "Get the front page",
		Long: `Get the front page for a course.

Examples:
  canvas pages front --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesFront(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesCreateCmd() *cobra.Command {
	opts := &options.PagesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new page",
		Long: `Create a new wiki page in a course.

Examples:
  canvas pages create --course-id 123 --title "Welcome"
  canvas pages create --course-id 123 --title "Syllabus" --body "<p>Course syllabus</p>"
  canvas pages create --course-id 123 --title "Home" --front-page --published`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Page title (required)")
	cmd.Flags().StringVar(&opts.Body, "body", "", "Page body (HTML)")
	cmd.Flags().StringVar(&opts.EditingRoles, "editing-roles", "", "Roles that can edit: teachers,students,members,public")
	cmd.Flags().BoolVar(&opts.NotifyUpdate, "notify", false, "Notify participants of update")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the page")
	cmd.Flags().BoolVar(&opts.FrontPage, "front-page", false, "Set as front page")
	cmd.Flags().StringVar(&opts.PublishAt, "publish-at", "", "Schedule publication date (ISO 8601)")
	cmd.MarkFlagRequired("course-id")
	cmd.MarkFlagRequired("title")

	return cmd
}

func newPagesUpdateCmd() *cobra.Command {
	opts := &options.PagesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <url-or-id>",
		Short: "Update an existing page",
		Long: `Update an existing wiki page.

Examples:
  canvas pages update --course-id 123 my-page --title "New Title"
  canvas pages update --course-id 123 my-page --body "<p>Updated content</p>"
  canvas pages update --course-id 123 my-page --published`,
		Args: ExactArgsWithUsage(1, "url-or-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URLOrID = args[0]
			// Track which fields were changed
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.BodySet = cmd.Flags().Changed("body")
			opts.EditingRolesSet = cmd.Flags().Changed("editing-roles")
			opts.NotifyUpdateSet = cmd.Flags().Changed("notify")
			opts.PublishedSet = cmd.Flags().Changed("published")
			opts.FrontPageSet = cmd.Flags().Changed("front-page")
			opts.PublishAtSet = cmd.Flags().Changed("publish-at")
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "New page title")
	cmd.Flags().StringVar(&opts.Body, "body", "", "New page body (HTML)")
	cmd.Flags().StringVar(&opts.EditingRoles, "editing-roles", "", "Roles that can edit: teachers,students,members,public")
	cmd.Flags().BoolVar(&opts.NotifyUpdate, "notify", false, "Notify participants of update")
	cmd.Flags().BoolVar(&opts.Published, "published", false, "Publish the page")
	cmd.Flags().BoolVar(&opts.FrontPage, "front-page", false, "Set as front page")
	cmd.Flags().StringVar(&opts.PublishAt, "publish-at", "", "Schedule publication date (ISO 8601)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesDeleteCmd() *cobra.Command {
	opts := &options.PagesDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <url-or-id>",
		Short: "Delete a page",
		Long: `Delete a wiki page from a course.

Examples:
  canvas pages delete --course-id 123 my-page-url`,
		Args: ExactArgsWithUsage(1, "url-or-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URLOrID = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesDuplicateCmd() *cobra.Command {
	opts := &options.PagesDuplicateOptions{}

	cmd := &cobra.Command{
		Use:   "duplicate <url-or-id>",
		Short: "Duplicate a page",
		Long: `Duplicate a wiki page.

Examples:
  canvas pages duplicate --course-id 123 my-page-url`,
		Args: ExactArgsWithUsage(1, "url-or-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URLOrID = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesDuplicate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesRevisionsCmd() *cobra.Command {
	opts := &options.PagesRevisionsOptions{}

	cmd := &cobra.Command{
		Use:   "revisions <url-or-id>",
		Short: "List page revisions",
		Long: `List all revisions for a wiki page.

Examples:
  canvas pages revisions --course-id 123 my-page-url`,
		Args: ExactArgsWithUsage(1, "url-or-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URLOrID = args[0]
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesRevisions(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newPagesRevertCmd() *cobra.Command {
	opts := &options.PagesRevertOptions{}

	cmd := &cobra.Command{
		Use:   "revert <url-or-id> <revision-id>",
		Short: "Revert to a specific revision",
		Long: `Revert a wiki page to a specific revision.

Examples:
  canvas pages revert --course-id 123 my-page-url 5`,
		Args: ExactArgsWithUsage(2, "url-or-id", "revision-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.URLOrID = args[0]
			revisionID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid revision ID: %s", args[1])
			}
			opts.RevisionID = revisionID
			if err := opts.Validate(); err != nil {
				return err
			}
			client, err := getAPIClient()
			if err != nil {
				return err
			}
			return runPagesRevert(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func runPagesList(ctx context.Context, client *api.Client, opts *options.PagesListOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.list", map[string]interface{}{
		"course_id": opts.CourseID,
		"sort":      opts.Sort,
	})

	pagesService := api.NewPagesService(client)

	apiOpts := &api.ListPagesOptions{
		Sort:       opts.Sort,
		Order:      opts.Order,
		SearchTerm: opts.SearchTerm,
		Include:    opts.Include,
	}

	if opts.Published != "" {
		pub := opts.Published == "true"
		apiOpts.Published = &pub
	}

	pages, err := pagesService.List(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "pages.list", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to list pages: %w", err)
	}

	if len(pages) == 0 {
		logger.LogCommandComplete(ctx, "pages.list", 0)
		fmt.Println("No pages found")
		return nil
	}

	printVerbose("Found %d pages:\n\n", len(pages))
	logger.LogCommandComplete(ctx, "pages.list", len(pages))
	return formatOutput(pages, nil)
}

func runPagesGet(ctx context.Context, client *api.Client, opts *options.PagesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.get", map[string]interface{}{
		"course_id": opts.CourseID,
		"url_or_id": opts.URLOrID,
	})

	pagesService := api.NewPagesService(client)

	page, err := pagesService.Get(ctx, opts.CourseID, opts.URLOrID)
	if err != nil {
		logger.LogCommandError(ctx, "pages.get", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"url_or_id": opts.URLOrID,
		})
		return fmt.Errorf("failed to get page: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.get", 1)
	return formatOutput(page, nil)
}

func runPagesFront(ctx context.Context, client *api.Client, opts *options.PagesFrontOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.front", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	pagesService := api.NewPagesService(client)

	page, err := pagesService.GetFrontPage(ctx, opts.CourseID)
	if err != nil {
		logger.LogCommandError(ctx, "pages.front", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get front page: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.front", 1)
	return formatOutput(page, nil)
}

func runPagesCreate(ctx context.Context, client *api.Client, opts *options.PagesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"title":     opts.Title,
	})

	pagesService := api.NewPagesService(client)

	params := &api.CreatePageParams{
		Title:          opts.Title,
		Body:           opts.Body,
		EditingRoles:   opts.EditingRoles,
		NotifyOfUpdate: opts.NotifyUpdate,
		Published:      opts.Published,
		FrontPage:      opts.FrontPage,
		PublishAt:      opts.PublishAt,
	}

	page, err := pagesService.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "pages.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"title":     opts.Title,
		})
		return fmt.Errorf("failed to create page: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.create", 1)
	return formatSuccessOutput(page, "Page created successfully!")
}

func runPagesUpdate(ctx context.Context, client *api.Client, opts *options.PagesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.update", map[string]interface{}{
		"course_id": opts.CourseID,
		"url_or_id": opts.URLOrID,
	})

	pagesService := api.NewPagesService(client)

	params := &api.UpdatePageParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.BodySet {
		params.Body = &opts.Body
	}
	if opts.EditingRolesSet {
		params.EditingRoles = &opts.EditingRoles
	}
	if opts.NotifyUpdateSet {
		params.NotifyOfUpdate = &opts.NotifyUpdate
	}
	if opts.PublishedSet {
		params.Published = &opts.Published
	}
	if opts.FrontPageSet {
		params.FrontPage = &opts.FrontPage
	}
	if opts.PublishAtSet {
		params.PublishAt = &opts.PublishAt
	}

	page, err := pagesService.Update(ctx, opts.CourseID, opts.URLOrID, params)
	if err != nil {
		logger.LogCommandError(ctx, "pages.update", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"url_or_id": opts.URLOrID,
		})
		return fmt.Errorf("failed to update page: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.update", 1)
	return formatSuccessOutput(page, "Page updated successfully!")
}

func runPagesDelete(ctx context.Context, client *api.Client, opts *options.PagesDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.delete", map[string]interface{}{
		"course_id": opts.CourseID,
		"url_or_id": opts.URLOrID,
		"force":     opts.Force,
	})

	// Confirm deletion
	confirmed, err := confirmDelete("page", opts.URLOrID, opts.Force)
	if err != nil {
		logger.LogCommandError(ctx, "pages.delete", err, map[string]interface{}{})
		return err
	}
	if !confirmed {
		logger.LogCommandComplete(ctx, "pages.delete", 0)
		fmt.Println("Delete cancelled")
		return nil
	}

	pagesService := api.NewPagesService(client)

	if err := pagesService.Delete(ctx, opts.CourseID, opts.URLOrID); err != nil {
		logger.LogCommandError(ctx, "pages.delete", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"url_or_id": opts.URLOrID,
		})
		return fmt.Errorf("failed to delete page: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.delete", 1)
	fmt.Printf("Page '%s' deleted successfully\n", opts.URLOrID)
	return nil
}

func runPagesDuplicate(ctx context.Context, client *api.Client, opts *options.PagesDuplicateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.duplicate", map[string]interface{}{
		"course_id": opts.CourseID,
		"url_or_id": opts.URLOrID,
	})

	pagesService := api.NewPagesService(client)

	page, err := pagesService.Duplicate(ctx, opts.CourseID, opts.URLOrID)
	if err != nil {
		logger.LogCommandError(ctx, "pages.duplicate", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"url_or_id": opts.URLOrID,
		})
		return fmt.Errorf("failed to duplicate page: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.duplicate", 1)
	return formatSuccessOutput(page, "Page duplicated successfully!")
}

func runPagesRevisions(ctx context.Context, client *api.Client, opts *options.PagesRevisionsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.revisions", map[string]interface{}{
		"course_id": opts.CourseID,
		"url_or_id": opts.URLOrID,
	})

	pagesService := api.NewPagesService(client)

	revisions, err := pagesService.ListRevisions(ctx, opts.CourseID, opts.URLOrID)
	if err != nil {
		logger.LogCommandError(ctx, "pages.revisions", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"url_or_id": opts.URLOrID,
		})
		return fmt.Errorf("failed to list revisions: %w", err)
	}

	if len(revisions) == 0 {
		logger.LogCommandComplete(ctx, "pages.revisions", 0)
		fmt.Println("No revisions found")
		return nil
	}

	printVerbose("Found %d revisions:\n\n", len(revisions))
	logger.LogCommandComplete(ctx, "pages.revisions", len(revisions))
	return formatOutput(revisions, nil)
}

func runPagesRevert(ctx context.Context, client *api.Client, opts *options.PagesRevertOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "pages.revert", map[string]interface{}{
		"course_id":   opts.CourseID,
		"url_or_id":   opts.URLOrID,
		"revision_id": opts.RevisionID,
	})

	pagesService := api.NewPagesService(client)

	revision, err := pagesService.RevertToRevision(ctx, opts.CourseID, opts.URLOrID, opts.RevisionID)
	if err != nil {
		logger.LogCommandError(ctx, "pages.revert", err, map[string]interface{}{
			"course_id":   opts.CourseID,
			"url_or_id":   opts.URLOrID,
			"revision_id": opts.RevisionID,
		})
		return fmt.Errorf("failed to revert to revision: %w", err)
	}

	logger.LogCommandComplete(ctx, "pages.revert", 1)
	return formatSuccessOutput(revision, "Page reverted successfully!")
}
