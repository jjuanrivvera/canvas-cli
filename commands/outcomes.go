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

// outcomesCmd represents the outcomes command group
var outcomesCmd = &cobra.Command{
	Use:   "outcomes",
	Short: "Manage Canvas learning outcomes",
	Long: `Manage Canvas learning outcomes.

Learning outcomes define what students should know or be able to do.
They can be organized into groups and linked to assignments, quizzes, and rubrics.

Examples:
  canvas outcomes get 123
  canvas outcomes groups list --account-id 1
  canvas outcomes results --course-id 123`,
}

// outcomesGroupsCmd represents the outcomes groups command group
var outcomesGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage outcome groups",
	Long: `Manage outcome groups.

Outcome groups organize learning outcomes into hierarchical structures.

Examples:
  canvas outcomes groups list --account-id 1
  canvas outcomes groups get 456 --account-id 1`,
}

func init() {
	rootCmd.AddCommand(outcomesCmd)
	outcomesCmd.AddCommand(newOutcomesGetCmd())
	outcomesCmd.AddCommand(newOutcomesCreateCmd())
	outcomesCmd.AddCommand(newOutcomesUpdateCmd())
	outcomesCmd.AddCommand(outcomesGroupsCmd)
	outcomesCmd.AddCommand(newOutcomesListCmd())
	outcomesCmd.AddCommand(newOutcomesLinkCmd())
	outcomesCmd.AddCommand(newOutcomesUnlinkCmd())
	outcomesCmd.AddCommand(newOutcomesResultsCmd())
	outcomesCmd.AddCommand(newOutcomesAlignmentsCmd())

	outcomesGroupsCmd.AddCommand(newOutcomesGroupsListCmd())
	outcomesGroupsCmd.AddCommand(newOutcomesGroupsGetCmd())
}

func newOutcomesGetCmd() *cobra.Command {
	opts := &options.OutcomesGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <outcome-id>",
		Short: "Get outcome details",
		Long: `Get details of a specific learning outcome.

Examples:
  canvas outcomes get 123`,
		Args: ExactArgsWithUsage(1, "outcome-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			outcomeID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid outcome ID: %s", args[0])
			}
			opts.OutcomeID = outcomeID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesGet(cmd.Context(), client, opts)
		},
	}

	return cmd
}

func newOutcomesCreateCmd() *cobra.Command {
	opts := &options.OutcomesCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new outcome",
		Long: `Create a new learning outcome in an outcome group.

Outcomes must be created within an outcome group. Use --group-id to specify
the target group, which can be the root group or any subgroup.

Calculation methods:
  - decaying_average: Weighted average favoring recent scores
  - n_mastery: Latest N scores must meet mastery threshold
  - latest: Only the most recent score counts
  - highest: Use the highest score achieved

Examples:
  canvas outcomes create --course-id 123 --group-id 456 --title "Problem Solving"
  canvas outcomes create --account-id 1 --group-id 789 --title "Critical Thinking" --mastery-points 4
  canvas outcomes create --course-id 123 --group-id 456 --title "Writing" --calculation-method decaying_average --calculation-int 65`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesCreate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Outcome group ID (required)")
	cmd.Flags().StringVar(&opts.Title, "title", "", "Outcome title (required)")
	cmd.Flags().StringVar(&opts.DisplayName, "display-name", "", "Display name")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Description")
	cmd.Flags().Float64Var(&opts.MasteryPoints, "mastery-points", 0, "Points for mastery")
	cmd.Flags().StringVar(&opts.CalculationMethod, "calculation-method", "", "Calculation method (decaying_average, n_mastery, latest, highest)")
	cmd.Flags().IntVar(&opts.CalculationInt, "calculation-int", 0, "Calculation parameter")
	cmd.MarkFlagRequired("group-id")
	cmd.MarkFlagRequired("title")

	return cmd
}

func newOutcomesUpdateCmd() *cobra.Command {
	opts := &options.OutcomesUpdateOptions{}

	cmd := &cobra.Command{
		Use:   "update <outcome-id>",
		Short: "Update an outcome",
		Long: `Update an existing learning outcome.

Examples:
  canvas outcomes update 123 --title "Updated Outcome"
  canvas outcomes update 123 --mastery-points 4
  canvas outcomes update 123 --calculation-method decaying_average --calculation-int 65`,
		Args: ExactArgsWithUsage(1, "outcome-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			outcomeID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid outcome ID: %s", args[0])
			}
			opts.OutcomeID = outcomeID

			// Track which flags were set
			opts.TitleSet = cmd.Flags().Changed("title")
			opts.DisplayNameSet = cmd.Flags().Changed("display-name")
			opts.DescriptionSet = cmd.Flags().Changed("description")
			opts.MasteryPointsSet = cmd.Flags().Changed("mastery-points")
			opts.CalculationMethodSet = cmd.Flags().Changed("calculation-method")
			opts.CalculationIntSet = cmd.Flags().Changed("calculation-int")

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesUpdate(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Title, "title", "", "Outcome title")
	cmd.Flags().StringVar(&opts.DisplayName, "display-name", "", "Display name")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Description")
	cmd.Flags().Float64Var(&opts.MasteryPoints, "mastery-points", 0, "Points for mastery")
	cmd.Flags().StringVar(&opts.CalculationMethod, "calculation-method", "", "Calculation method (decaying_average, n_mastery, latest, highest)")
	cmd.Flags().IntVar(&opts.CalculationInt, "calculation-int", 0, "Calculation parameter")

	return cmd
}

func newOutcomesListCmd() *cobra.Command {
	opts := &options.OutcomesListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List outcomes in a group",
		Long: `List all outcomes in a specific outcome group.

Examples:
  canvas outcomes list --account-id 1 --group-id 456
  canvas outcomes list --course-id 123 --group-id 456`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Outcome group ID (required)")
	cmd.MarkFlagRequired("group-id")

	return cmd
}

func newOutcomesLinkCmd() *cobra.Command {
	opts := &options.OutcomesLinkOptions{}

	cmd := &cobra.Command{
		Use:   "link <outcome-id>",
		Short: "Link outcome to a group",
		Long: `Link an existing outcome to an outcome group.

Examples:
  canvas outcomes link 789 --account-id 1 --group-id 456
  canvas outcomes link 789 --course-id 123 --group-id 456`,
		Args: ExactArgsWithUsage(1, "outcome-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			outcomeID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid outcome ID: %s", args[0])
			}
			opts.OutcomeID = outcomeID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesLink(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Outcome group ID (required)")
	cmd.MarkFlagRequired("group-id")

	return cmd
}

func newOutcomesUnlinkCmd() *cobra.Command {
	opts := &options.OutcomesUnlinkOptions{}

	cmd := &cobra.Command{
		Use:   "unlink <outcome-id>",
		Short: "Unlink outcome from a group",
		Long: `Remove an outcome link from an outcome group.

Examples:
  canvas outcomes unlink 789 --account-id 1 --group-id 456
  canvas outcomes unlink 789 --course-id 123 --group-id 456`,
		Args: ExactArgsWithUsage(1, "outcome-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			outcomeID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid outcome ID: %s", args[0])
			}
			opts.OutcomeID = outcomeID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesUnlink(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")
	cmd.Flags().Int64Var(&opts.GroupID, "group-id", 0, "Outcome group ID (required)")
	cmd.Flags().BoolVar(&opts.Force, "force", false, "Skip confirmation prompt")
	cmd.MarkFlagRequired("group-id")

	return cmd
}

func newOutcomesGroupsListCmd() *cobra.Command {
	opts := &options.OutcomesGroupsListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List outcome groups",
		Long: `List all outcome groups in a course or account.

If neither --account-id nor --course-id is specified, uses default account.

Examples:
  canvas outcomes groups list                  # Uses default account
  canvas outcomes groups list --account-id 1
  canvas outcomes groups list --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesGroupsList(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")

	return cmd
}

func newOutcomesGroupsGetCmd() *cobra.Command {
	opts := &options.OutcomesGroupsGetOptions{}

	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get outcome group details",
		Long: `Get details of a specific outcome group.

Examples:
  canvas outcomes groups get 456 --account-id 1
  canvas outcomes groups get 456 --course-id 123`,
		Args: ExactArgsWithUsage(1, "group-id"),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid group ID: %s", args[0])
			}
			opts.GroupID = groupID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesGroupsGet(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID")
	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID")

	return cmd
}

func newOutcomesResultsCmd() *cobra.Command {
	opts := &options.OutcomesResultsOptions{}

	cmd := &cobra.Command{
		Use:   "results",
		Short: "Get outcome results",
		Long: `Get student outcome results for a course.

Examples:
  canvas outcomes results --course-id 123
  canvas outcomes results --course-id 123 --user-ids 100,101
  canvas outcomes results --course-id 123 --outcome-ids 200,201
  canvas outcomes results --course-id 123 --include outcomes,alignments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesResults(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().Int64SliceVar(&opts.UserIDs, "user-ids", []int64{}, "Filter by user IDs")
	cmd.Flags().Int64SliceVar(&opts.OutcomeIDs, "outcome-ids", []int64{}, "Filter by outcome IDs")
	cmd.Flags().StringSliceVar(&opts.Include, "include", []string{}, "Include additional data (alignments, outcomes, outcomes.alignments, outcome_groups)")
	cmd.Flags().BoolVar(&opts.IncludeHidden, "include-hidden", false, "Include hidden outcomes")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newOutcomesAlignmentsCmd() *cobra.Command {
	opts := &options.OutcomesAlignmentsOptions{}

	cmd := &cobra.Command{
		Use:   "alignments",
		Short: "Get outcome alignments",
		Long: `Get outcome alignments for a course.

Shows which assignments, quizzes, etc. are aligned to outcomes.

Examples:
  canvas outcomes alignments --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runOutcomesAlignments(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

// Run functions

func runOutcomesGet(ctx context.Context, client *api.Client, opts *options.OutcomesGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.get", map[string]interface{}{
		"outcome_id": opts.OutcomeID,
	})

	service := api.NewOutcomesService(client)

	outcome, err := service.Get(ctx, opts.OutcomeID)
	if err != nil {
		logger.LogCommandError(ctx, "outcomes.get", err, map[string]interface{}{
			"outcome_id": opts.OutcomeID,
		})
		return fmt.Errorf("failed to get outcome: %w", err)
	}

	if err := formatOutput(outcome, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.get", 1)
	return nil
}

func runOutcomesCreate(ctx context.Context, client *api.Client, opts *options.OutcomesCreateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.create", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"group_id":   opts.GroupID,
		"title":      opts.Title,
	})

	service := api.NewOutcomesService(client)

	params := &api.CreateOutcomeParams{
		Title:             opts.Title,
		DisplayName:       opts.DisplayName,
		Description:       opts.Description,
		MasteryPoints:     opts.MasteryPoints,
		CalculationMethod: opts.CalculationMethod,
		CalculationInt:    opts.CalculationInt,
	}

	var link *api.OutcomeLink
	var err error

	if opts.CourseID > 0 {
		link, err = service.CreateOutcomeCourse(ctx, opts.CourseID, opts.GroupID, params)
	} else {
		link, err = service.CreateOutcomeAccount(ctx, opts.AccountID, opts.GroupID, params)
	}

	if err != nil {
		logger.LogCommandError(ctx, "outcomes.create", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"group_id":   opts.GroupID,
		})
		return fmt.Errorf("failed to create outcome: %w", err)
	}

	if link.Outcome != nil {
		printInfo("Outcome created successfully (ID: %d)\n", link.Outcome.ID)
	} else {
		printInfo("Outcome created successfully\n")
	}

	if err := formatOutput(link, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.create", 1)
	return nil
}

func runOutcomesUpdate(ctx context.Context, client *api.Client, opts *options.OutcomesUpdateOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.update", map[string]interface{}{
		"outcome_id": opts.OutcomeID,
	})

	service := api.NewOutcomesService(client)

	params := &api.UpdateOutcomeParams{}

	if opts.TitleSet {
		params.Title = &opts.Title
	}
	if opts.DisplayNameSet {
		params.DisplayName = &opts.DisplayName
	}
	if opts.DescriptionSet {
		params.Description = &opts.Description
	}
	if opts.MasteryPointsSet {
		params.MasteryPoints = &opts.MasteryPoints
	}
	if opts.CalculationMethodSet {
		params.CalculationMethod = &opts.CalculationMethod
	}
	if opts.CalculationIntSet {
		params.CalculationInt = &opts.CalculationInt
	}

	outcome, err := service.Update(ctx, opts.OutcomeID, params)
	if err != nil {
		logger.LogCommandError(ctx, "outcomes.update", err, map[string]interface{}{
			"outcome_id": opts.OutcomeID,
		})
		return fmt.Errorf("failed to update outcome: %w", err)
	}

	printInfo("Outcome updated successfully (ID: %d)\n", outcome.ID)
	if err := formatOutput(outcome, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.update", 1)
	return nil
}

func runOutcomesList(ctx context.Context, client *api.Client, opts *options.OutcomesListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"group_id":   opts.GroupID,
	})

	service := api.NewOutcomesService(client)

	var links []api.OutcomeLink
	var err error

	if opts.CourseID > 0 {
		links, err = service.ListOutcomesInGroupCourse(ctx, opts.CourseID, opts.GroupID, nil)
	} else {
		links, err = service.ListOutcomesInGroupAccount(ctx, opts.AccountID, opts.GroupID, nil)
	}

	if err != nil {
		logger.LogCommandError(ctx, "outcomes.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"group_id":   opts.GroupID,
		})
		return fmt.Errorf("failed to list outcomes: %w", err)
	}

	if err := formatEmptyOrOutput(links, "No outcomes found in group"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.list", len(links))
	return nil
}

func runOutcomesLink(ctx context.Context, client *api.Client, opts *options.OutcomesLinkOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.link", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"group_id":   opts.GroupID,
		"outcome_id": opts.OutcomeID,
	})

	service := api.NewOutcomesService(client)

	var link *api.OutcomeLink
	var err error

	if opts.CourseID > 0 {
		link, err = service.LinkOutcomeCourse(ctx, opts.CourseID, opts.GroupID, opts.OutcomeID)
	} else {
		link, err = service.LinkOutcomeAccount(ctx, opts.AccountID, opts.GroupID, opts.OutcomeID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "outcomes.link", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"group_id":   opts.GroupID,
			"outcome_id": opts.OutcomeID,
		})
		return fmt.Errorf("failed to link outcome: %w", err)
	}

	fmt.Printf("Outcome %d linked to group %d\n", opts.OutcomeID, opts.GroupID)
	if err := formatOutput(link, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.link", 1)
	return nil
}

func runOutcomesUnlink(ctx context.Context, client *api.Client, opts *options.OutcomesUnlinkOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.unlink", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"group_id":   opts.GroupID,
		"outcome_id": opts.OutcomeID,
	})

	if !opts.Force {
		fmt.Printf("WARNING: This will unlink outcome %d from group %d.\n", opts.OutcomeID, opts.GroupID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Unlink cancelled")
			return nil
		}
	}

	service := api.NewOutcomesService(client)

	var link *api.OutcomeLink
	var err error

	if opts.CourseID > 0 {
		link, err = service.UnlinkOutcomeCourse(ctx, opts.CourseID, opts.GroupID, opts.OutcomeID)
	} else {
		link, err = service.UnlinkOutcomeAccount(ctx, opts.AccountID, opts.GroupID, opts.OutcomeID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "outcomes.unlink", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"group_id":   opts.GroupID,
			"outcome_id": opts.OutcomeID,
		})
		return fmt.Errorf("failed to unlink outcome: %w", err)
	}

	fmt.Printf("Outcome %d unlinked from group %d\n", opts.OutcomeID, opts.GroupID)
	if err := formatOutput(link, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.unlink", 1)
	return nil
}

func runOutcomesGroupsList(ctx context.Context, client *api.Client, opts *options.OutcomesGroupsListOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.groups.list", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
	})

	// Use default account ID if neither course nor account is specified
	if opts.CourseID == 0 && opts.AccountID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			logger.LogCommandError(ctx, "outcomes.groups.list", err, map[string]interface{}{})
			return fmt.Errorf("must specify --course-id or --account-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		opts.AccountID = defaultID
		printVerbose("Using default account ID: %d\n", defaultID)
	}

	service := api.NewOutcomesService(client)

	var groups []api.OutcomeGroup
	var err error

	if opts.CourseID > 0 {
		groups, err = service.ListGroupsCourse(ctx, opts.CourseID, nil)
	} else {
		groups, err = service.ListGroupsAccount(ctx, opts.AccountID, nil)
	}

	if err != nil {
		logger.LogCommandError(ctx, "outcomes.groups.list", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
		})
		return fmt.Errorf("failed to list outcome groups: %w", err)
	}

	if err := formatEmptyOrOutput(groups, "No outcome groups found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.groups.list", len(groups))
	return nil
}

func runOutcomesGroupsGet(ctx context.Context, client *api.Client, opts *options.OutcomesGroupsGetOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.groups.get", map[string]interface{}{
		"course_id":  opts.CourseID,
		"account_id": opts.AccountID,
		"group_id":   opts.GroupID,
	})

	service := api.NewOutcomesService(client)

	var group *api.OutcomeGroup
	var err error

	if opts.CourseID > 0 {
		group, err = service.GetGroupCourse(ctx, opts.CourseID, opts.GroupID)
	} else {
		group, err = service.GetGroupAccount(ctx, opts.AccountID, opts.GroupID)
	}

	if err != nil {
		logger.LogCommandError(ctx, "outcomes.groups.get", err, map[string]interface{}{
			"course_id":  opts.CourseID,
			"account_id": opts.AccountID,
			"group_id":   opts.GroupID,
		})
		return fmt.Errorf("failed to get outcome group: %w", err)
	}

	if err := formatOutput(group, nil); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.groups.get", 1)
	return nil
}

func runOutcomesResults(ctx context.Context, client *api.Client, opts *options.OutcomesResultsOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.results", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewOutcomesService(client)

	apiOpts := &api.OutcomeResultsOptions{
		UserIDs:       opts.UserIDs,
		OutcomeIDs:    opts.OutcomeIDs,
		Include:       opts.Include,
		IncludeHidden: opts.IncludeHidden,
	}

	response, err := service.GetResults(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "outcomes.results", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get outcome results: %w", err)
	}

	if err := formatEmptyOrOutput(response, "No outcome results found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.results", len(response.OutcomeResults))
	return nil
}

func runOutcomesAlignments(ctx context.Context, client *api.Client, opts *options.OutcomesAlignmentsOptions) error {
	logger := logging.NewCommandLogger(verbose)

	logger.LogCommandStart(ctx, "outcomes.alignments", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewOutcomesService(client)

	alignments, err := service.GetAlignments(ctx, opts.CourseID, 0)
	if err != nil {
		logger.LogCommandError(ctx, "outcomes.alignments", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get outcome alignments: %w", err)
	}

	if err := formatEmptyOrOutput(alignments, "No outcome alignments found"); err != nil {
		return fmt.Errorf("failed to print results: %w", err)
	}

	logger.LogCommandComplete(ctx, "outcomes.alignments", len(alignments))
	return nil
}
