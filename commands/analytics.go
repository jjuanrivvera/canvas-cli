package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
)

var (
	analyticsCourseID   int64
	analyticsAccountID  int64
	analyticsType       string
	analyticsSortColumn string
	analyticsTermID     int64
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "View Canvas analytics",
	Long: `View analytics data for courses, students, and departments.

Analytics provide insights into participation, assignment performance,
and engagement metrics.

Examples:
  canvas analytics activity --course-id 123
  canvas analytics assignments --course-id 123
  canvas analytics students --course-id 123
  canvas analytics user 456 --course-id 123`,
}

var analyticsActivityCmd = &cobra.Command{
	Use:   "activity",
	Short: "View course activity over time",
	Long: `View course participation and page view activity over time.

Examples:
  canvas analytics activity --course-id 123`,
	RunE: runAnalyticsActivity,
}

var analyticsAssignmentsCmd = &cobra.Command{
	Use:   "assignments",
	Short: "View assignment analytics",
	Long: `View assignment statistics including score distribution and submission tardiness.

Examples:
  canvas analytics assignments --course-id 123`,
	RunE: runAnalyticsAssignments,
}

var analyticsStudentsCmd = &cobra.Command{
	Use:   "students",
	Short: "View student summary analytics",
	Long: `View student engagement and grade summaries for a course.

Sort columns: name, name_descending, score, score_descending,
              participations, participations_descending,
              page_views, page_views_descending

Examples:
  canvas analytics students --course-id 123
  canvas analytics students --course-id 123 --sort score_descending`,
	RunE: runAnalyticsStudents,
}

var analyticsUserCmd = &cobra.Command{
	Use:   "user <user-id>",
	Short: "View user analytics",
	Long: `View analytics for a specific user in a course.

Analytics types:
  - activity: Page views and participation over time
  - assignments: Assignment scores and submission data
  - communication: Messaging statistics

Examples:
  canvas analytics user 456 --course-id 123
  canvas analytics user 456 --course-id 123 --type assignments`,
	Args: ExactArgsWithUsage(1, "user-id"),
	RunE: runAnalyticsUser,
}

var analyticsDepartmentCmd = &cobra.Command{
	Use:   "department",
	Short: "View department analytics",
	Long: `View department-level analytics for an account.

Analytics types:
  - statistics: Overall department statistics
  - activity: Department participation over time
  - grades: Grade distribution

Examples:
  canvas analytics department --account-id 1
  canvas analytics department --account-id 1 --type statistics
  canvas analytics department --account-id 1 --type grades --term-id 123`,
	RunE: runAnalyticsDepartment,
}

func init() {
	rootCmd.AddCommand(analyticsCmd)
	analyticsCmd.AddCommand(analyticsActivityCmd)
	analyticsCmd.AddCommand(analyticsAssignmentsCmd)
	analyticsCmd.AddCommand(analyticsStudentsCmd)
	analyticsCmd.AddCommand(analyticsUserCmd)
	analyticsCmd.AddCommand(analyticsDepartmentCmd)

	// Activity flags
	analyticsActivityCmd.Flags().Int64Var(&analyticsCourseID, "course-id", 0, "Course ID (required)")
	analyticsActivityCmd.MarkFlagRequired("course-id")

	// Assignments flags
	analyticsAssignmentsCmd.Flags().Int64Var(&analyticsCourseID, "course-id", 0, "Course ID (required)")
	analyticsAssignmentsCmd.MarkFlagRequired("course-id")

	// Students flags
	analyticsStudentsCmd.Flags().Int64Var(&analyticsCourseID, "course-id", 0, "Course ID (required)")
	analyticsStudentsCmd.Flags().StringVar(&analyticsSortColumn, "sort", "", "Sort column (name, score, participations, page_views)")
	analyticsStudentsCmd.MarkFlagRequired("course-id")

	// User flags
	analyticsUserCmd.Flags().Int64Var(&analyticsCourseID, "course-id", 0, "Course ID (required)")
	analyticsUserCmd.Flags().StringVar(&analyticsType, "type", "activity", "Analytics type: activity, assignments, communication")
	analyticsUserCmd.MarkFlagRequired("course-id")

	// Department flags
	analyticsDepartmentCmd.Flags().Int64Var(&analyticsAccountID, "account-id", 0, "Account ID (required)")
	analyticsDepartmentCmd.Flags().StringVar(&analyticsType, "type", "statistics", "Analytics type: statistics, activity, grades")
	analyticsDepartmentCmd.Flags().Int64Var(&analyticsTermID, "term-id", 0, "Filter by term ID")
	analyticsDepartmentCmd.MarkFlagRequired("account-id")
}

func runAnalyticsActivity(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, analyticsCourseID); err != nil {
		return err
	}

	service := api.NewAnalyticsService(client)

	ctx := context.Background()
	activity, err := service.GetCourseActivity(ctx, analyticsCourseID)
	if err != nil {
		return fmt.Errorf("failed to get course activity: %w", err)
	}

	if len(activity) == 0 {
		fmt.Println("No activity data found")
		return nil
	}

	printVerbose("Course activity data (%d days):\n\n", len(activity))
	return formatOutput(activity, nil)
}

func runAnalyticsAssignments(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, analyticsCourseID); err != nil {
		return err
	}

	service := api.NewAnalyticsService(client)

	ctx := context.Background()
	assignments, err := service.GetCourseAssignments(ctx, analyticsCourseID)
	if err != nil {
		return fmt.Errorf("failed to get assignment analytics: %w", err)
	}

	if len(assignments) == 0 {
		fmt.Println("No assignment data found")
		return nil
	}

	printVerbose("Assignment analytics (%d assignments):\n\n", len(assignments))
	return formatOutput(assignments, nil)
}

func runAnalyticsStudents(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, analyticsCourseID); err != nil {
		return err
	}

	service := api.NewAnalyticsService(client)

	opts := &api.ListStudentSummariesOptions{
		SortColumn: analyticsSortColumn,
	}

	ctx := context.Background()
	summaries, err := service.GetStudentSummaries(ctx, analyticsCourseID, opts)
	if err != nil {
		return fmt.Errorf("failed to get student summaries: %w", err)
	}

	if len(summaries) == 0 {
		fmt.Println("No student data found")
		return nil
	}

	printVerbose("Student summaries (%d students):\n\n", len(summaries))
	return formatOutput(summaries, nil)
}

func runAnalyticsUser(cmd *cobra.Command, args []string) error {
	userID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	client, err := getAPIClient()
	if err != nil {
		return err
	}

	if _, err := validateCourseID(client, analyticsCourseID); err != nil {
		return err
	}

	service := api.NewAnalyticsService(client)
	ctx := context.Background()

	switch analyticsType {
	case "activity":
		activity, err := service.GetUserActivity(ctx, analyticsCourseID, userID)
		if err != nil {
			return fmt.Errorf("failed to get user activity: %w", err)
		}
		if len(activity) == 0 {
			fmt.Println("No activity data found")
			return nil
		}
		return formatOutput(activity, nil)

	case "assignments":
		assignments, err := service.GetUserAssignments(ctx, analyticsCourseID, userID)
		if err != nil {
			return fmt.Errorf("failed to get user assignments: %w", err)
		}
		if len(assignments) == 0 {
			fmt.Println("No assignment data found")
			return nil
		}
		return formatOutput(assignments, nil)

	case "communication":
		communication, err := service.GetUserCommunication(ctx, analyticsCourseID, userID)
		if err != nil {
			return fmt.Errorf("failed to get user communication: %w", err)
		}
		return formatOutput(communication, nil)

	default:
		return fmt.Errorf("invalid analytics type: %s (use: activity, assignments, communication)", analyticsType)
	}
}

func runAnalyticsDepartment(cmd *cobra.Command, args []string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	service := api.NewAnalyticsService(client)
	ctx := context.Background()

	opts := &api.DepartmentAnalyticsOptions{}
	if analyticsTermID > 0 {
		opts.TermID = analyticsTermID
	}

	switch analyticsType {
	case "statistics":
		stats, err := service.GetDepartmentStatistics(ctx, analyticsAccountID, opts)
		if err != nil {
			return fmt.Errorf("failed to get department statistics: %w", err)
		}
		return formatOutput(stats, nil)

	case "activity":
		activity, err := service.GetDepartmentActivity(ctx, analyticsAccountID, opts)
		if err != nil {
			return fmt.Errorf("failed to get department activity: %w", err)
		}
		if len(activity) == 0 {
			fmt.Println("No activity data found")
			return nil
		}
		return formatOutput(activity, nil)

	case "grades":
		grades, err := service.GetDepartmentGrades(ctx, analyticsAccountID, opts)
		if err != nil {
			return fmt.Errorf("failed to get department grades: %w", err)
		}
		if len(grades) == 0 {
			fmt.Println("No grade data found")
			return nil
		}
		return formatOutput(grades, nil)

	default:
		return fmt.Errorf("invalid analytics type: %s (use: statistics, activity, grades)", analyticsType)
	}
}
