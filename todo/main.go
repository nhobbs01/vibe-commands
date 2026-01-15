package main

import (
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

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	nextCmd := flag.NewFlagSet("next", flag.ExitOnError)
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)

	var err error

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		args := addCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: title required")
			fmt.Println("Usage: todo add <title>")
			os.Exit(1)
		}
		title := strings.Join(args, " ")
		err = AddItem(title)

	case "list":
		listCmd.Parse(os.Args[2:])
		err = ListItems()

	case "done":
		doneCmd.Parse(os.Args[2:])
		args := doneCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: index required")
			fmt.Println("Usage: todo done <index>")
			os.Exit(1)
		}
		index, parseErr := strconv.Atoi(args[0])
		if parseErr != nil {
			fmt.Printf("Error: invalid index '%s'\n", args[0])
			os.Exit(1)
		}
		err = DoneItem(index)

	case "next":
		nextCmd.Parse(os.Args[2:])
		err = NextItem()

	case "edit":
		editCmd.Parse(os.Args[2:])
		args := editCmd.Args()
		if len(args) < 2 {
			fmt.Println("Error: index and title required")
			fmt.Println("Usage: todo edit <index> <title>")
			os.Exit(1)
		}
		index, parseErr := strconv.Atoi(args[0])
		if parseErr != nil {
			fmt.Printf("Error: invalid index '%s'\n", args[0])
			os.Exit(1)
		}
		title := strings.Join(args[1:], " ")
		err = EditItem(index, title)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: todo <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <title>          Add a new item")
	fmt.Println("  list                 List all items")
	fmt.Println("  done <index>         Remove an item")
	fmt.Println("  next                 Output the next item title (for piping)")
	fmt.Println("  edit <index> <title> Edit an item's title")
}
