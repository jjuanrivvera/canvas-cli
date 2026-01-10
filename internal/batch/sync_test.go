package batch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyncOperation(t *testing.T) {
	// Since api.Client requires real configuration, we'll test with nil
	// In real usage, proper clients would be created
	sync := NewSyncOperation(nil, nil, true)

	assert.NotNil(t, sync)
	assert.Nil(t, sync.sourceClient)
	assert.Nil(t, sync.targetClient)
	assert.True(t, sync.interactive)
}

func TestNewSyncOperation_NonInteractive(t *testing.T) {
	sync := NewSyncOperation(nil, nil, false)

	assert.NotNil(t, sync)
	assert.False(t, sync.interactive)
}

func TestConflictResolution_Constants(t *testing.T) {
	// Test that conflict resolution constants are defined
	assert.Equal(t, ConflictResolution(0), ResolutionSkip)
	assert.Equal(t, ConflictResolution(1), ResolutionOverwrite)
	assert.Equal(t, ConflictResolution(2), ResolutionMerge)
}

func TestSyncResult_Structure(t *testing.T) {
	result := &SyncResult{
		TotalItems:   10,
		SyncedItems:  7,
		SkippedItems: 2,
		FailedItems:  1,
		Errors:       []error{assert.AnError},
	}

	assert.Equal(t, 10, result.TotalItems)
	assert.Equal(t, 7, result.SyncedItems)
	assert.Equal(t, 2, result.SkippedItems)
	assert.Equal(t, 1, result.FailedItems)
	assert.Len(t, result.Errors, 1)
}

func TestSyncResult_EmptyErrors(t *testing.T) {
	result := &SyncResult{
		TotalItems:   5,
		SyncedItems:  5,
		SkippedItems: 0,
		FailedItems:  0,
		Errors:       nil,
	}

	assert.Equal(t, 5, result.TotalItems)
	assert.Equal(t, 5, result.SyncedItems)
	assert.Nil(t, result.Errors)
}
