package models

import "time"

type Article struct {
	ID         int64     `json:"id"`
	AuthorID   int64     `json:"author_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CategoryID int64     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"update_at"`
}
