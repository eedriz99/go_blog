package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	payload "github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (p *PostStore) Create(ctx context.Context, m *model.Post) error {
	query := `
				INSERT 
				INTO posts (content, title, user_id, tags)
				VALUES ($1, $2, $3, $4) 
				RETURNING id, created_at, updated_at;
				`

	err := p.db.QueryRowContext(ctx, query,
		m.Content,
		m.Title,
		m.UserID,
		pq.Array(m.Tags),
	).Scan(
		&m.ID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostStore) GetByID(ctx context.Context, payload string) (*model.Post, error) {
	query := `
				SELECT id, title, content, tags, updated_at 
				FROM posts 
				WHERE id = $1;
				`
	var m model.Post

	err := p.db.QueryRowContext(ctx, query, payload).Scan(
		&m.ID,
		&m.Title,
		&m.Content,
		pq.Array(&m.Tags),
		&m.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			log.Println("DB Error: ", err)
			return nil, ErrorInternal
		}
	}
	return &m, nil
}

func (p *PostStore) GetAll(ctx context.Context, payload string) ([]model.Post, error) {
	query := `
				SELECT id, title, content, tags, created_at, updated_at 
				FROM posts 
				WHERE user_id=$1 
				ORDER BY updated_at DESC;
				`
	posts := []model.Post{}

	rows, err := p.db.QueryContext(ctx, query, payload)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m model.Post
		if err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Content,
			pq.Array(&m.Tags),
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostStore) Update(ctx context.Context, payload payload.UpdatePostPayload) (*model.Post, error) {

	setParts := []string{}
	args := []any{}
	i := 1

	var m model.Post

	m.UserID = ""

	if payload.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title=$%d", i))
		args = append(args, *payload.Title)
		i++
	}

	if payload.Content != nil {
		setParts = append(setParts, fmt.Sprintf("content=$%d", i))
		args = append(args, *payload.Content)
		i++
	}

	if payload.Tags != nil {
		setParts = append(setParts, fmt.Sprintf("tags=$%d", i))
		args = append(args, pq.Array(*payload.Tags))
		i++
	}

	if len(setParts) == 0 {
		return nil, ErrorBadRequest
	}

	query := fmt.Sprintf(`
						UPDATE posts SET %s, 
						updated_at = NOW() WHERE id=$%d 
						RETURNING id,title, content, tags, created_at, updated_at;
							`, strings.Join(setParts, ", "), i)
	args = append(args, payload.ID)

	// log.Printf("Args: %v", args)
	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&m.ID,
		&m.Title,
		&m.Content,
		pq.Array(&m.Tags),
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		// log.Printf("store Error: %v", err.Error())
		return nil, err
	}

	return &m, nil
}

func (p *PostStore) Delete(ctx context.Context, payload payload.DeletePostPayload) error {

	query := `DELETE FROM posts WHERE id = $1 AND user_id = $2;`

	result, err := p.db.ExecContext(ctx, query, payload.ID, payload.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrorNotFound
	}

	return nil
}
