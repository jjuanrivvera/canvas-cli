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
	rubricsCourseID  int64
	rubricsAccountID int64
	rubricsInclude   []string

	// Create/Update flags
	rubricsTitle                     string
	rubricsPointsPossible            float64
	rubricsFreeFormCriterionComments bool
	rubricsHideScoreTotal            bool

	// Associate flags
	rubricsAssignmentID  int64
	rubricsUseForGrading bool
	rubricsHidePoints    bool

	// Delete flags
	rubricsForce bool
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

// rubricsListCmd represents the rubrics list command
var rubricsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List rubrics",
	Long: `List all rubrics in a course or account.

Examples:
  canvas rubrics list --course-id 123
  canvas rubrics list --account-id 1
  canvas rubrics list --course-id 123 --include assessments,associations`,
	RunE: runRubricsList,
}

// rubricsGetCmd represents the rubrics get command
var rubricsGetCmd = &cobra.Command{
	Use:   "get <rubric-id>",
	Short: "Get rubric details",
	Long: `Get details of a specific rubric.

Examples:
  canvas rubrics get 456 --course-id 123
  canvas rubrics get 456 --account-id 1 --include assessments`,
	Args: ExactArgsWithUsage(1, "rubric-id"),
	RunE: runRubricsGet,
}

// rubricsCreateCmd represents the rubrics create command
var rubricsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new rubric",
	Long: `Create a new rubric in a course.

Examples:
  canvas rubrics create --course-id 123 --title "Essay Rubric" --points 100
  canvas rubrics create --course-id 123 --title "Discussion Rubric" --free-form`,
	RunE: runRubricsCreate,
}

// rubricsUpdateCmd represents the rubrics update command
var rubricsUpdateCmd = &cobra.Command{
	Use:   "update <rubric-id>",
	Short: "Update a rubric",
	Long: `Update an existing rubric.

Examples:
  canvas rubrics update 456 --course-id 123 --title "Updated Title"
  canvas rubrics update 456 --course-id 123 --points 150`,
	Args: ExactArgsWithUsage(1, "rubric-id"),
	RunE: runRubricsUpdate,
}

// rubricsDeleteCmd represents the rubrics delete command
var rubricsDeleteCmd = &cobra.Command{
	Use:   "delete <rubric-id>",
	Short: "Delete a rubric",
	Long: `Delete a rubric.

Examples:
  canvas rubrics delete 456 --course-id 123
  canvas rubrics delete 456 --course-id 123 --force`,
	Args: ExactArgsWithUsage(1, "rubric-id"),
	RunE: runRubricsDelete,
}

// rubricsAssociateCmd represents the rubrics associate command
var rubricsAssociateCmd = &cobra.Command{
	Use:   "associate <rubric-id>",
	Short: "Associate rubric with assignment",
	Long: `Associate a rubric with an assignment.

Examples:
  canvas rubrics associate 456 --course-id 123 --assignment-id 789
  canvas rubrics associate 456 --course-id 123 --assignment-id 789 --use-for-grading`,
	Args: ExactArgsWithUsage(1, "rubric-id"),
	RunE: runRubricsAssociate,
}

func init() {
	rootCmd.AddCommand(rubricsCmd)
	rubricsCmd.AddCommand(rubricsListCmd)
	rubricsCmd.AddCommand(rubricsGetCmd)
	rubricsCmd.AddCommand(rubricsCreateCmd)
	rubricsCmd.AddCommand(rubricsUpdateCmd)
	rubricsCmd.AddCommand(rubricsDeleteCmd)
	rubricsCmd.AddCommand(rubricsAssociateCmd)

	// List flags
	rubricsListCmd.Flags().Int64Var(&rubricsCourseID, "course-id", 0, "Course ID")
	rubricsListCmd.Flags().Int64Var(&rubricsAccountID, "account-id", 0, "Account ID")
	rubricsListCmd.Flags().StringSliceVar(&rubricsInclude, "include", []string{}, "Include additional data (assessments, associations, assignment_associations)")

	// Get flags
	rubricsGetCmd.Flags().Int64Var(&rubricsCourseID, "course-id", 0, "Course ID")
	rubricsGetCmd.Flags().Int64Var(&rubricsAccountID, "account-id", 0, "Account ID")
	rubricsGetCmd.Flags().StringSliceVar(&rubricsInclude, "include", []string{}, "Include additional data")

	// Create flags
	rubricsCreateCmd.Flags().Int64Var(&rubricsCourseID, "course-id", 0, "Course ID (required)")
	rubricsCreateCmd.MarkFlagRequired("course-id")
	rubricsCreateCmd.Flags().StringVar(&rubricsTitle, "title", "", "Rubric title (required)")
	rubricsCreateCmd.MarkFlagRequired("title")
	rubricsCreateCmd.Flags().Float64Var(&rubricsPointsPossible, "points", 0, "Total points possible")
	rubricsCreateCmd.Flags().BoolVar(&rubricsFreeFormCriterionComments, "free-form", false, "Allow free-form criterion comments")
	rubricsCreateCmd.Flags().BoolVar(&rubricsHideScoreTotal, "hide-score-total", false, "Hide score total")

	// Update flags
	rubricsUpdateCmd.Flags().Int64Var(&rubricsCourseID, "course-id", 0, "Course ID (required)")
	rubricsUpdateCmd.MarkFlagRequired("course-id")
	rubricsUpdateCmd.Flags().StringVar(&rubricsTitle, "title", "", "Rubric title")
	rubricsUpdateCmd.Flags().Float64Var(&rubricsPointsPossible, "points", 0, "Total points possible")
	rubricsUpdateCmd.Flags().BoolVar(&rubricsFreeFormCriterionComments, "free-form", false, "Allow free-form criterion comments")
	rubricsUpdateCmd.Flags().BoolVar(&rubricsHideScoreTotal, "hide-score-total", false, "Hide score total")

	// Delete flags
	rubricsDeleteCmd.Flags().Int64Var(&rubricsCourseID, "course-id", 0, "Course ID (required)")
	rubricsDeleteCmd.MarkFlagRequired("course-id")
	rubricsDeleteCmd.Flags().BoolVar(&rubricsForce, "force", false, "Skip confirmation prompt")

	// Associate flags
	rubricsAssociateCmd.Flags().Int64Var(&rubricsCourseID, "course-id", 0, "Course ID (required)")
	rubricsAssociateCmd.MarkFlagRequired("course-id")
	rubricsAssociateCmd.Flags().Int64Var(&rubricsAssignmentID, "assignment-id", 0, "Assignment ID (required)")
	rubricsAssociateCmd.MarkFlagRequired("assignment-id")
	rubricsAssociateCmd.Flags().BoolVar(&rubricsUseForGrading, "use-for-grading", false, "Use rubric for grading")
	rubricsAssociateCmd.Flags().BoolVar(&rubricsHideScoreTotal, "hide-score-total", false, "Hide score total")
	rubricsAssociateCmd.Flags().BoolVar(&rubricsHidePoints, "hide-points", false, "Hide points")
}

func runRubricsList(cmd *cobra.Command, args []string) error {
	if rubricsCourseID == 0 && rubricsAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRubricsService(client)

	opts := &api.ListRubricsOptions{
		Include: rubricsInclude,
	}

	ctx := context.Background()
	var rubrics []api.Rubric

	if rubricsCourseID > 0 {
		rubrics, err = service.ListCourse(ctx, rubricsCourseID, opts)
	} else {
		rubrics, err = service.ListAccount(ctx, rubricsAccountID, opts)
	}

	if err != nil {
		return fmt.Errorf("failed to list rubrics: %w", err)
	}

	if len(rubrics) == 0 {
		fmt.Println("No rubrics found")
		return nil
	}

	printVerbose("Found %d rubrics:\n\n", len(rubrics))
	return formatOutput(rubrics, nil)
}

func runRubricsGet(cmd *cobra.Command, args []string) error {
	rubricID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid rubric ID: %w", err)
	}

	if rubricsCourseID == 0 && rubricsAccountID == 0 {
		return fmt.Errorf("must specify either --course-id or --account-id")
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRubricsService(client)

	ctx := context.Background()
	var rubric *api.Rubric

	if rubricsCourseID > 0 {
		rubric, err = service.GetCourse(ctx, rubricsCourseID, rubricID, rubricsInclude)
	} else {
		rubric, err = service.GetAccount(ctx, rubricsAccountID, rubricID, rubricsInclude)
	}

	if err != nil {
		return fmt.Errorf("failed to get rubric: %w", err)
	}

	return formatOutput(rubric, nil)
}

func runRubricsCreate(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRubricsService(client)

	params := &api.CreateRubricParams{
		Title:                     rubricsTitle,
		PointsPossible:            rubricsPointsPossible,
		FreeFormCriterionComments: rubricsFreeFormCriterionComments,
		HideScoreTotal:            rubricsHideScoreTotal,
	}

	ctx := context.Background()
	rubric, err := service.Create(ctx, rubricsCourseID, params)
	if err != nil {
		return fmt.Errorf("failed to create rubric: %w", err)
	}

	fmt.Printf("Rubric created successfully (ID: %d)\n", rubric.ID)
	return formatOutput(rubric, nil)
}

func runRubricsUpdate(cmd *cobra.Command, args []string) error {
	rubricID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid rubric ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRubricsService(client)

	params := &api.UpdateRubricParams{}

	if cmd.Flags().Changed("title") {
		params.Title = &rubricsTitle
	}
	if cmd.Flags().Changed("points") {
		params.PointsPossible = &rubricsPointsPossible
	}
	if cmd.Flags().Changed("free-form") {
		params.FreeFormCriterionComments = &rubricsFreeFormCriterionComments
	}
	if cmd.Flags().Changed("hide-score-total") {
		params.HideScoreTotal = &rubricsHideScoreTotal
	}

	ctx := context.Background()
	rubric, err := service.Update(ctx, rubricsCourseID, rubricID, params)
	if err != nil {
		return fmt.Errorf("failed to update rubric: %w", err)
	}

	fmt.Printf("Rubric updated successfully (ID: %d)\n", rubric.ID)
	return formatOutput(rubric, nil)
}

func runRubricsDelete(cmd *cobra.Command, args []string) error {
	rubricID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid rubric ID: %w", err)
	}

	if !rubricsForce {
		fmt.Printf("WARNING: This will delete rubric %d.\n", rubricID)
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRubricsService(client)

	ctx := context.Background()
	rubric, err := service.Delete(ctx, rubricsCourseID, rubricID)
	if err != nil {
		return fmt.Errorf("failed to delete rubric: %w", err)
	}

	fmt.Printf("Rubric %d deleted\n", rubric.ID)
	return nil
}

func runRubricsAssociate(cmd *cobra.Command, args []string) error {
	rubricID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid rubric ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewRubricsService(client)

	params := &api.AssociateParams{
		AssociationType: "Assignment",
		AssociationID:   rubricsAssignmentID,
		UseForGrading:   rubricsUseForGrading,
		HideScoreTotal:  rubricsHideScoreTotal,
		HidePoints:      rubricsHidePoints,
		Purpose:         "grading",
	}

	ctx := context.Background()
	association, err := service.Associate(ctx, rubricsCourseID, rubricID, params)
	if err != nil {
		return fmt.Errorf("failed to associate rubric: %w", err)
	}

	fmt.Printf("Rubric %d associated with assignment %d (Association ID: %d)\n", rubricID, rubricsAssignmentID, association.ID)
	return formatOutput(association, nil)
}
