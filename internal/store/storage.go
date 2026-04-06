package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
)

// ===========================> Payloads & Responses structs <==============================
// ===========================> Create Payloads <==============================
type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type CreateCommentPayload struct {
	Content string `json:"content"`
	PostID  string `json:"post_id"`
}

// ===========================> Get Payloads <==============================
type GetPostPayload struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type GetCommentPayload struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	PostID string `json:"post_id"`
}

type GetPostCommentsRes struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

// ===========================> Update Payloads <==============================
type UpdatePostPayload struct {
	ID      string    `json:"id"`
	Title   *string   `json:"title,omitempty"`
	Content *string   `json:"content,omitempty"`
	Tags    *[]string `json:"tags,omitempty"`
}

type UpdateCommentPayload struct {
	ID      string `json:"id"`
	PostID  string `json:"post_id"`
	UserID  string `json:user_id`
	Content string `json:"content"`
}

// ===========================> Delete Payloads <==============================

type DeletePayload struct {
	ID string `json:"id"`
}

var (
	ErrorNotFound   = errors.New("resource not found")
	ErrorInternal   = errors.New("internal server error")
	ErrorBadRequest = errors.New("bad request")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *model.Post) error
		GetByID(context.Context, string) (*model.Post, error)
		GetAll(context.Context, string) ([]model.Post, error)
		Update(context.Context, payload.UpdatePostPayload) (*model.Post, error)
		Delete(context.Context, payload.DeletePostPayload) error
	}

	Users interface {
		Create(context.Context, *model.User) error
		GetByID(context.Context, string) (*model.User, error)
		Update(context.Context, *model.User) error
		Delete(context.Context, string) error
	}

	Comments interface {
		Create(context.Context, *model.Comment) error
		Update(context.Context, UpdateCommentPayload) (*model.Comment, error)
		Delete(context.Context, string) error
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
