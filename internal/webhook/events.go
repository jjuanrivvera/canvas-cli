package webhook

// Canvas webhook event types
const (
	// Assignment events
	EventAssignmentCreated = "assignment_created"
	EventAssignmentUpdated = "assignment_updated"
	EventAssignmentDeleted = "assignment_deleted"

	// Submission events
	EventSubmissionCreated = "submission_created"
	EventSubmissionUpdated = "submission_updated"
	EventGradeChange       = "grade_change"

	// Enrollment events
	EventEnrollmentCreated = "enrollment_created"
	EventEnrollmentUpdated = "enrollment_updated"
	EventEnrollmentDeleted = "enrollment_deleted"

	// User events
	EventUserCreated = "user_created"
	EventUserUpdated = "user_updated"

	// Course events
	EventCourseCreated   = "course_created"
	EventCourseUpdated   = "course_updated"
	EventCourseCompleted = "course_completed"

	// Discussion events
	EventDiscussionTopicCreated = "discussion_topic_created"
	EventDiscussionEntryCreated = "discussion_entry_created"

	// Quiz events
	EventQuizSubmitted = "quiz_submitted"

	// Conversation events
	EventConversationCreated        = "conversation_created"
	EventConversationMessageCreated = "conversation_message_created"
)

// EventTypeNames maps event types to human-readable names
var EventTypeNames = map[string]string{
	EventAssignmentCreated:          "Assignment Created",
	EventAssignmentUpdated:          "Assignment Updated",
	EventAssignmentDeleted:          "Assignment Deleted",
	EventSubmissionCreated:          "Submission Created",
	EventSubmissionUpdated:          "Submission Updated",
	EventGradeChange:                "Grade Change",
	EventEnrollmentCreated:          "Enrollment Created",
	EventEnrollmentUpdated:          "Enrollment Updated",
	EventEnrollmentDeleted:          "Enrollment Deleted",
	EventUserCreated:                "User Created",
	EventUserUpdated:                "User Updated",
	EventCourseCreated:              "Course Created",
	EventCourseUpdated:              "Course Updated",
	EventCourseCompleted:            "Course Completed",
	EventDiscussionTopicCreated:     "Discussion Topic Created",
	EventDiscussionEntryCreated:     "Discussion Entry Created",
	EventQuizSubmitted:              "Quiz Submitted",
	EventConversationCreated:        "Conversation Created",
	EventConversationMessageCreated: "Conversation Message Created",
}

// GetEventName returns the human-readable name for an event type
func GetEventName(eventType string) string {
	if name, exists := EventTypeNames[eventType]; exists {
		return name
	}
	return eventType
}

// IsValidEventType checks if an event type is valid
func IsValidEventType(eventType string) bool {
	_, exists := EventTypeNames[eventType]
	return exists
}

// AllEventTypes returns a list of all supported event types
func AllEventTypes() []string {
	types := make([]string, 0, len(EventTypeNames))
	for t := range EventTypeNames {
		types = append(types, t)
	}
	return types
}
