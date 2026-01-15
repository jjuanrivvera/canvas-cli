package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	// Common flags
	outcomesCourseID  int64
	outcomesAccountID int64
	outcomesGroupID   int64

	// Update flags
	outcomesTitle             string
	outcomesDisplayName       string
	outcomesDescription       string
	outcomesMasteryPoints     float64
	outcomesCalculationMethod string
	outcomesCalculationInt    int

	// Results flags
	outcomesUserIDs       []int64
	outcomesOutcomeIDs    []int64
	outcomesInclude       []string
	outcomesIncludeHidden bool

	// Delete flags
	outcomesForce bool
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

// outcomesGetCmd represents the outcomes get command
var outcomesGetCmd = &cobra.Command{
	Use:   "get <outcome-id>",
	Short: "Get outcome details",
	Long: `Get details of a specific learning outcome.

Examples:
  canvas outcomes get 123`,
	Args: ExactArgsWithUsage(1, "outcome-id"),
	RunE: runOutcomesGet,
}

// outcomesCreateCmd represents the outcomes create command
var outcomesCreateCmd = &cobra.Command{
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
	RunE: runOutcomesCreate,
}

// outcomesUpdateCmd represents the outcomes update command
var outcomesUpdateCmd = &cobra.Command{
	Use:   "update <outcome-id>",
	Short: "Update an outcome",
	Long: `Update an existing learning outcome.

Examples:
  canvas outcomes update 123 --title "Updated Outcome"
  canvas outcomes update 123 --mastery-points 4
  canvas outcomes update 123 --calculation-method decaying_average --calculation-int 65`,
	Args: ExactArgsWithUsage(1, "outcome-id"),
	RunE: runOutcomesUpdate,
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

// outcomesGroupsListCmd represents the outcomes groups list command
var outcomesGroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List outcome groups",
	Long: `List all outcome groups in a course or account.

If neither --account-id nor --course-id is specified, uses default account.

Examples:
  canvas outcomes groups list                  # Uses default account
  canvas outcomes groups list --account-id 1
  canvas outcomes groups list --course-id 123`,
	RunE: runOutcomesGroupsList,
}

// outcomesGroupsGetCmd represents the outcomes groups get command
var outcomesGroupsGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get outcome group details",
	Long: `Get details of a specific outcome group.

Examples:
  canvas outcomes groups get 456 --account-id 1
  canvas outcomes groups get 456 --course-id 123`,
	Args: ExactArgsWithUsage(1, "group-id"),
	RunE: runOutcomesGroupsGet,
}

// outcomesListCmd represents the outcomes list command
var outcomesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List outcomes in a group",
	Long: `List all outcomes in a specific outcome group.

Examples:
  canvas outcomes list --account-id 1 --group-id 456
  canvas outcomes list --course-id 123 --group-id 456`,
	RunE: runOutcomesList,
}

// outcomesLinkCmd represents the outcomes link command
var outcomesLinkCmd = &cobra.Command{
	Use:   "link <outcome-id>",
	Short: "Link outcome to a group",
	Long: `Link an existing outcome to an outcome group.

Examples:
  canvas outcomes link 789 --account-id 1 --group-id 456
  canvas outcomes link 789 --course-id 123 --group-id 456`,
	Args: ExactArgsWithUsage(1, "outcome-id"),
	RunE: runOutcomesLink,
}

// outcomesUnlinkCmd represents the outcomes unlink command
var outcomesUnlinkCmd = &cobra.Command{
	Use:   "unlink <outcome-id>",
	Short: "Unlink outcome from a group",
	Long: `Remove an outcome link from an outcome group.

Examples:
  canvas outcomes unlink 789 --account-id 1 --group-id 456
  canvas outcomes unlink 789 --course-id 123 --group-id 456`,
	Args: ExactArgsWithUsage(1, "outcome-id"),
	RunE: runOutcomesUnlink,
}

// outcomesResultsCmd represents the outcomes results command
var outcomesResultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Get outcome results",
	Long: `Get student outcome results for a course.

Examples:
  canvas outcomes results --course-id 123
  canvas outcomes results --course-id 123 --user-ids 100,101
  canvas outcomes results --course-id 123 --outcome-ids 200,201
  canvas outcomes results --course-id 123 --include outcomes,alignments`,
	RunE: runOutcomesResults,
}

// outcomesAlignmentsCmd represents the outcomes alignments command
var outcomesAlignmentsCmd = &cobra.Command{
	Use:   "alignments",
	Short: "Get outcome alignments",
	Long: `Get outcome alignments for a course.

Shows which assignments, quizzes, etc. are aligned to outcomes.

Examples:
  canvas outcomes alignments --course-id 123`,
	RunE: runOutcomesAlignments,
}

func init() {
	rootCmd.AddCommand(outcomesCmd)
	outcomesCmd.AddCommand(outcomesGetCmd)
	outcomesCmd.AddCommand(outcomesCreateCmd)
	outcomesCmd.AddCommand(outcomesUpdateCmd)
	outcomesCmd.AddCommand(outcomesGroupsCmd)
	outcomesCmd.AddCommand(outcomesListCmd)
	outcomesCmd.AddCommand(outcomesLinkCmd)
	outcomesCmd.AddCommand(outcomesUnlinkCmd)
	outcomesCmd.AddCommand(outcomesResultsCmd)
	outcomesCmd.AddCommand(outcomesAlignmentsCmd)

	outcomesGroupsCmd.AddCommand(outcomesGroupsListCmd)
	outcomesGroupsCmd.AddCommand(outcomesGroupsGetCmd)

	// Create flags
	outcomesCreateCmd.Flags().Int64Var(&outcomesAccountID, "account-id", 0, "Account ID")
	outcomesCreateCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID")
	outcomesCreateCmd.Flags().Int64Var(&outcomesGroupID, "group-id", 0, "Outcome group ID (required)")
	outcomesCreateCmd.Flags().StringVar(&outcomesTitle, "title", "", "Outcome title (required)")
	outcomesCreateCmd.Flags().StringVar(&outcomesDisplayName, "display-name", "", "Display name")
	outcomesCreateCmd.Flags().StringVar(&outcomesDescription, "description", "", "Description")
	outcomesCreateCmd.Flags().Float64Var(&outcomesMasteryPoints, "mastery-points", 0, "Points for mastery")
	outcomesCreateCmd.Flags().StringVar(&outcomesCalculationMethod, "calculation-method", "", "Calculation method (decaying_average, n_mastery, latest, highest)")
	outcomesCreateCmd.Flags().IntVar(&outcomesCalculationInt, "calculation-int", 0, "Calculation parameter")
	outcomesCreateCmd.MarkFlagRequired("group-id")
	outcomesCreateCmd.MarkFlagRequired("title")

	// Update flags
	outcomesUpdateCmd.Flags().StringVar(&outcomesTitle, "title", "", "Outcome title")
	outcomesUpdateCmd.Flags().StringVar(&outcomesDisplayName, "display-name", "", "Display name")
	outcomesUpdateCmd.Flags().StringVar(&outcomesDescription, "description", "", "Description")
	outcomesUpdateCmd.Flags().Float64Var(&outcomesMasteryPoints, "mastery-points", 0, "Points for mastery")
	outcomesUpdateCmd.Flags().StringVar(&outcomesCalculationMethod, "calculation-method", "", "Calculation method (decaying_average, n_mastery, latest, highest)")
	outcomesUpdateCmd.Flags().IntVar(&outcomesCalculationInt, "calculation-int", 0, "Calculation parameter")

	// Groups list flags
	outcomesGroupsListCmd.Flags().Int64Var(&outcomesAccountID, "account-id", 0, "Account ID")
	outcomesGroupsListCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID")

	// Groups get flags
	outcomesGroupsGetCmd.Flags().Int64Var(&outcomesAccountID, "account-id", 0, "Account ID")
	outcomesGroupsGetCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID")

	// List flags
	outcomesListCmd.Flags().Int64Var(&outcomesAccountID, "account-id", 0, "Account ID")
	outcomesListCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID")
	outcomesListCmd.Flags().Int64Var(&outcomesGroupID, "group-id", 0, "Outcome group ID (required)")
	outcomesListCmd.MarkFlagRequired("group-id")

	// Link flags
	outcomesLinkCmd.Flags().Int64Var(&outcomesAccountID, "account-id", 0, "Account ID")
	outcomesLinkCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID")
	outcomesLinkCmd.Flags().Int64Var(&outcomesGroupID, "group-id", 0, "Outcome group ID (required)")
	outcomesLinkCmd.MarkFlagRequired("group-id")

	// Unlink flags
	outcomesUnlinkCmd.Flags().Int64Var(&outcomesAccountID, "account-id", 0, "Account ID")
	outcomesUnlinkCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID")
	outcomesUnlinkCmd.Flags().Int64Var(&outcomesGroupID, "group-id", 0, "Outcome group ID (required)")
	outcomesUnlinkCmd.MarkFlagRequired("group-id")
	outcomesUnlinkCmd.Flags().BoolVar(&outcomesForce, "force", false, "Skip confirmation prompt")

	// Results flags
	outcomesResultsCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID (required)")
	outcomesResultsCmd.MarkFlagRequired("course-id")
	outcomesResultsCmd.Flags().Int64SliceVar(&outcomesUserIDs, "user-ids", []int64{}, "Filter by user IDs")
	outcomesResultsCmd.Flags().Int64SliceVar(&outcomesOutcomeIDs, "outcome-ids", []int64{}, "Filter by outcome IDs")
	outcomesResultsCmd.Flags().StringSliceVar(&outcomesInclude, "include", []string{}, "Include additional data (alignments, outcomes, outcomes.alignments, outcome_groups)")
	outcomesResultsCmd.Flags().BoolVar(&outcomesIncludeHidden, "include-hidden", false, "Include hidden outcomes")

	// Alignments flags
	outcomesAlignmentsCmd.Flags().Int64Var(&outcomesCourseID, "course-id", 0, "Course ID (required)")
	outcomesAlignmentsCmd.MarkFlagRequired("course-id")
}

func runOutcomesGet(cmd *cobra.Command, args []string) error {
	outcomeID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid outcome ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	outcome, err := service.Get(ctx, outcomeID)
	if err != nil {
		return fmt.Errorf("failed to get outcome: %w", err)
	}

	return formatOutput(outcome, nil)
}

func runOutcomesCreate(cmd *cobra.Command, args []string) error {
	if outcomesCourseID == 0 && outcomesAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	params := &api.CreateOutcomeParams{
		Title:             outcomesTitle,
		DisplayName:       outcomesDisplayName,
		Description:       outcomesDescription,
		MasteryPoints:     outcomesMasteryPoints,
		CalculationMethod: outcomesCalculationMethod,
		CalculationInt:    outcomesCalculationInt,
	}

	ctx := context.Background()
	var link *api.OutcomeLink

	if outcomesCourseID > 0 {
		link, err = service.CreateOutcomeCourse(ctx, outcomesCourseID, outcomesGroupID, params)
	} else {
		link, err = service.CreateOutcomeAccount(ctx, outcomesAccountID, outcomesGroupID, params)
	}

	if err != nil {
		return fmt.Errorf("failed to create outcome: %w", err)
	}

	fmt.Printf("Outcome created successfully")
	if link.Outcome != nil {
		fmt.Printf(" (ID: %d)", link.Outcome.ID)
	}
	fmt.Println()
	return formatOutput(link, nil)
}

func runOutcomesUpdate(cmd *cobra.Command, args []string) error {
	outcomeID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid outcome ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	params := &api.UpdateOutcomeParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &outcomesTitle
	}
	if cmd.Flags().Changed("display-name") {
		params.DisplayName = &outcomesDisplayName
	}
	if cmd.Flags().Changed("description") {
		params.Description = &outcomesDescription
	}
	if cmd.Flags().Changed("mastery-points") {
		params.MasteryPoints = &outcomesMasteryPoints
	}
	if cmd.Flags().Changed("calculation-method") {
		params.CalculationMethod = &outcomesCalculationMethod
	}
	if cmd.Flags().Changed("calculation-int") {
		params.CalculationInt = &outcomesCalculationInt
	}

	ctx := context.Background()
	outcome, err := service.Update(ctx, outcomeID, params)
	if err != nil {
		return fmt.Errorf("failed to update outcome: %w", err)
	}

	fmt.Printf("Outcome updated successfully (ID: %d)\n", outcome.ID)
	return formatOutput(outcome, nil)
}

func runOutcomesGroupsList(cmd *cobra.Command, args []string) error {
	// Use default account ID if neither course nor account is specified
	if outcomesCourseID == 0 && outcomesAccountID == 0 {
		defaultID, err := getDefaultAccountID()
		if err != nil || defaultID == 0 {
			return fmt.Errorf("must specify --course-id or --account-id (no default account configured). Use 'canvas config account --detect' to set one")
		}
		outcomesAccountID = defaultID
		printVerbose("Using default account ID: %d\n", defaultID)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	var groups []api.OutcomeGroup

	if outcomesCourseID > 0 {
		groups, err = service.ListGroupsCourse(ctx, outcomesCourseID, nil)
	} else {
		groups, err = service.ListGroupsAccount(ctx, outcomesAccountID, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to list outcome groups: %w", err)
	}

	if len(groups) == 0 {
		fmt.Println("No outcome groups found")
		return nil
	}

	printVerbose("Found %d outcome groups:\n\n", len(groups))
	return formatOutput(groups, nil)
}

func runOutcomesGroupsGet(cmd *cobra.Command, args []string) error {
	groupID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	if outcomesCourseID == 0 && outcomesAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	var group *api.OutcomeGroup

	if outcomesCourseID > 0 {
		group, err = service.GetGroupCourse(ctx, outcomesCourseID, groupID)
	} else {
		group, err = service.GetGroupAccount(ctx, outcomesAccountID, groupID)
	}

	if err != nil {
		return fmt.Errorf("failed to get outcome group: %w", err)
	}

	return formatOutput(group, nil)
}

func runOutcomesList(cmd *cobra.Command, args []string) error {
	if outcomesCourseID == 0 && outcomesAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	var links []api.OutcomeLink

	if outcomesCourseID > 0 {
		links, err = service.ListOutcomesInGroupCourse(ctx, outcomesCourseID, outcomesGroupID, nil)
	} else {
		links, err = service.ListOutcomesInGroupAccount(ctx, outcomesAccountID, outcomesGroupID, nil)
	}

	if err != nil {
		return fmt.Errorf("failed to list outcomes: %w", err)
	}

	if len(links) == 0 {
		fmt.Println("No outcomes found in group")
		return nil
	}

	printVerbose("Found %d outcomes:\n\n", len(links))
	return formatOutput(links, nil)
}

func runOutcomesLink(cmd *cobra.Command, args []string) error {
	outcomeID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid outcome ID: %w", err)
	}

	if outcomesCourseID == 0 && outcomesAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	var link *api.OutcomeLink

	if outcomesCourseID > 0 {
		link, err = service.LinkOutcomeCourse(ctx, outcomesCourseID, outcomesGroupID, outcomeID)
	} else {
		link, err = service.LinkOutcomeAccount(ctx, outcomesAccountID, outcomesGroupID, outcomeID)
	}

	if err != nil {
		return fmt.Errorf("failed to link outcome: %w", err)
	}

	fmt.Printf("Outcome %d linked to group %d\n", outcomeID, outcomesGroupID)
	return formatOutput(link, nil)
}

func runOutcomesUnlink(cmd *cobra.Command, args []string) error {
	outcomeID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid outcome ID: %w", err)
	}

	if outcomesCourseID == 0 && outcomesAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	if !outcomesForce {
		fmt.Printf("WARNING: This will unlink outcome %d from group %d.\n", outcomeID, outcomesGroupID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Unlink cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	var link *api.OutcomeLink

	if outcomesCourseID > 0 {
		link, err = service.UnlinkOutcomeCourse(ctx, outcomesCourseID, outcomesGroupID, outcomeID)
	} else {
		link, err = service.UnlinkOutcomeAccount(ctx, outcomesAccountID, outcomesGroupID, outcomeID)
	}

	if err != nil {
		return fmt.Errorf("failed to unlink outcome: %w", err)
	}

	fmt.Printf("Outcome %d unlinked from group %d\n", outcomeID, outcomesGroupID)
	return formatOutput(link, nil)
}

func runOutcomesResults(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	opts := &api.OutcomeResultsOptions{
		UserIDs:       outcomesUserIDs,
		OutcomeIDs:    outcomesOutcomeIDs,
		Include:       outcomesInclude,
		IncludeHidden: outcomesIncludeHidden,
	}

	ctx := context.Background()
	response, err := service.GetResults(ctx, outcomesCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to get outcome results: %w", err)
	}

	if len(response.OutcomeResults) == 0 {
		fmt.Println("No outcome results found")
		return nil
	}

	printVerbose("Found %d outcome results:\n\n", len(response.OutcomeResults))
	return formatOutput(response, nil)
}

func runOutcomesAlignments(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewOutcomesService(client)

	ctx := context.Background()
	alignments, err := service.GetAlignments(ctx, outcomesCourseID, 0)
	if err != nil {
		return fmt.Errorf("failed to get outcome alignments: %w", err)
	}

	if len(alignments) == 0 {
		fmt.Println("No outcome alignments found")
		return nil
	}

	printVerbose("Found %d outcome alignments:\n\n", len(alignments))
	return formatOutput(alignments, nil)
}
