package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// newStateManagerWithPath builds a StateManager pointing at an arbitrary path,
// bypassing the config.GetConfigDir dependency so tests stay self-contained.
func newStateManagerWithPath(path string) *StateManager {
	return &StateManager{statePath: path}
}

// --- Load ---

func TestStateManager_Load_MissingFile(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "missing.json"))

	state, err := m.Load()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	// Zero-value struct returned when file doesn't exist
	if !state.LastCheckTime.IsZero() {
		t.Errorf("expected zero LastCheckTime, got %v", state.LastCheckTime)
	}
}

func TestStateManager_Load_CorruptedFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")
	if err := os.WriteFile(p, []byte("not valid json {{{{"), 0600); err != nil {
		t.Fatal(err)
	}

	m := newStateManagerWithPath(p)
	state, err := m.Load()
	if err != nil {
		t.Fatalf("expected nil error for corrupted JSON (fresh state returned), got %v", err)
	}
	if state == nil {
		t.Fatal("expected non-nil state for corrupted JSON")
	}
}

func TestStateManager_Load_ValidFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "state.json")

	now := time.Now().UTC().Truncate(time.Second)
	s := &State{
		LastCheckTime:    now,
		LastVersion:      "1.0.0",
		UpdatedToVersion: "1.1.0",
	}
	data, _ := json.MarshalIndent(s, "", "  ")
	if err := os.WriteFile(p, data, 0600); err != nil {
		t.Fatal(err)
	}

	m := newStateManagerWithPath(p)
	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.LastVersion != "1.0.0" {
		t.Errorf("expected LastVersion 1.0.0, got %s", loaded.LastVersion)
	}
	if loaded.UpdatedToVersion != "1.1.0" {
		t.Errorf("expected UpdatedToVersion 1.1.0, got %s", loaded.UpdatedToVersion)
	}
}

// --- Save ---

func TestStateManager_Save_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "nested", "dir", "state.json")
	m := newStateManagerWithPath(p)

	s := &State{LastVersion: "2.0.0"}
	if err := m.Save(s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists and is readable
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("could not read saved file: %v", err)
	}

	var loaded State
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("saved JSON is invalid: %v", err)
	}
	if loaded.LastVersion != "2.0.0" {
		t.Errorf("expected LastVersion 2.0.0, got %s", loaded.LastVersion)
	}
}

func TestStateManager_Save_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	now := time.Now().UTC().Truncate(time.Second)
	original := &State{
		LastCheckTime:    now,
		LastUpdateTime:   now,
		LastVersion:      "1.0.0",
		UpdatedToVersion: "1.1.0",
		LastError:        "something went wrong",
		LastErrorTime:    now,
	}

	if err := m.Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := m.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.LastError != "something went wrong" {
		t.Errorf("unexpected LastError: %q", loaded.LastError)
	}
	if loaded.UpdatedToVersion != "1.1.0" {
		t.Errorf("unexpected UpdatedToVersion: %q", loaded.UpdatedToVersion)
	}
}

// --- ShouldCheck ---

func TestStateManager_ShouldCheck_NoStateFile(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "nonexistent.json"))

	// No file → always should check
	if !m.ShouldCheck(time.Hour) {
		t.Error("expected ShouldCheck to return true when state file missing")
	}
}

func TestStateManager_ShouldCheck_RecentCheck(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	// Record a check just now
	s := &State{LastCheckTime: time.Now()}
	if err := m.Save(s); err != nil {
		t.Fatal(err)
	}

	// Should NOT check again within the interval
	if m.ShouldCheck(time.Hour) {
		t.Error("expected ShouldCheck to return false when check was recent")
	}
}

func TestStateManager_ShouldCheck_StaleCheck(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	// Record a check 2 hours ago
	s := &State{LastCheckTime: time.Now().Add(-2 * time.Hour)}
	if err := m.Save(s); err != nil {
		t.Fatal(err)
	}

	// Should check because enough time has passed
	if !m.ShouldCheck(time.Hour) {
		t.Error("expected ShouldCheck to return true when interval has passed")
	}
}

// --- RecordCheck ---

func TestStateManager_RecordCheck_SuccessfulUpdate(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	result := &UpdateResult{
		Updated:     true,
		FromVersion: "1.0.0",
		ToVersion:   "1.1.0",
	}

	if err := m.RecordCheck(result); err != nil {
		t.Fatalf("RecordCheck failed: %v", err)
	}

	state, err := m.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if state.LastVersion != "1.0.0" {
		t.Errorf("expected LastVersion 1.0.0, got %s", state.LastVersion)
	}
	if state.UpdatedToVersion != "1.1.0" {
		t.Errorf("expected UpdatedToVersion 1.1.0, got %s", state.UpdatedToVersion)
	}
	// Successful update should clear errors
	if state.LastError != "" {
		t.Errorf("expected empty LastError after successful update, got %q", state.LastError)
	}
	if state.LastCheckTime.IsZero() {
		t.Error("expected non-zero LastCheckTime")
	}
}

func TestStateManager_RecordCheck_WithError(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	result := &UpdateResult{
		Error: errorf("network timeout"),
	}

	if err := m.RecordCheck(result); err != nil {
		t.Fatalf("RecordCheck failed: %v", err)
	}

	state, err := m.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if state.LastError != "network timeout" {
		t.Errorf("expected LastError 'network timeout', got %q", state.LastError)
	}
	if state.LastErrorTime.IsZero() {
		t.Error("expected non-zero LastErrorTime")
	}
}

func TestStateManager_RecordCheck_NoUpdate(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	result := &UpdateResult{Updated: false}
	if err := m.RecordCheck(result); err != nil {
		t.Fatalf("RecordCheck failed: %v", err)
	}

	state, err := m.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	// LastUpdateTime should remain zero when not updated
	if !state.LastUpdateTime.IsZero() {
		t.Error("expected zero LastUpdateTime when no update occurred")
	}
}

// --- GetPendingNotification ---

func TestStateManager_GetPendingNotification_NoUpdate(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	_, _, hasNotification := m.GetPendingNotification()
	if hasNotification {
		t.Error("expected no notification when no state file exists")
	}
}

func TestStateManager_GetPendingNotification_RecentUpdate(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	s := &State{
		LastUpdateTime:   time.Now(),
		LastVersion:      "1.0.0",
		UpdatedToVersion: "1.1.0",
	}
	if err := m.Save(s); err != nil {
		t.Fatal(err)
	}

	fromVersion, toVersion, hasNotification := m.GetPendingNotification()
	if !hasNotification {
		t.Fatal("expected hasNotification=true for recent update")
	}
	if fromVersion != "1.0.0" {
		t.Errorf("expected fromVersion 1.0.0, got %s", fromVersion)
	}
	if toVersion != "1.1.0" {
		t.Errorf("expected toVersion 1.1.0, got %s", toVersion)
	}

	// Notification should be cleared after retrieval
	_, _, hasNotification2 := m.GetPendingNotification()
	if hasNotification2 {
		t.Error("expected notification to be cleared after first retrieval")
	}
}

func TestStateManager_GetPendingNotification_OldUpdate(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	// Update was 10 minutes ago — outside the 5-minute window
	s := &State{
		LastUpdateTime:   time.Now().Add(-10 * time.Minute),
		LastVersion:      "1.0.0",
		UpdatedToVersion: "1.1.0",
	}
	if err := m.Save(s); err != nil {
		t.Fatal(err)
	}

	_, _, hasNotification := m.GetPendingNotification()
	if hasNotification {
		t.Error("expected no notification for update older than 5 minutes")
	}
}

// --- GetLastError ---

func TestStateManager_GetLastError_NoError(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	_, hasError := m.GetLastError()
	if hasError {
		t.Error("expected no error when state file is missing")
	}
}

func TestStateManager_GetLastError_RecentError(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	s := &State{
		LastError:     "failed to reach GitHub",
		LastErrorTime: time.Now(),
	}
	if err := m.Save(s); err != nil {
		t.Fatal(err)
	}

	errMsg, hasError := m.GetLastError()
	if !hasError {
		t.Fatal("expected hasError=true for recent error")
	}
	if errMsg != "failed to reach GitHub" {
		t.Errorf("expected error message 'failed to reach GitHub', got %q", errMsg)
	}

	// Error should be cleared after retrieval
	_, hasError2 := m.GetLastError()
	if hasError2 {
		t.Error("expected error to be cleared after first retrieval")
	}
}

func TestStateManager_GetLastError_OldError(t *testing.T) {
	dir := t.TempDir()
	m := newStateManagerWithPath(filepath.Join(dir, "state.json"))

	// Error was 10 minutes ago — outside the 5-minute window
	s := &State{
		LastError:     "old error",
		LastErrorTime: time.Now().Add(-10 * time.Minute),
	}
	if err := m.Save(s); err != nil {
		t.Fatal(err)
	}

	_, hasError := m.GetLastError()
	if hasError {
		t.Error("expected no error for error older than 5 minutes")
	}
}

// --- NewStateManager (smoke test, uses real home dir but cleans nothing) ---

func TestNewStateManager(t *testing.T) {
	m, err := NewStateManager()
	if err != nil {
		t.Fatalf("NewStateManager failed: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil StateManager")
	}
	if m.statePath == "" {
		t.Error("expected non-empty statePath")
	}
}

// errorf creates a simple error value for testing without importing errors package.
func errorf(msg string) error {
	return &simpleError{msg: msg}
}

type simpleError struct{ msg string }

func (e *simpleError) Error() string { return e.msg }
