package main

import "time"

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
