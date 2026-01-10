package repl

import (
	"sync"
	"testing"
)

func TestNewSession(t *testing.T) {
	session := NewSession()

	if session == nil {
		t.Fatal("expected non-nil session")
	}

	if session.Variables == nil {
		t.Error("expected variables map to be initialized")
	}

	if session.CourseID != 0 {
		t.Errorf("expected CourseID to be 0, got %d", session.CourseID)
	}

	if session.UserID != 0 {
		t.Errorf("expected UserID to be 0, got %d", session.UserID)
	}

	if session.AssignmentID != 0 {
		t.Errorf("expected AssignmentID to be 0, got %d", session.AssignmentID)
	}
}

func TestSession_SetGet(t *testing.T) {
	session := NewSession()

	// Test setting and getting a variable
	session.Set("test_key", "test_value")

	value, exists := session.Get("test_key")
	if !exists {
		t.Error("expected variable to exist")
	}

	if value != "test_value" {
		t.Errorf("expected 'test_value', got '%v'", value)
	}
}

func TestSession_SetGetCourseID(t *testing.T) {
	session := NewSession()

	// Test setting course_id as special variable
	session.Set("course_id", int64(123))

	value, exists := session.Get("course_id")
	if !exists {
		t.Error("expected course_id to exist")
	}

	if value != int64(123) {
		t.Errorf("expected 123, got %v", value)
	}

	if session.CourseID != 123 {
		t.Errorf("expected CourseID to be 123, got %d", session.CourseID)
	}
}

func TestSession_SetGetUserID(t *testing.T) {
	session := NewSession()

	// Test setting user_id as special variable
	session.Set("user_id", int64(456))

	value, exists := session.Get("user_id")
	if !exists {
		t.Error("expected user_id to exist")
	}

	if value != int64(456) {
		t.Errorf("expected 456, got %v", value)
	}

	if session.UserID != 456 {
		t.Errorf("expected UserID to be 456, got %d", session.UserID)
	}
}

func TestSession_SetGetAssignmentID(t *testing.T) {
	session := NewSession()

	// Test setting assignment_id as special variable
	session.Set("assignment_id", int64(789))

	value, exists := session.Get("assignment_id")
	if !exists {
		t.Error("expected assignment_id to exist")
	}

	if value != int64(789) {
		t.Errorf("expected 789, got %v", value)
	}

	if session.AssignmentID != 789 {
		t.Errorf("expected AssignmentID to be 789, got %d", session.AssignmentID)
	}
}

func TestSession_GetNonexistent(t *testing.T) {
	session := NewSession()

	_, exists := session.Get("nonexistent")
	if exists {
		t.Error("expected nonexistent variable to not exist")
	}
}

func TestSession_Clear(t *testing.T) {
	session := NewSession()

	// Set some data
	session.Set("test_key", "test_value")
	session.Set("course_id", int64(123))
	session.Set("user_id", int64(456))
	session.Set("assignment_id", int64(789))

	// Clear session
	session.Clear()

	// Verify everything is cleared
	if session.CourseID != 0 {
		t.Errorf("expected CourseID to be 0 after clear, got %d", session.CourseID)
	}

	if session.UserID != 0 {
		t.Errorf("expected UserID to be 0 after clear, got %d", session.UserID)
	}

	if session.AssignmentID != 0 {
		t.Errorf("expected AssignmentID to be 0 after clear, got %d", session.AssignmentID)
	}

	if len(session.Variables) != 0 {
		t.Errorf("expected empty variables map after clear, got %d items", len(session.Variables))
	}

	_, exists := session.Get("test_key")
	if exists {
		t.Error("expected variable to be cleared")
	}
}

func TestSession_SetCourseID(t *testing.T) {
	session := NewSession()

	session.SetCourseID(123)

	if session.GetCourseID() != 123 {
		t.Errorf("expected CourseID 123, got %d", session.GetCourseID())
	}
}

func TestSession_SetUserID(t *testing.T) {
	session := NewSession()

	session.SetUserID(456)

	if session.GetUserID() != 456 {
		t.Errorf("expected UserID 456, got %d", session.GetUserID())
	}
}

func TestSession_SetAssignmentID(t *testing.T) {
	session := NewSession()

	session.SetAssignmentID(789)

	if session.GetAssignmentID() != 789 {
		t.Errorf("expected AssignmentID 789, got %d", session.GetAssignmentID())
	}
}

func TestSession_GetCourseID_Zero(t *testing.T) {
	session := NewSession()

	if session.GetCourseID() != 0 {
		t.Errorf("expected CourseID 0 for new session, got %d", session.GetCourseID())
	}

	// Verify Get returns false for zero course_id
	_, exists := session.Get("course_id")
	if exists {
		t.Error("expected Get to return false for zero course_id")
	}
}

func TestSession_GetUserID_Zero(t *testing.T) {
	session := NewSession()

	if session.GetUserID() != 0 {
		t.Errorf("expected UserID 0 for new session, got %d", session.GetUserID())
	}

	// Verify Get returns false for zero user_id
	_, exists := session.Get("user_id")
	if exists {
		t.Error("expected Get to return false for zero user_id")
	}
}

func TestSession_GetAssignmentID_Zero(t *testing.T) {
	session := NewSession()

	if session.GetAssignmentID() != 0 {
		t.Errorf("expected AssignmentID 0 for new session, got %d", session.GetAssignmentID())
	}

	// Verify Get returns false for zero assignment_id
	_, exists := session.Get("assignment_id")
	if exists {
		t.Error("expected Get to return false for zero assignment_id")
	}
}

func TestSession_SetMultipleVariables(t *testing.T) {
	session := NewSession()

	// Set multiple variables
	session.Set("var1", "value1")
	session.Set("var2", 123)
	session.Set("var3", true)

	// Verify all can be retrieved
	value1, exists1 := session.Get("var1")
	value2, exists2 := session.Get("var2")
	value3, exists3 := session.Get("var3")

	if !exists1 || !exists2 || !exists3 {
		t.Error("expected all variables to exist")
	}

	if value1 != "value1" {
		t.Errorf("expected 'value1', got '%v'", value1)
	}
	if value2 != 123 {
		t.Errorf("expected 123, got %v", value2)
	}
	if value3 != true {
		t.Errorf("expected true, got %v", value3)
	}
}

func TestSession_OverwriteVariable(t *testing.T) {
	session := NewSession()

	// Set a variable
	session.Set("test_key", "value1")

	// Overwrite it
	session.Set("test_key", "value2")

	value, exists := session.Get("test_key")
	if !exists {
		t.Error("expected variable to exist")
	}

	if value != "value2" {
		t.Errorf("expected 'value2', got '%v'", value)
	}
}

func TestSession_ConcurrentAccess(t *testing.T) {
	session := NewSession()

	// Test concurrent reads and writes
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// Writer goroutine
		go func(val int) {
			defer wg.Done()
			session.Set("concurrent_key", val)
		}(i)

		// Reader goroutine
		go func() {
			defer wg.Done()
			session.Get("concurrent_key")
		}()
	}

	wg.Wait()

	// Session should still be functional
	session.Set("test", "value")
	value, exists := session.Get("test")
	if !exists || value != "value" {
		t.Error("session should still be functional after concurrent access")
	}
}

func TestSession_SetInvalidTypeForSpecialVariable(t *testing.T) {
	session := NewSession()

	// Try setting course_id with wrong type (should be ignored)
	session.Set("course_id", "not an int")

	if session.CourseID != 0 {
		t.Error("expected CourseID to remain 0 when set with wrong type")
	}
}

func TestSession_ClearMultipleTimes(t *testing.T) {
	session := NewSession()

	// Set some data
	session.Set("key1", "value1")
	session.SetCourseID(123)

	// Clear multiple times
	session.Clear()
	session.Clear()

	// Should still work correctly
	if session.CourseID != 0 {
		t.Error("expected CourseID to be 0 after multiple clears")
	}

	// Setting new data after clear should work
	session.Set("key2", "value2")
	value, exists := session.Get("key2")
	if !exists || value != "value2" {
		t.Error("setting new data after clear should work")
	}
}
