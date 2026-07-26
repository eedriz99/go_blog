package store

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/eedriz99/go_blog/internal/dto/payload"
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
		GetByUserID(context.Context, string) ([]model.Post, error)
		GetAll(context.Context) ([]model.Post, error)
		Update(context.Context, payload.UpdatePostPayload) (*model.Post, error)
		Delete(context.Context, string) error
	}

	Users interface {
		CreateWithInvitation(context.Context, *model.User, time.Duration) (*string, error)
		GetByID(context.Context, string) (*model.User, error)
		GetByEmail(context.Context, string) (*model.User, error)
		ActivateByInvitationToken(context.Context, string) error
		Update(context.Context, payload.UpdateUserPayload) error
		Delete(context.Context, string) error
		DeleteToken(context.Context, string) error
	}

	Comments interface {
		Create(context.Context, *model.Comment) error
		GetByPostID(ctx context.Context, postID string) ([]CommentWithUsername, error)
		Update(context.Context, payload.UpdateCommentPayload) (*model.Comment, error)
		Delete(context.Context, payload.DeleteCommentPayload) error
		GetByUser(context.Context, string) ([]model.Comment, error)
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}

func withTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("rollback failed: %v", rbErr)
		}
		return err
	}

	return tx.Commit()
}

func withReturningTx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) (*string, error)) (*string, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	result, err := fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return nil, rbErr
		}
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}
