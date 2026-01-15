package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Define subcommand flag sets
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listLimit := listCmd.Int("n", 10, "number of entries to show (0 for all)")
	viewCmd := flag.NewFlagSet("view", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editTitle := editCmd.String("title", "", "new title for the entry")
	editStart := editCmd.Int("start", 0, "adjust start time by minutes (negative = earlier)")

	noteCmd := flag.NewFlagSet("note", flag.ExitOnError)

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)
	summaryToday := summaryCmd.Bool("today", false, "show today's summary")
	summaryWeek := summaryCmd.Bool("week", false, "show this week's summary")
	summaryLast := summaryCmd.Bool("last", false, "show last working day's summary")

	var err error
	command := os.Args[1]

	switch command {
	case "start":
		startCmd.Parse(os.Args[2:])
		args := startCmd.Args()
		var title string
		if len(args) > 0 {
			title = strings.Join(args, " ")
		} else {
			// Check if stdin is a pipe (not a terminal)
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				reader := bufio.NewReader(os.Stdin)
				title, _ = reader.ReadString('\n')
				title = strings.TrimSpace(title)
			}
		}
		if title == "" {
			fmt.Println("Error: missing task title")
			fmt.Println("Usage: timetrack start <title>")
			os.Exit(1)
		}
		err = StartTask(title)

	case "stop":
		stopCmd.Parse(os.Args[2:])
		err = StopTask()

	case "status":
		statusCmd.Parse(os.Args[2:])
		err = Status()

	case "list":
		listCmd.Parse(os.Args[2:])
		err = ListTasks(*listLimit)

	case "view":
		viewCmd.Parse(os.Args[2:])
		args := viewCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: missing entry index")
			fmt.Println("Usage: timetrack view <index>")
			os.Exit(1)
		}
		index, parseErr := strconv.Atoi(args[0])
		if parseErr != nil {
			fmt.Println("Error: index must be a number")
			os.Exit(1)
		}
		err = ViewTask(index)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		args := deleteCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: missing entry index")
			fmt.Println("Usage: timetrack delete <index>")
			os.Exit(1)
		}
		index, parseErr := strconv.Atoi(args[0])
		if parseErr != nil {
			fmt.Println("Error: index must be a number")
			os.Exit(1)
		}
		err = DeleteTask(index)

	case "edit":
		editCmd.Parse(os.Args[2:])
		args := editCmd.Args()
		if *editTitle == "" && *editStart == 0 {
			fmt.Println("Error: must specify --title or --start")
			fmt.Println("Usage: timetrack edit [--title \"new title\"] [--start <mins>] <index>")
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Println("Error: missing entry index")
			fmt.Println("Usage: timetrack edit [--title \"new title\"] [--start <mins>] <index>")
			os.Exit(1)
		}
		index, parseErr := strconv.Atoi(args[0])
		if parseErr != nil {
			fmt.Println("Error: index must be a number")
			os.Exit(1)
		}
		err = EditTask(index, *editTitle, *editStart)

	case "note":
		noteCmd.Parse(os.Args[2:])
		args := noteCmd.Args()
		if len(args) < 2 {
			fmt.Println("Error: missing index and/or note text")
			fmt.Println("Usage: timetrack note <index> \"note text\"")
			os.Exit(1)
		}
		index, parseErr := strconv.Atoi(args[0])
		if parseErr != nil {
			fmt.Println("Error: index must be a number")
			os.Exit(1)
		}
		noteText := strings.Join(args[1:], " ")
		err = NoteTask(index, noteText)

	case "summary":
		summaryCmd.Parse(os.Args[2:])
		filter := ""
		if *summaryToday {
			filter = "today"
		} else if *summaryWeek {
			filter = "week"
		} else if *summaryLast {
			filter = "last"
		}
		err = Summary(filter)

	case "help", "--help", "-h":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`timetrack - Simple time tracking CLI

Usage:
  timetrack <command> [arguments]

Commands:
  start <title>              Start a new task (auto-stops current task)
  stop                       Stop the current running task
  status                     Show the current running task
  list [-n <limit>]          List time entries (default: 10, most recent first)
  view <index>               View full details of an entry
  delete <index>             Delete an entry by index
  edit [--title <title>] [--start <mins>] <index>
                           Edit an entry (--start -30 = started 30 mins earlier)
  note <index> <text>        Add a note to an entry (appends if note exists)
  summary [--today|--week|--last]
                             Show time summary (--last = last day with entries)

Examples:
  timetrack start "Working on feature X"
  timetrack stop
  timetrack list
  timetrack list -n 20    # Show 20 entries
  timetrack list -n 0     # Show all entries
  timetrack note 0 "Fixed the login bug"
  timetrack summary --today`)
}
