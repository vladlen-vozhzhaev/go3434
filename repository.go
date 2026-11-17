package main

import (
	"context"
	"time"
)

func getAllPosts(ctx context.Context) ([]Post, error) {
	rows, err := db.Query(ctx, "SELECT * FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func createPost(ctx context.Context, title, content string) error {
	_, err := db.Exec(ctx,
		"INSERT INTO posts (title, content, created_at) VALUES ($1, $2, $3)",
		title, content, time.Now(),
	)
	return err
}

func getPostByID(ctx context.Context, id int) (Post, error) {
	var p Post

	err := db.QueryRow(ctx,
		"SELECT * FROM posts WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt)
	return p, err
}

func updatePost(ctx context.Context, id int, title, content string) error {
	_, err := db.Exec(ctx,
		"UPDATE posts SET title = $1, content = $2 WHERE id = $3",
		title, content, id,
	)
	return err
}

func deletePost(ctx context.Context, id int) error {
	_, err := db.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	return err
}
