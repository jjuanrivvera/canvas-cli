package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jjuanrivvera/canvas-cli/internal/config"
)

// State tracks update-related state
type State struct {
	LastCheckTime    time.Time `json:"last_check_time"`
	LastUpdateTime   time.Time `json:"last_update_time,omitempty"`
	LastVersion      string    `json:"last_version,omitempty"`
	UpdatedToVersion string    `json:"updated_to_version,omitempty"`
	LastError        string    `json:"last_error,omitempty"`
	LastErrorTime    time.Time `json:"last_error_time,omitempty"`
}

// StateManager handles loading and saving update state
type StateManager struct {
	statePath string
}

// NewStateManager creates a new StateManager
func NewStateManager() (*StateManager, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}

	statePath := filepath.Join(configDir, "update-state.json")
	return &StateManager{statePath: statePath}, nil
}

// Load loads the update state from disk
func (m *StateManager) Load() (*State, error) {
	data, err := os.ReadFile(m.statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{}, nil
		}
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		// If file is corrupted, return fresh state
		return &State{}, nil
	}

	return &state, nil
}

// Save saves the update state to disk
func (m *StateManager) Save(state *State) error {
	// Ensure directory exists
	dir := filepath.Dir(m.statePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.statePath, data, 0600)
}

// ShouldCheck returns true if enough time has passed since the last check
func (m *StateManager) ShouldCheck(interval time.Duration) bool {
	state, err := m.Load()
	if err != nil {
		return true // Check if we can't read state
	}

	return time.Since(state.LastCheckTime) >= interval
}

// RecordCheck records that an update check was performed
func (m *StateManager) RecordCheck(result *UpdateResult) error {
	state, _ := m.Load()

	state.LastCheckTime = time.Now()

	if result.Updated {
		state.LastUpdateTime = time.Now()
		state.LastVersion = result.FromVersion
		state.UpdatedToVersion = result.ToVersion
		state.LastError = ""
		state.LastErrorTime = time.Time{}
	}

	if result.Error != nil {
		state.LastError = result.Error.Error()
		state.LastErrorTime = time.Now()
	}

	return m.Save(state)
}

// GetPendingNotification returns update info if there's a recent update to notify about
// It clears the notification after returning it
func (m *StateManager) GetPendingNotification() (fromVersion, toVersion string, hasNotification bool) {
	state, err := m.Load()
	if err != nil {
		return "", "", false
	}

	// Check if there's a recent update (within last 5 minutes) to notify about
	if state.UpdatedToVersion != "" && time.Since(state.LastUpdateTime) < 5*time.Minute {
		fromVersion = state.LastVersion
		toVersion = state.UpdatedToVersion

		// Clear the notification (ignore save error - non-critical)
		state.UpdatedToVersion = ""
		_ = m.Save(state)

		return fromVersion, toVersion, true
	}

	return "", "", false
}

// GetLastError returns the last error if it occurred recently
func (m *StateManager) GetLastError() (string, bool) {
	state, err := m.Load()
	if err != nil {
		return "", false
	}

	// Only report errors from the last 5 minutes
	if state.LastError != "" && time.Since(state.LastErrorTime) < 5*time.Minute {
		errorMsg := state.LastError

		// Clear the error after reporting (ignore save error - non-critical)
		state.LastError = ""
		state.LastErrorTime = time.Time{}
		_ = m.Save(state)

		return errorMsg, true
	}

	return "", false
}
