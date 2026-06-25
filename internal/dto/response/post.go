package dto

import (
	"time"

	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
)

type PostResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostWithCommentsResponse struct {
	ID        string              `json:"id"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	Tags      []string            `json:"tags"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Comments  CommentListResponse `json:"comments,omitempty"` // include comments in the response
}

type PostWithMetadataResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username,omitempty"`    // include author's username in the response
	UserAvatar   string    `json:"user_avatar,omitempty"` // include author's avatar URL in the response
	CommentCount int       `json:"comment_count"`         // include comment count in the response

}

type PostListResponse struct {
	Data  []PostResponse `json:"data"`
	Total int            `json:"total"`
}

func NewPostResponse(p *model.Post) PostResponse {
	return PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func NewListPostResponse(posts []model.Post) PostListResponse {
	data := make([]PostResponse, len(posts))

	for i := range posts {
		data[i] = NewPostResponse(&posts[i])
	}

	return PostListResponse{
		Data:  data,
		Total: len(data),
	}
}

func NewPostWithCommentsResponse(p PostResponse, comments []store.CommentWithUsername) PostWithCommentsResponse {
	commentResponses := NewCommentListResponse(comments)

	return PostWithCommentsResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Comments:  commentResponses,
	}
}
