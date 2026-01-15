package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

func generateID() string {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Error with rand function")
	}
	return hex.EncodeToString(bytes)
}

func findRunningTask(data *TimeData) *TimeEntry {
	for i := range data.Entries {
		if data.Entries[i].IsRunning() {
			return &data.Entries[i]
		}
	}
	return nil
}

// getSortedIndices returns the indices of entries sorted by start time (most recent first)
func getSortedIndices(entries []TimeEntry) []int {
	indices := make([]int, len(entries))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return entries[indices[i]].StartTime.After(entries[indices[j]].StartTime)
	})
	return indices
}

func StartTask(title string) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	running := findRunningTask(data)
	if running != nil {
		now := time.Now()
		running.EndTime = &now
		fmt.Printf("Stopped: %s (ran for %s)\n", running.Title, formatDuration(running.Duration()))
	}

	entry := TimeEntry{
		ID:        generateID(),
		Title:     title,
		StartTime: time.Now(),
	}

	data.Entries = append(data.Entries, entry)

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Started: %s [%s]\n", title, entry.ID)
	return nil
}

func StopTask() error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	running := findRunningTask(data)
	if running == nil {
		fmt.Println("No task is currently running")
		return nil
	}

	now := time.Now()
	running.EndTime = &now

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Stopped: %s (ran for %s)\n", running.Title, formatDuration(running.Duration()))
	return nil
}

func Status() error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	running := findRunningTask(data)
	if running == nil {
		fmt.Println("No task is currently running")
		return nil
	}

	fmt.Printf("Running: %s [%s]\n", running.Title, running.ID)
	fmt.Printf("Started: %s (%s ago)\n", running.StartTime.Format("15:04:05"), formatDuration(running.Duration()))
	return nil
}

func ListTasks(limit int) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if len(data.Entries) == 0 {
		fmt.Println("No time entries found")
		return nil
	}

	// Get indices sorted by start time, most recent first
	sortedIndices := getSortedIndices(data.Entries)

	// Apply limit
	totalEntries := len(sortedIndices)
	displayIndices := sortedIndices
	if limit > 0 && limit < totalEntries {
		displayIndices = sortedIndices[:limit]
	}

	fmt.Printf("%-5s %-30s %-20s %-20s %-10s\n", "IDX", "TITLE", "START", "END", "DURATION")
	fmt.Println(strings.Repeat("-", 90))

	for i, origIdx := range displayIndices {
		entry := data.Entries[origIdx]
		endStr := "running"
		if entry.EndTime != nil {
			endStr = entry.EndTime.Format("2006-01-02 15:04")
		}

		title := entry.Title
		if len(title) > 28 {
			title = title[:28] + ".."
		}

		fmt.Printf("%-5d %-30s %-20s %-20s %-10s\n",
			i,
			title,
			entry.StartTime.Format("2006-01-02 15:04"),
			endStr,
			formatDuration(entry.Duration()),
		)
	}

	if limit > 0 && totalEntries > limit {
		fmt.Printf("\nShowing %d of %d entries. Use -n <number> to show more.\n", limit, totalEntries)
	}

	return nil
}

func ViewTask(index int) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if len(data.Entries) == 0 {
		fmt.Println("No time entries found")
		return nil
	}

	sortedIndices := getSortedIndices(data.Entries)

	if index < 0 || index >= len(sortedIndices) {
		return fmt.Errorf("invalid index: %d (valid range: 0-%d)", index, len(sortedIndices)-1)
	}

	origIdx := sortedIndices[index]
	entry := data.Entries[origIdx]

	endStr := "running"
	if entry.EndTime != nil {
		endStr = entry.EndTime.Format("2006-01-02 15:04:05")
	}

	fmt.Printf("Index:    %d\n", index)
	fmt.Printf("ID:       %s\n", entry.ID)
	fmt.Printf("Title:    %s\n", entry.Title)
	fmt.Printf("Start:    %s\n", entry.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("End:      %s\n", endStr)
	fmt.Printf("Duration: %s\n", formatDuration(entry.Duration()))
	if entry.Notes != "" {
		fmt.Printf("Notes:\n%s\n", entry.Notes)
	}

	return nil
}

func DeleteTask(index int) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if index < 0 || index >= len(data.Entries) {
		return fmt.Errorf("invalid index: %d", index)
	}

	// Get the original index from sorted order (most recent first)
	sortedIndices := getSortedIndices(data.Entries)
	origIdx := sortedIndices[index]

	fmt.Printf("Deleted: %s\n", data.Entries[origIdx].Title)
	data.Entries = append(data.Entries[:origIdx], data.Entries[origIdx+1:]...)

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

func EditTask(index int, newTitle string, startAdjustMins int) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if index < 0 || index >= len(data.Entries) {
		return fmt.Errorf("invalid index: %d", index)
	}

	// Get the original index from sorted order (most recent first)
	sortedIndices := getSortedIndices(data.Entries)
	origIdx := sortedIndices[index]

	entry := &data.Entries[origIdx]

	if newTitle != "" {
		oldTitle := entry.Title
		entry.Title = newTitle
		fmt.Printf("Updated title: '%s' -> '%s'\n", oldTitle, newTitle)
	}

	if startAdjustMins != 0 {
		oldStart := entry.StartTime
		entry.StartTime = entry.StartTime.Add(time.Duration(startAdjustMins) * time.Minute)
		if entry.EndTime != nil && entry.StartTime.After(*entry.EndTime) {
			entry.StartTime = oldStart
			return fmt.Errorf("adjusted start time would be after end time")
		}
		fmt.Printf("Updated start: %s -> %s\n", oldStart.Format("15:04"), entry.StartTime.Format("15:04"))
	}

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

func NoteTask(index int, note string) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if index < 0 || index >= len(data.Entries) {
		return fmt.Errorf("invalid index: %d", index)
	}

	sortedIndices := getSortedIndices(data.Entries)
	origIdx := sortedIndices[index]

	entry := &data.Entries[origIdx]

	if entry.Notes == "" {
		entry.Notes = note
	} else {
		entry.Notes = entry.Notes + "\n" + note
	}

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Added note to: %s\n", entry.Title)
	return nil
}

func Summary(filter string) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if len(data.Entries) == 0 {
		fmt.Println("No time entries found")
		return nil
	}

	now := time.Now()
	var startFilter, endFilter time.Time
	filterLabel := "All time"

	switch filter {
	case "today":
		startFilter = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		filterLabel = "Today"
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startFilter = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		filterLabel = "This week"
	case "last":
		lastDay := findLastWorkingDay(data.Entries, now)
		if lastDay.IsZero() {
			fmt.Println("No entries found before today")
			return nil
		}
		startFilter = lastDay
		endFilter = lastDay.AddDate(0, 0, 1)
		filterLabel = "Last working day (" + lastDay.Format("Mon 2 Jan") + ")"
	}

	var totalDuration time.Duration
	taskDurations := make(map[string]time.Duration)
	taskNotes := make(map[string][]string)
	count := 0

	for _, entry := range data.Entries {
		if !startFilter.IsZero() && entry.StartTime.Before(startFilter) {
			continue
		}
		if !endFilter.IsZero() && !entry.StartTime.Before(endFilter) {
			continue
		}
		duration := entry.Duration()
		totalDuration += duration
		taskDurations[entry.Title] += duration
		if entry.Notes != "" {
			taskNotes[entry.Title] = append(taskNotes[entry.Title], entry.Notes)
		}
		count++
	}

	if count == 0 {
		fmt.Println("No entries found for the selected period")
		return nil
	}

	fmt.Printf("=== %s Summary ===\n\n", filterLabel)
	fmt.Printf("Total time: %s (%d entries)\n\n", formatDuration(totalDuration), count)

	fmt.Println("By task:")
	fmt.Println(strings.Repeat("-", 50))

	for title, duration := range taskDurations {
		fmt.Printf("%s: %s\n", title, formatDuration(duration))
		if notes, ok := taskNotes[title]; ok {
			for _, note := range notes {
				for _, line := range strings.Split(note, "\n") {
					fmt.Printf("  - %s\n", line)
				}
			}
		}
	}

	return nil
}

func findLastWorkingDay(entries []TimeEntry, now time.Time) time.Time {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var lastDay time.Time
	for _, entry := range entries {
		entryDay := time.Date(entry.StartTime.Year(), entry.StartTime.Month(), entry.StartTime.Day(), 0, 0, 0, 0, entry.StartTime.Location())
		if entryDay.Before(today) && entryDay.After(lastDay) {
			lastDay = entryDay
		}
	}
	return lastDay
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
