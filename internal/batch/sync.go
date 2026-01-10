package batch

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/api"
	"golang.org/x/term"
)

// ConflictResolution defines how to handle conflicts during sync
type ConflictResolution int

const (
	ResolutionSkip ConflictResolution = iota
	ResolutionOverwrite
	ResolutionMerge
)

// promptTimeout is the maximum time to wait for user input in interactive mode
const promptTimeout = 60 * time.Second

// isTerminal checks if stdin is a terminal (TTY)
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// promptWithTimeout reads user input with a timeout
// Returns empty string on timeout or error
func promptWithTimeout(ctx context.Context, timeout time.Duration) (string, error) {
	if !isTerminal() {
		return "", fmt.Errorf("not running in a terminal, cannot prompt for input")
	}

	inputCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			errCh <- err
			return
		}
		inputCh <- strings.TrimSpace(line)
	}()

	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case input := <-inputCh:
		return input, nil
	case err := <-errCh:
		return "", err
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("input timeout after %v", timeout)
		}
		return "", timeoutCtx.Err()
	}
}

// SyncOperation handles cross-instance synchronization
type SyncOperation struct {
	sourceClient *api.Client
	targetClient *api.Client
	interactive  bool
}

// NewSyncOperation creates a new sync operation
func NewSyncOperation(source, target *api.Client, interactive bool) *SyncOperation {
	return &SyncOperation{
		sourceClient: source,
		targetClient: target,
		interactive:  interactive,
	}
}

// CopyAssignment copies an assignment from source to target instance
func (s *SyncOperation) CopyAssignment(ctx context.Context, sourceCourseID, targetCourseID, assignmentID int64) error {
	sourceAssignments := api.NewAssignmentsService(s.sourceClient)
	targetAssignments := api.NewAssignmentsService(s.targetClient)

	// Fetch source assignment
	assignment, err := sourceAssignments.Get(ctx, sourceCourseID, assignmentID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch source assignment: %w", err)
	}

	// Check if exists in target
	existing, err := targetAssignments.Get(ctx, targetCourseID, assignmentID, nil)
	if err == nil && existing != nil {
		// Conflict detected
		if s.interactive {
			resolution := s.promptConflict(ctx, assignment, existing)
			switch resolution {
			case ResolutionSkip:
				return nil
			case ResolutionOverwrite:
				// Continue to update
			case ResolutionMerge:
				return fmt.Errorf("merge not yet implemented")
			}
		} else {
			return fmt.Errorf("conflict: assignment %d already exists in target", assignmentID)
		}
	}

	// Create assignment in target (Canvas API doesn't support direct copy with same ID)
	// We need to create a new assignment with the same properties
	return s.createAssignmentInTarget(ctx, targetCourseID, assignment)
}

// createAssignmentInTarget creates a new assignment in the target instance
// copying properties from the source assignment
func (s *SyncOperation) createAssignmentInTarget(ctx context.Context, courseID int64, assignment *api.Assignment) error {
	targetAssignments := api.NewAssignmentsService(s.targetClient)

	// Map the source assignment properties to CreateAssignmentParams
	// Note: Some properties may not transfer directly (e.g., IDs, course-specific settings)
	params := &api.CreateAssignmentParams{
		Name:                   assignment.Name,
		Description:            assignment.Description,
		PointsPossible:         assignment.PointsPossible,
		GradingType:            assignment.GradingType,
		SubmissionTypes:        assignment.SubmissionTypes,
		AllowedExtensions:      assignment.AllowedExtensions,
		PeerReviews:            assignment.PeerReviews,
		AutomaticPeerReviews:   assignment.AutomaticPeerReviews,
		Published:              assignment.Published,
		OmitFromFinalGrade:     assignment.OmitFromFinalGrade,
		ModeratedGrading:       assignment.ModeratedGrading,
		GraderCount:            assignment.GraderCount,
		AnonymousGrading:       assignment.AnonymousGrading,
		AllowedAttempts:        assignment.AllowedAttempts,
	}

	// Convert time fields to ISO8601 strings if set
	if !assignment.DueAt.IsZero() {
		params.DueAt = assignment.DueAt.Format("2006-01-02T15:04:05Z")
	}
	if !assignment.LockAt.IsZero() {
		params.LockAt = assignment.LockAt.Format("2006-01-02T15:04:05Z")
	}
	if !assignment.UnlockAt.IsZero() {
		params.UnlockAt = assignment.UnlockAt.Format("2006-01-02T15:04:05Z")
	}

	_, err := targetAssignments.Create(ctx, courseID, params)
	if err != nil {
		return fmt.Errorf("failed to create assignment in target: %w", err)
	}

	return nil
}

// CopyCourse copies course structure from source to target
func (s *SyncOperation) CopyCourse(ctx context.Context, sourceCourseID, targetCourseID int64) error {
	sourceCourses := api.NewCoursesService(s.sourceClient)
	targetCourses := api.NewCoursesService(s.targetClient)

	// Fetch source course
	sourceCourse, err := sourceCourses.Get(ctx, sourceCourseID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch source course: %w", err)
	}

	// Fetch target course
	targetCourse, err := targetCourses.Get(ctx, targetCourseID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch target course: %w", err)
	}

	// Check for conflicts - use UpdatedAt to detect actual modifications
	if s.interactive && !sourceCourse.UpdatedAt.Equal(targetCourse.UpdatedAt) {
		fmt.Printf("\n⚠️  Course modified in both instances\n")
		fmt.Printf("Source: %s (modified: %s)\n", sourceCourse.Name, sourceCourse.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Target: %s (modified: %s)\n", targetCourse.Name, targetCourse.UpdatedAt.Format("2006-01-02 15:04:05"))

		resolution := s.promptCourseConflict(ctx)
		if resolution == ResolutionSkip {
			return nil
		}
	}

	// Sync course assignments
	sourceAssignments := api.NewAssignmentsService(s.sourceClient)
	assignments, err := sourceAssignments.List(ctx, sourceCourseID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch assignments: %w", err)
	}

	for _, assignment := range assignments {
		if err := s.CopyAssignment(ctx, sourceCourseID, targetCourseID, assignment.ID); err != nil {
			if s.interactive {
				fmt.Printf("⚠️  Failed to copy assignment %s: %v\n", assignment.Name, err)
				continue
			}
			return err
		}
	}

	return nil
}

// promptConflict prompts user for conflict resolution
func (s *SyncOperation) promptConflict(ctx context.Context, source, target *api.Assignment) ConflictResolution {
	fmt.Printf("\n⚠️  Conflict detected for assignment: %s\n", source.Name)
	fmt.Printf("Source: %s (modified: %v)\n", source.Name, source.UpdatedAt)
	fmt.Printf("Target: %s (modified: %v)\n", target.Name, target.UpdatedAt)
	fmt.Println("\nChoose action:")
	fmt.Println("  [s] Skip this assignment (default)")
	fmt.Println("  [o] Overwrite target with source")
	fmt.Println("  [m] Merge (interactive)")
	fmt.Printf("\nYour choice (timeout in %v): ", promptTimeout)

	choice, err := promptWithTimeout(ctx, promptTimeout)
	if err != nil {
		fmt.Printf("\n⚠️  %v, defaulting to skip\n", err)
		return ResolutionSkip
	}

	switch strings.ToLower(choice) {
	case "o":
		return ResolutionOverwrite
	case "m":
		return ResolutionMerge
	default:
		return ResolutionSkip
	}
}

// promptCourseConflict prompts user for course-level conflict resolution
func (s *SyncOperation) promptCourseConflict(ctx context.Context) ConflictResolution {
	fmt.Println("\nChoose action:")
	fmt.Println("  [s] Skip sync (default)")
	fmt.Println("  [o] Overwrite target with source")
	fmt.Printf("\nYour choice (timeout in %v): ", promptTimeout)

	choice, err := promptWithTimeout(ctx, promptTimeout)
	if err != nil {
		fmt.Printf("\n⚠️  %v, defaulting to skip\n", err)
		return ResolutionSkip
	}

	switch strings.ToLower(choice) {
	case "o":
		return ResolutionOverwrite
	default:
		return ResolutionSkip
	}
}

// SyncResult contains the result of a sync operation
type SyncResult struct {
	TotalItems     int
	SyncedItems    int
	SkippedItems   int
	FailedItems    int
	Errors         []error
}

// defaultConcurrency is the default number of concurrent workers for batch sync
const defaultConcurrency = 5

// SyncAssignments synchronizes all assignments from source to target course
// Uses concurrent batch processing for improved performance
func (s *SyncOperation) SyncAssignments(ctx context.Context, sourceCourseID, targetCourseID int64) (*SyncResult, error) {
	result := &SyncResult{}

	sourceAssignments := api.NewAssignmentsService(s.sourceClient)

	// Fetch all assignments from source
	assignments, err := sourceAssignments.List(ctx, sourceCourseID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source assignments: %w", err)
	}

	result.TotalItems = len(assignments)

	if len(assignments) == 0 {
		return result, nil
	}

	// Convert assignments to interface slice for batch processor
	items := make([]interface{}, len(assignments))
	for i, a := range assignments {
		items[i] = a
	}

	// Use batch processor for concurrent sync
	// In interactive mode, don't stop on first error (allow user to resolve conflicts)
	// In non-interactive mode, stop on first error
	processor := New(defaultConcurrency, !s.interactive, NewConsoleProgress(time.Second))

	summary, err := processor.Process(ctx, items, func(ctx context.Context, item interface{}) error {
		assignment := item.(*api.Assignment)
		return s.CopyAssignment(ctx, sourceCourseID, targetCourseID, assignment.ID)
	})

	if summary != nil {
		result.SyncedItems = summary.Succeeded
		result.FailedItems = summary.Failed
		result.SkippedItems = result.TotalItems - result.SyncedItems - result.FailedItems
		result.Errors = summary.Errors()
	}

	if err != nil {
		return result, err
	}

	return result, nil
}
