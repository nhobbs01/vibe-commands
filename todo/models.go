package main

import "time"

type TodoItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoData struct {
	Items []TodoItem `json:"items"`
}
