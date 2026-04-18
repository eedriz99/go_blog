package dto

import (
	"time"

	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
)

type CommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Username  string `json:"username,omitempty"` // optional field for response
}

type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Total    int               `json:"total"`
}

func NewCommentResponse(model *store.CommentWithUsername) CommentResponse {
	return CommentResponse{
		ID:        model.ID,
		PostID:    model.PostID,
		Content:   model.Content,
		CreatedAt: model.CreatedAt.Format(time.RFC3339),
		UpdatedAt: model.UpdatedAt.Format(time.RFC3339),
		Username:  model.Username,
	}
}

func NewCommentListResponse(models []store.CommentWithUsername) CommentListResponse {
	res := CommentListResponse{
		Comments: make([]CommentResponse, len(models)),
		Total:    len(models),
	}

	for i := range models {
		res.Comments[i] = NewCommentResponse(&models[i])
	}
	return res
}

func NewCommentWithoutUsernameResponse(model *model.Comment) CommentResponse {
	return CommentResponse{
		ID:        model.ID,
		PostID:    model.PostID,
		Content:   model.Content,
		CreatedAt: model.CreatedAt.Format(time.RFC3339),
		UpdatedAt: model.UpdatedAt.Format(time.RFC3339),
	}
}
