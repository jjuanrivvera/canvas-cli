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

// analyticsCmd represents the analytics command group
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

func init() {
	rootCmd.AddCommand(analyticsCmd)
	analyticsCmd.AddCommand(newAnalyticsActivityCmd())
	analyticsCmd.AddCommand(newAnalyticsAssignmentsCmd())
	analyticsCmd.AddCommand(newAnalyticsStudentsCmd())
	analyticsCmd.AddCommand(newAnalyticsUserCmd())
	analyticsCmd.AddCommand(newAnalyticsDepartmentCmd())
}

func newAnalyticsActivityCmd() *cobra.Command {
	opts := &options.AnalyticsActivityOptions{}

	cmd := &cobra.Command{
		Use:   "activity",
		Short: "View course activity over time",
		Long: `View course participation and page view activity over time.

Examples:
  canvas analytics activity --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			if _, err := validateCourseID(client, opts.CourseID); err != nil {
				return err
			}

			return runAnalyticsActivity(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnalyticsAssignmentsCmd() *cobra.Command {
	opts := &options.AnalyticsAssignmentsOptions{}

	cmd := &cobra.Command{
		Use:   "assignments",
		Short: "View assignment analytics",
		Long: `View assignment statistics including score distribution and submission tardiness.

Examples:
  canvas analytics assignments --course-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			if _, err := validateCourseID(client, opts.CourseID); err != nil {
				return err
			}

			return runAnalyticsAssignments(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnalyticsStudentsCmd() *cobra.Command {
	opts := &options.AnalyticsStudentsOptions{}

	cmd := &cobra.Command{
		Use:   "students",
		Short: "View student summary analytics",
		Long: `View student engagement and grade summaries for a course.

Sort columns: name, name_descending, score, score_descending,
              participations, participations_descending,
              page_views, page_views_descending

Examples:
  canvas analytics students --course-id 123
  canvas analytics students --course-id 123 --sort score_descending`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			if _, err := validateCourseID(client, opts.CourseID); err != nil {
				return err
			}

			return runAnalyticsStudents(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.SortColumn, "sort", "", "Sort column (name, score, participations, page_views)")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnalyticsUserCmd() *cobra.Command {
	opts := &options.AnalyticsUserOptions{}

	cmd := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid user ID: %w", err)
			}
			opts.UserID = userID

			if err := opts.Validate(); err != nil {
				return err
			}

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			if _, err := validateCourseID(client, opts.CourseID); err != nil {
				return err
			}

			return runAnalyticsUser(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.CourseID, "course-id", 0, "Course ID (required)")
	cmd.Flags().StringVar(&opts.Type, "type", "activity", "Analytics type: activity, assignments, communication")
	cmd.MarkFlagRequired("course-id")

	return cmd
}

func newAnalyticsDepartmentCmd() *cobra.Command {
	opts := &options.AnalyticsDepartmentOptions{}

	cmd := &cobra.Command{
		Use:   "department",
		Short: "View department analytics",
		Long: `View department-level analytics for an account.

If --account-id is not specified, uses the default account ID from config.

Analytics types:
  - statistics: Overall department statistics
  - activity: Department participation over time
  - grades: Grade distribution

Examples:
  canvas analytics department                     # Uses default account
  canvas analytics department --account-id 1
  canvas analytics department --type grades --term-id 123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			accountID, err := resolveAccountID(opts.AccountID, "analytics department")
			if err != nil {
				return err
			}
			opts.AccountID = accountID

			client, err := getAPIClient()
			if err != nil {
				return err
			}

			return runAnalyticsDepartment(cmd.Context(), client, opts)
		},
	}

	cmd.Flags().Int64Var(&opts.AccountID, "account-id", 0, "Account ID (uses default if configured)")
	cmd.Flags().StringVar(&opts.Type, "type", "statistics", "Analytics type: statistics, activity, grades")
	cmd.Flags().Int64Var(&opts.TermID, "term-id", 0, "Filter by term ID")

	return cmd
}

func runAnalyticsActivity(ctx context.Context, client *api.Client, opts *options.AnalyticsActivityOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "analytics.activity", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewAnalyticsService(client)

	activity, err := service.GetCourseActivity(ctx, opts.CourseID)
	if err != nil {
		logger.LogCommandError(ctx, "analytics.activity", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get course activity: %w", err)
	}

	if len(activity) == 0 {
		fmt.Println("No activity data found")
		logger.LogCommandComplete(ctx, "analytics.activity", 0)
		return nil
	}

	printVerbose("Course activity data (%d days):\n\n", len(activity))
	logger.LogCommandComplete(ctx, "analytics.activity", len(activity))
	return formatOutput(activity, nil)
}

func runAnalyticsAssignments(ctx context.Context, client *api.Client, opts *options.AnalyticsAssignmentsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "analytics.assignments", map[string]interface{}{
		"course_id": opts.CourseID,
	})

	service := api.NewAnalyticsService(client)

	assignments, err := service.GetCourseAssignments(ctx, opts.CourseID)
	if err != nil {
		logger.LogCommandError(ctx, "analytics.assignments", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get assignment analytics: %w", err)
	}

	if len(assignments) == 0 {
		fmt.Println("No assignment data found")
		logger.LogCommandComplete(ctx, "analytics.assignments", 0)
		return nil
	}

	printVerbose("Assignment analytics (%d assignments):\n\n", len(assignments))
	logger.LogCommandComplete(ctx, "analytics.assignments", len(assignments))
	return formatOutput(assignments, nil)
}

func runAnalyticsStudents(ctx context.Context, client *api.Client, opts *options.AnalyticsStudentsOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "analytics.students", map[string]interface{}{
		"course_id":   opts.CourseID,
		"sort_column": opts.SortColumn,
	})

	service := api.NewAnalyticsService(client)

	apiOpts := &api.ListStudentSummariesOptions{
		SortColumn: opts.SortColumn,
	}

	summaries, err := service.GetStudentSummaries(ctx, opts.CourseID, apiOpts)
	if err != nil {
		logger.LogCommandError(ctx, "analytics.students", err, map[string]interface{}{
			"course_id": opts.CourseID,
		})
		return fmt.Errorf("failed to get student summaries: %w", err)
	}

	if len(summaries) == 0 {
		fmt.Println("No student data found")
		logger.LogCommandComplete(ctx, "analytics.students", 0)
		return nil
	}

	printVerbose("Student summaries (%d students):\n\n", len(summaries))
	logger.LogCommandComplete(ctx, "analytics.students", len(summaries))
	return formatOutput(summaries, nil)
}

func runAnalyticsUser(ctx context.Context, client *api.Client, opts *options.AnalyticsUserOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "analytics.user", map[string]interface{}{
		"course_id": opts.CourseID,
		"user_id":   opts.UserID,
		"type":      opts.Type,
	})

	service := api.NewAnalyticsService(client)

	switch opts.Type {
	case "activity":
		activity, err := service.GetUserActivity(ctx, opts.CourseID, opts.UserID)
		if err != nil {
			logger.LogCommandError(ctx, "analytics.user", err, map[string]interface{}{
				"course_id": opts.CourseID,
				"user_id":   opts.UserID,
				"type":      opts.Type,
			})
			return fmt.Errorf("failed to get user activity: %w", err)
		}
		if len(activity) == 0 {
			fmt.Println("No activity data found")
			logger.LogCommandComplete(ctx, "analytics.user", 0)
			return nil
		}
		logger.LogCommandComplete(ctx, "analytics.user", len(activity))
		return formatOutput(activity, nil)

	case "assignments":
		assignments, err := service.GetUserAssignments(ctx, opts.CourseID, opts.UserID)
		if err != nil {
			logger.LogCommandError(ctx, "analytics.user", err, map[string]interface{}{
				"course_id": opts.CourseID,
				"user_id":   opts.UserID,
				"type":      opts.Type,
			})
			return fmt.Errorf("failed to get user assignments: %w", err)
		}
		if len(assignments) == 0 {
			fmt.Println("No assignment data found")
			logger.LogCommandComplete(ctx, "analytics.user", 0)
			return nil
		}
		logger.LogCommandComplete(ctx, "analytics.user", len(assignments))
		return formatOutput(assignments, nil)

	case "communication":
		communication, err := service.GetUserCommunication(ctx, opts.CourseID, opts.UserID)
		if err != nil {
			logger.LogCommandError(ctx, "analytics.user", err, map[string]interface{}{
				"course_id": opts.CourseID,
				"user_id":   opts.UserID,
				"type":      opts.Type,
			})
			return fmt.Errorf("failed to get user communication: %w", err)
		}
		logger.LogCommandComplete(ctx, "analytics.user", 1)
		return formatOutput(communication, nil)

	default:
		logger.LogCommandError(ctx, "analytics.user", fmt.Errorf("invalid analytics type"), map[string]interface{}{
			"course_id": opts.CourseID,
			"user_id":   opts.UserID,
			"type":      opts.Type,
		})
		return fmt.Errorf("invalid analytics type: %s (use: activity, assignments, communication)", opts.Type)
	}
}

func runAnalyticsDepartment(ctx context.Context, client *api.Client, opts *options.AnalyticsDepartmentOptions) error {
	logger := logging.NewCommandLogger(verbose)
	logger.LogCommandStart(ctx, "analytics.department", map[string]interface{}{
		"account_id": opts.AccountID,
		"type":       opts.Type,
		"term_id":    opts.TermID,
	})

	service := api.NewAnalyticsService(client)

	apiOpts := &api.DepartmentAnalyticsOptions{}
	if opts.TermID > 0 {
		apiOpts.TermID = opts.TermID
	}

	switch opts.Type {
	case "statistics":
		stats, err := service.GetDepartmentStatistics(ctx, opts.AccountID, apiOpts)
		if err != nil {
			logger.LogCommandError(ctx, "analytics.department", err, map[string]interface{}{
				"account_id": opts.AccountID,
				"type":       opts.Type,
			})
			return fmt.Errorf("failed to get department statistics: %w", err)
		}
		logger.LogCommandComplete(ctx, "analytics.department", 1)
		return formatOutput(stats, nil)

	case "activity":
		activity, err := service.GetDepartmentActivity(ctx, opts.AccountID, apiOpts)
		if err != nil {
			logger.LogCommandError(ctx, "analytics.department", err, map[string]interface{}{
				"account_id": opts.AccountID,
				"type":       opts.Type,
			})
			return fmt.Errorf("failed to get department activity: %w", err)
		}
		if len(activity) == 0 {
			fmt.Println("No activity data found")
			logger.LogCommandComplete(ctx, "analytics.department", 0)
			return nil
		}
		logger.LogCommandComplete(ctx, "analytics.department", len(activity))
		return formatOutput(activity, nil)

	case "grades":
		grades, err := service.GetDepartmentGrades(ctx, opts.AccountID, apiOpts)
		if err != nil {
			logger.LogCommandError(ctx, "analytics.department", err, map[string]interface{}{
				"account_id": opts.AccountID,
				"type":       opts.Type,
			})
			return fmt.Errorf("failed to get department grades: %w", err)
		}
		if len(grades) == 0 {
			fmt.Println("No grade data found")
			logger.LogCommandComplete(ctx, "analytics.department", 0)
			return nil
		}
		logger.LogCommandComplete(ctx, "analytics.department", len(grades))
		return formatOutput(grades, nil)

	default:
		logger.LogCommandError(ctx, "analytics.department", fmt.Errorf("invalid analytics type"), map[string]interface{}{
			"account_id": opts.AccountID,
			"type":       opts.Type,
		})
		return fmt.Errorf("invalid analytics type: %s (use: statistics, activity, grades)", opts.Type)
	}
}
