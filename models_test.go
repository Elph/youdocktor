package main

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

func TestNewTimeSheetEntry_ShouldCreateAsExpected(t *testing.T) {

	// arrange
	tests := []struct {
		date     time.Time
		id       string
		length   string // in seconds
		taskID   string
		taskName string

		expectedStatus  status
		expectedIssueID string
	}{
		{time.Now(), "ID", "600", "1", "[XX-123] name", Unknown, "XX-123"},
	}

	for _, theory := range tests {

		// act
		item := WorkLogItem{
			ID:       theory.id,
			Length:   theory.length,
			TaskID:   theory.taskID,
			TaskName: theory.taskName,
		}
		obj := NewTimeSheetEntry(theory.date, item)

		duration, _ := time.ParseDuration(theory.length + "s")
		assert.NotNil(t, obj)
		assert.Equal(t, theory.expectedStatus, obj.Status)
		assert.Equal(t, duration, obj.SpentTime)
		assert.Equal(t, theory.expectedIssueID, obj.YouTrackIssueID)

	}
}

func TestNewTimeSheetEntry_ShouldChangeStatusAsExpected(t *testing.T) {

	// act
	item := WorkLogItem{
		ID:       "1",
		Length:   "10",
		TaskID:   "theory.taskID",
		TaskName: "theory.taskName",
	}
	obj := NewTimeSheetEntry(time.Now(), item)

	obj.SetModified()
	assert.Equal(t, Modified, obj.Status)
	assert.Equal(t, ">", obj.Status.String())

	obj.SetNew()
	assert.Equal(t, New, obj.Status)
	assert.Equal(t, "+", obj.Status.String())

	obj.SetNotModified()
	assert.Equal(t, NotModified, obj.Status)
	assert.Equal(t, "=", obj.Status.String())

	obj.SetUnknown()
	assert.Equal(t, Unknown, obj.Status)
	assert.Equal(t, "?", obj.Status.String())
}
