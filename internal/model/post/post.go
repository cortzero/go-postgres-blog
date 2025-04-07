package post

import "time"

type Post struct {
	ID        uint      `json:"id,omitempty"`
	UserID    uint      `json:"user_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
