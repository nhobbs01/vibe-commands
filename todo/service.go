package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func generateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func AddItem(title string) error {
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	item := TodoItem{
		ID:        generateID(),
		Title:     title,
		CreatedAt: time.Now(),
	}

	data.Items = append(data.Items, item)

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Added: %s [%s]\n", title, item.ID)
	return nil
}

func ListItems() error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if len(data.Items) == 0 {
		fmt.Println("No items")
		return nil
	}

	fmt.Printf("%-5s %s\n", "IDX", "TITLE")
	fmt.Println("---   -----")

	for i, item := range data.Items {
		title := item.Title
		if len(title) > 60 {
			title = title[:58] + ".."
		}
		fmt.Printf("%-5d %s\n", i, title)
	}

	return nil
}

func DoneItem(index int) error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if index < 0 || index >= len(data.Items) {
		return fmt.Errorf("invalid index: %d", index)
	}

	item := data.Items[index]
	data.Items = append(data.Items[:index], data.Items[index+1:]...)

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Done: %s\n", item.Title)
	return nil
}

func NextItem() error {
	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if len(data.Items) == 0 {
		return fmt.Errorf("no items")
	}

	fmt.Print(data.Items[0].Title)
	return nil
}

func EditItem(index int, newTitle string) error {
	if newTitle == "" {
		return fmt.Errorf("title cannot be empty")
	}

	data, err := LoadData()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	if index < 0 || index >= len(data.Items) {
		return fmt.Errorf("invalid index: %d", index)
	}

	oldTitle := data.Items[index].Title
	data.Items[index].Title = newTitle

	if err := SaveData(data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	fmt.Printf("Updated: '%s' -> '%s'\n", oldTitle, newTitle)
	return nil
}
