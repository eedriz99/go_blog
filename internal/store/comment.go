package store

import (
	"context"
	"database/sql"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
)

type CommentStore struct {
	db *sql.DB
}

func (c *CommentStore) Create(ctx context.Context, m *model.Comment) error {
	query := `
		INSERT INTO comments (content, user_id, post_id)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`

	err := c.db.QueryRowContext(
		ctx, query,
		m.Content,
		m.UserID,
		m.PostID,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (c *CommentStore) GetByPost(ctx context.Context, postID string) ([]CommentWithUsername, error) {
	query := `
			SELECT c.id, c.post_id, c.content,c.created_at, c.updated_at, u.username 
			FROM comments AS c
			JOIN users AS u ON c.user_id = u.id
			WHERE post_id = $1
			ORDER BY c.created_at DESC;`

	rows, err := c.db.QueryContext(ctx, query, postID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []CommentWithUsername{}
	for rows.Next() {
		var m CommentWithUsername // use the mediator struct to include username
		if err := rows.Scan(
			&m.ID,
			&m.PostID,
			&m.Content,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Username,
		); err != nil {
			return nil, err
		}

		comments = append(comments, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentStore) Update(ctx context.Context, comment payload.UpdateCommentPayload) (*model.Comment, error) {
	query := `
				UPDATE comments SET content = $1, update_at = NOW()
				WHERE id = &2, user_id = $3 AND post_id = $4
				RETURNING id, content, user_id, post_id, created_at, update_at;
	`
	var m model.Comment
	err := c.db.QueryRowContext(ctx, query,
		&comment.Content,
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
	).Scan(&m.ID,
		&m.Content,
		&m.UserID,
		&m.PostID,
		&m.CreatedAt,
		&m.UpdatedAt)

	if err != nil {
		return nil, ErrorInternal
	}
	return &m, nil
}

func (c *CommentStore) Delete(ctx context.Context, id string) error {
	return nil
}

func (c *CommentStore) GetByUser(ctx context.Context, userID string) ([]model.Comment, error) {
	return nil, nil
}

// # Context-aware access (primary path)
// GET  /posts/{post_id}/comments/       # comments for this post  -- Done
// POST /posts/{post_id}/comments/       # create (post_id inferred from URL) -- Done

// # Cross-cutting access (secondary path)
// GET  /comments/{id}/                  # direct access by ID
// PUT  /comments/{id}/                  # edit (no need to know post)
// DELETE /comments/{id}/                # delete (no need to know post)
// GET  /users/{user_id}/comments/       # user's comment history
