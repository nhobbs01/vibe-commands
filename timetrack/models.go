package main

import "time"

type TimeEntry struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

type TimeData struct {
	Entries []TimeEntry `json:"entries"`
}

func (e *TimeEntry) IsRunning() bool {
	return e.EndTime == nil
}

func (e *TimeEntry) Duration() time.Duration {
	if e.EndTime == nil {
		return time.Since(e.StartTime)
	}
	return e.EndTime.Sub(e.StartTime)
}
