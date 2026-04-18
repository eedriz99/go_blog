package store

import (
	"context"
	"database/sql"
	"errors"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
)

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
		GetByPost(ctx context.Context, postID string) ([]CommentWithUsername, error)
		Update(context.Context, payload.UpdateCommentPayload) (*model.Comment, error)
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
