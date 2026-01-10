package repl

import "sync"

// Session represents a REPL session with persistent state
type Session struct {
	mu           sync.RWMutex
	CourseID     int64
	UserID       int64
	AssignmentID int64
	Variables    map[string]interface{}
}

// NewSession creates a new REPL session
func NewSession() *Session {
	return &Session{
		Variables: make(map[string]interface{}),
	}
}

// Set sets a session variable
func (s *Session) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Handle special session variables
	switch key {
	case "course_id":
		if id, ok := value.(int64); ok {
			s.CourseID = id
		}
	case "user_id":
		if id, ok := value.(int64); ok {
			s.UserID = id
		}
	case "assignment_id":
		if id, ok := value.(int64); ok {
			s.AssignmentID = id
		}
	default:
		s.Variables[key] = value
	}
}

// Get retrieves a session variable
func (s *Session) Get(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Handle special session variables
	switch key {
	case "course_id":
		return s.CourseID, s.CourseID > 0
	case "user_id":
		return s.UserID, s.UserID > 0
	case "assignment_id":
		return s.AssignmentID, s.AssignmentID > 0
	default:
		value, exists := s.Variables[key]
		return value, exists
	}
}

// Clear clears all session variables
func (s *Session) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CourseID = 0
	s.UserID = 0
	s.AssignmentID = 0
	s.Variables = make(map[string]interface{})
}

// SetCourseID sets the current course ID
func (s *Session) SetCourseID(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CourseID = id
}

// SetUserID sets the current user ID
func (s *Session) SetUserID(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UserID = id
}

// SetAssignmentID sets the current assignment ID
func (s *Session) SetAssignmentID(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AssignmentID = id
}

// GetCourseID returns the current course ID
func (s *Session) GetCourseID() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CourseID
}

// GetUserID returns the current user ID
func (s *Session) GetUserID() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.UserID
}

// GetAssignmentID returns the current assignment ID
func (s *Session) GetAssignmentID() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.AssignmentID
}
