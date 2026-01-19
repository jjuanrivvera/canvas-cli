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

// rubricsCmd represents the rubrics command group
var rubricsCmd = &cobra.Command{
	Use:   "rubrics",
	Short: "Manage Canvas rubrics",
	Long: `Manage Canvas rubrics for grading assignments.

Rubrics provide standardized criteria for grading and can be attached to
assignments to ensure consistent evaluation.

Examples:
  canvas rubrics list --course-id 123
  canvas rubrics get 456 --course-id 123
  canvas rubrics create --course-id 123 --title "Essay Rubric"`,
}

func init() {
	rootCmd.AddCommand(rubricsCmd)
	rubricsCmd.AddCommand(newRubricsListCmd())
	rubricsCmd.AddCommand(newRubricsGetCmd())
	rubricsCmd.AddCommand(newRubricsCreateCmd())
	rubricsCmd.AddCommand(newRubricsUpdateCmd())
	rubricsCmd.AddCommand(newRubricsDeleteCmd())
	rubricsCmd.AddCommand(newRubricsAssociateCmd())
}

func newRubricsListCmd() *cobra.Command {
	opts := &options.RubricsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List rubrics",
		Long: `List all rubrics in a course or account.

If neither --account-id nor --course-id is specified, uses default account.

Examples:
  canvas rubrics list                          # Uses default account
  canvas rubrics list --course-id 123
  canvas rubrics list --account-id 1
  canvas rubrics list --course-id 123 --include assessments,associations`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRubricsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data (assessments, associations, assignment_associations)")

	return cmd
}

func newRubricsGetCmd() *cobra.Command {
	opts := &options.RubricsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <rubric-id>",
		Short: "Get rubric details",
		Long: `Get details of a specific rubric.

Examples:
  canvas rubrics get 456 --course-id 123
  canvas rubrics get 456 --account-id 1 --include assessments`,
		Args: ExactArgsWithUsage(1, "rubric-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			rubricID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rubric ID: %w", err)
			}
			opts.RubricID = rubricID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRubricsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data")

	return cmd
}

func newRubricsCreateCmd() *cobra.Command {
	opts := &options.RubricsCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new rubric",
		Long: `Create a new rubric in a course.

Examples:
  canvas rubrics create --course-id 123 --title "Essay Rubric" --points 100
  canvas rubrics create --course-id 123 --title "Discussion Rubric" --free-form`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRubricsCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Rubric title (required)")
	cmd.MarkFlagRequired("title")
	cmd.Flags().Float64Var(&opts.PointsPossible, "points", 0, "Total points possible")
	cmd.Flags().BoolVar(&opts.FreeFormCriterionComments, "free-form", false, "Allow free-form criterion comments")
	cmd.Flags().BoolVar(&opts.HideScoreTotal, "hide-score-total", false, "Hide score total")

	return cmd
}

func newRubricsUpdateCmd() *cobra.Command {
	opts := &options.RubricsUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <rubric-id>",
		Short: "Update a rubric",
		Long: `Update an existing rubric.

Examples:
  canvas rubrics update 456 --course-id 123 --title "Updated Title"
  canvas rubrics update 456 --course-id 123 --points 150`,
		Args: ExactArgsWithUsage(1, "rubric-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			rubricID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rubric ID: %w", err)
			}
			opts.RubricID = rubricID

			// Track which fields were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.PointsPossibleSet = cmd.Flags().Changed("points")
			opts.FreeFormCriterionCommentsSet = cmd.Flags().Changed("free-form")
			opts.HideScoreTotalSet = cmd.Flags().Changed("hide-score-total")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRubricsUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Rubric title")
	cmd.Flags().Float64Var(&opts.PointsPossible, "points", 0, "Total points possible")
	cmd.Flags().BoolVar(&opts.FreeFormCriterionComments, "free-form", false, "Allow free-form criterion comments")
	cmd.Flags().BoolVar(&opts.HideScoreTotal, "hide-score-total", false, "Hide score total")

	return cmd
}

func newRubricsDeleteCmd() *cobra.Command {
	opts := &options.RubricsDeleteOptions{}

	cmd := &cobra.Command{
		Use:   "delete <rubric-id>",
		Short: "Delete a rubric",
		Long: `Delete a rubric.

Examples:
  canvas rubrics delete 456 --course-id 123
  canvas rubrics delete 456 --course-id 123 --force`,
		Args: ExactArgsWithUsage(1, "rubric-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			rubricID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rubric ID: %w", err)
			}
			opts.RubricID = rubricID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRubricsDelete(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")

	return cmd
}

func newRubricsAssociateCmd() *cobra.Command {
	opts := &options.RubricsAssociateOptions{}

	cmd := &cobra.Command{
		Use:   "associate <rubric-id>",
		Short: "Associate rubric with assignment",
		Long: `Associate a rubric with an assignment.

Examples:
  canvas rubrics associate 456 --course-id 123 --assignment-id 789
  canvas rubrics associate 456 --course-id 123 --assignment-id 789 --use-for-grading`,
		Args: ExactArgsWithUsage(1, "rubric-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			rubricID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid rubric ID: %w", err)
			}
			opts.RubricID = rubricID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runRubricsAssociate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")
	cmd.Flags().Int64Var(&opts.AssignmentID, "assignment-id", 0, "Assignment ID (required)")
	cmd.MarkFlagRequired("assignment-id")
	cmd.Flags().BoolVar(&opts.UseForGrading, "use-for-grading", false, "Use rubric for grading")
	cmd.Flags().BoolVar(&opts.HideScoreTotal, "hide-score-total", false, "Hide score total")
	cmd.Flags().BoolVar(&opts.HidePoints, "hide-points", false, "Hide points")

	return cmd
}

func runRubricsList(ctx context.Context, client *api.Client, opts *options.RubricsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	// Use default account ID if neither course nor account is specified
	if opts.CourseID == 0 && opts.AccountID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			return fmt.Errorf("must specify --course-id or --account-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		opts.AccountID = defaultID
		printVerbose("Using default account ID: %d\n", defaultID)
	}

	logger.LogCommandStart(ctx, "rubrics.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"include":    opts.Include,
	})

	service := api.NewRubricsService(client)

	apiOpts := &api.ListRubricsOptions{
		Include: opts.Include,
	}

	var rubrics []api.Rubric
	var err error

	if opts.CourseID > 0 {
		rubrics, err = service.ListCourse(ctx, opts.CourseID, apiOpts)
	} else {
		rubrics, err = service.ListAccount(ctx, opts.AccountID, apiOpts)
	}

	if err != nil {
		logger.LogCommandError(ctx, "rubrics.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list rubrics: %w", err)
	}

	if len(rubrics) == 0 {
		fmt.Println("No rubrics found")
		logger.LogCommandComplete(ctx, "rubrics.list", 0)
		return nil
	}

	printVerbose("Found %d rubrics:\n\n", len(rubrics))
	logger.LogCommandComplete(ctx, "rubrics.list", len(rubrics))
	return formatOutput(rubrics, nil)
}

func runRubricsGet(ctx context.Context, client *api.Client, opts *options.RubricsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubrics.get", map[string]interface{}{
		"rubric_id":  opts.RubricID,
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"include":    opts.Include,
	})

	service := api.NewRubricsService(client)

	var rubric *api.Rubric
	var err error

	if opts.CourseID > 0 {
		rubric, err = service.GetCourse(ctx, opts.CourseID, opts.RubricID, opts.Include)
	} else {
		rubric, err = service.GetAccount(ctx, opts.AccountID, opts.RubricID, opts.Include)
	}

	if err != nil {
		logger.LogCommandError(ctx, "rubrics.get", err, map[string]interface{}{
			"rubric_id":  opts.RubricID,
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to get rubric: %w", err)
	}

	logger.LogCommandComplete(ctx, "rubrics.get", 1)
	return formatOutput(rubric, nil)
}

func runRubricsCreate(ctx context.Context, client *api.Client, opts *options.RubricsCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubrics.create", map[string]interface{}{
		"course_id": opts.CourseID,
		"title":     opts.Title,
		"points":    opts.PointsPossible,
	})

	service := api.NewRubricsService(client)

	params := &api.CreateRubricParams{
		Title:                     opts.Title,
		PointsPossible:            opts.PointsPossible,
		FreeFormCriterionComments: opts.FreeFormCriterionComments,
		HideScoreTotal:            opts.HideScoreTotal,
	}

	rubric, err := service.Create(ctx, opts.CourseID, params)
	if err != nil {
		logger.LogCommandError(ctx, "rubrics.create", err, map[string]interface{}{
			"course_id": opts.CourseID,
			"title":     opts.Title,
		})
		return fmt.Errorf("failed to create rubric: %w", err)
	}

	fmt.Printf("Rubric created successfully (ID: %d)\n", rubric.ID)
	logger.LogCommandComplete(ctx, "rubrics.create", 1)
	return formatOutput(rubric, nil)
}

func runRubricsUpdate(ctx context.Context, client *api.Client, opts *options.RubricsUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubrics.update", map[string]interface{}{
		"rubric_id": opts.RubricID,
		"course_id": opts.CourseID,
	})

	service := api.NewRubricsService(client)

	params := &api.UpdateRubricParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.PointsPossibleSet {
		params.PointsPossible = &opts.PointsPossible
	}
	if opts.FreeFormCriterionCommentsSet {
		params.FreeFormCriterionComments = &opts.FreeFormCriterionComments
	}
	if opts.HideScoreTotalSet {
		params.HideScoreTotal = &opts.HideScoreTotal
	}

	rubric, err := service.Update(ctx, opts.CourseID, opts.RubricID, params)
	if err != nil {
		logger.LogCommandError(ctx, "rubrics.update", err, map[string]interface{}{
			"rubric_id": opts.RubricID,
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to update rubric: %w", err)
	}

	fmt.Printf("Rubric updated successfully (ID: %d)\n", rubric.ID)
	logger.LogCommandComplete(ctx, "rubrics.update", 1)
	return formatOutput(rubric, nil)
}

func runRubricsDelete(ctx context.Context, client *api.Client, opts *options.RubricsDeleteOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubrics.delete", map[string]interface{}{
		"rubric_id": opts.RubricID,
		"course_id": opts.CourseID,
		"force":     opts.Force,
	})

	if !opts.Force {
		fmt.Printf("WARNING: This will delete rubric %d.\n", opts.RubricID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			logger.LogCommandComplete(ctx, "rubrics.delete", 0)
			return nil
		}
	}

	service := api.NewRubricsService(client)

	rubric, err := service.Delete(ctx, opts.CourseID, opts.RubricID)
	if err != nil {
		logger.LogCommandError(ctx, "rubrics.delete", err, map[string]interface{}{
			"rubric_id": opts.RubricID,
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to delete rubric: %w", err)
	}

	fmt.Printf("Rubric %d deleted\n", rubric.ID)
	logger.LogCommandComplete(ctx, "rubrics.delete", 1)
	return nil
}

func runRubricsAssociate(ctx context.Context, client *api.Client, opts *options.RubricsAssociateOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "rubrics.associate", map[string]interface{}{
		"rubric_id":       opts.RubricID,
		"course_id":       opts.CourseID,
		"assignment_id":   opts.AssignmentID,
		"use_for_grading": opts.UseForGrading,
	})

	service := api.NewRubricsService(client)

	params := &api.AssociateParams{
		AssociationType: "Assignment",
		AssociationID:   opts.AssignmentID,
		UseForGrading:   opts.UseForGrading,
		HideScoreTotal:  opts.HideScoreTotal,
		HidePoints:      opts.HidePoints,
		Purpose:         "grading",
	}

	association, err := service.Associate(ctx, opts.CourseID, opts.RubricID, params)
	if err != nil {
		logger.LogCommandError(ctx, "rubrics.associate", err, map[string]interface{}{
			"rubric_id":     opts.RubricID,
			"course_id":     opts.CourseID,
			"assignment_id": opts.AssignmentID,
		})
		return fmt.Errorf("failed to associate rubric: %w", err)
	}

	fmt.Printf("Rubric %d associated with assignment %d (Association ID: %d)\n", opts.RubricID, opts.AssignmentID, association.ID)
	logger.LogCommandComplete(ctx, "rubrics.associate", 1)
	return formatOutput(association, nil)
}
