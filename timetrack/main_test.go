package main

import (
	"os"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{5 * time.Second, "5s"},
		{65 * time.Second, "1m 5s"},
		{3600 * time.Second, "1h 0m"},
		{3665 * time.Second, "1h 1m"},
		{7200*time.Second + 30*time.Minute, "2h 30m"},
	}

	for _, tt := range tests {
		result := formatDuration(tt.duration)
		if result != tt.expected {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
		}
	}
}

func TestTimeEntryIsRunning(t *testing.T) {
	entry := TimeEntry{
		ID:        "test123",
		Title:     "Test task",
		StartTime: time.Now(),
		EndTime:   nil,
	}

	if !entry.IsRunning() {
		t.Error("Expected entry with nil EndTime to be running")
	}

	now := time.Now()
	entry.EndTime = &now

	if entry.IsRunning() {
		t.Error("Expected entry with EndTime to not be running")
	}
}

func TestTimeEntryDuration(t *testing.T) {
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now()

	entry := TimeEntry{
		ID:        "test123",
		Title:     "Test task",
		StartTime: start,
		EndTime:   &end,
	}

	duration := entry.Duration()
	if duration < 59*time.Minute || duration > 61*time.Minute {
		t.Errorf("Expected duration around 1 hour, got %v", duration)
	}
}

// Helper to set up a temp data file for integration tests
func setupTestStorage(t *testing.T) (cleanup func()) {
	t.Helper()

	tmpDir := t.TempDir()

	// Override the home directory for tests
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	return func() {
		os.Setenv("HOME", origHome)
	}
}

func TestStorageLoadSave(t *testing.T) {
	cleanup := setupTestStorage(t)
	defer cleanup()

	// Load empty data
	data, err := LoadData()
	if err != nil {
		t.Fatalf("LoadData() error = %v", err)
	}
	if len(data.Entries) != 0 {
		t.Errorf("Expected empty entries, got %d", len(data.Entries))
	}

	// Add an entry and save
	now := time.Now()
	data.Entries = append(data.Entries, TimeEntry{
		ID:        "abc123",
		Title:     "Test task",
		StartTime: now,
		EndTime:   nil,
	})

	err = SaveData(data)
	if err != nil {
		t.Fatalf("SaveData() error = %v", err)
	}

	// Load again and verify
	data2, err := LoadData()
	if err != nil {
		t.Fatalf("LoadData() after save error = %v", err)
	}
	if len(data2.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(data2.Entries))
	}
	if data2.Entries[0].ID != "abc123" {
		t.Errorf("Expected ID 'abc123', got %q", data2.Entries[0].ID)
	}
	if data2.Entries[0].Title != "Test task" {
		t.Errorf("Expected title 'Test task', got %q", data2.Entries[0].Title)
	}
}

func TestFindRunningTask(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)

	data := &TimeData{
		Entries: []TimeEntry{
			{ID: "1", Title: "Completed", StartTime: past, EndTime: &now},
			{ID: "2", Title: "Running", StartTime: now, EndTime: nil},
		},
	}

	running := findRunningTask(data)
	if running == nil {
		t.Fatal("Expected to find running task")
	}
	if running.ID != "2" {
		t.Errorf("Expected running task ID '2', got %q", running.ID)
	}

	// Test with no running tasks
	data.Entries[1].EndTime = &now
	running = findRunningTask(data)
	if running != nil {
		t.Error("Expected no running task")
	}
}
