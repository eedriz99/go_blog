package model

import (
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostWithMetadata struct {
	Post
	Username     string `json:"username,omitempty"`    // include author's username in the response
	UserAvatar   string `json:"user_avatar,omitempty"` // include author's avatar URL in the response
	CommentCount int    `json:"comment_count"`         // include comment count in the response
}
