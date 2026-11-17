package main

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
}
