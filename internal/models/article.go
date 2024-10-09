// internal/models/article.go
package models

import "time"

type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    int       `json:"users_id"`
	Author    Author    `json:"author"`
}

type Author struct {
	Name string `json:"name"`
}
