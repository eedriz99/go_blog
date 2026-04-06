package store

import (
	"context"
	"database/sql"

	"github.com/eedriz99/go_blog/internal/model"
	//"github.com/lib/pq"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) Create(ctx context.Context, m *model.User) error {
	query := `INSERT INTO users (email, first_name, last_name, username) VALUES ($1, $2, $3) RETURNING id`

	err := u.db.QueryRowContext(ctx, query,
		m.Email,
		m.FirstName,
		m.LastName,
		m.Username,
	).Scan(
		&m.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) GetByID(ctx context.Context, payload string) (*model.User, error) {
	return nil, nil
}

func (u *UserStore) Update(ctx context.Context, m *model.User) error {
	return nil
}
func (u *UserStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1;`

	result, err := u.db.ExecContext(ctx, query, id)
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
