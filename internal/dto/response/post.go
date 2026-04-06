package dto

import (
	"time"

	"github.com/eedriz99/go_blog/internal/model"
)

type PostResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
