# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test Commands

```bash
go build -o timetrack       # Build the binary
go test ./...               # Run all tests
go test -v ./...            # Run tests with verbose output
go test -run TestName       # Run a specific test
./install.sh                # Build and install to /usr/local/bin (requires sudo)
```

## Architecture

This is a simple CLI time tracking application written in Go. Data is stored as JSON in `~/.timetrack.json`.

**File structure:**
- `main.go` - CLI entry point with subcommand parsing using `flag.FlagSet`
- `service.go` - Business logic for all commands (start, stop, list, edit, delete, summary)
- `models.go` - Data types (`TimeEntry`, `TimeData`) with helper methods
- `storage.go` - JSON file persistence (`LoadData`, `SaveData`)

**Key patterns:**
- Entries are stored in an array but displayed sorted by start time (most recent first)
- The `list` command shows a display index (0-based, sorted order) that maps to original array indices via `getSortedIndices()`
- Only one task can be running at a time; starting a new task auto-stops the current one
- Running tasks have `EndTime == nil`
