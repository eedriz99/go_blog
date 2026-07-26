package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
	//"github.com/lib/pq"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, email, first_name, last_name, username, is_active FROM users WHERE id = $1`

	var m model.User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.Email,
		&m.FirstName,
		&m.LastName,
		&m.Username,
		&m.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return &m, nil
}

func (u *UserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, first_name, last_name, username, password, is_active FROM users WHERE email = $1`

	var m model.User
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&m.ID,
		&m.Email,
		&m.FirstName,
		&m.LastName,
		&m.Username,
		&m.Password,
		&m.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return &m, nil
}

func (u *UserStore) Update(ctx context.Context, payload payload.UpdateUserPayload) error {
	setParts := []string{}
	args := []any{}
	i := 1

	var m model.User

	if payload.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email=$%d", i))
		args = append(args, *payload.Email)
		i++
	}

	if payload.FirstName != nil {
		setParts = append(setParts, fmt.Sprintf("first_name=$%d", i))
		args = append(args, *payload.FirstName)
		i++
	}

	if payload.LastName != nil {
		setParts = append(setParts, fmt.Sprintf("last_name=$%d", i))
		args = append(args, *payload.LastName)
		i++
	}

	if payload.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username=$%d", i))
		args = append(args, *payload.Username)
		i++
	}

	if payload.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active=$%d", i))
		args = append(args, *payload.IsActive)
		i++
	}

	query := fmt.Sprintf(`
						UPDATE users SET %s 
						WHERE id=$%d 
						RETURNING id, email, username, first_name, last_name, is_active;
							`, strings.Join(setParts, ", "), i)
	args = append(args, payload.ID)
	return withTx(ctx, u.db, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, query, args...).Scan(
			&m.ID,
			&m.Email,
			&m.Username,
			&m.FirstName,
			&m.LastName,
			&m.IsActive,
		); err != nil {
			return err
		}
		return nil
	})

}

func (u *UserStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1;`

	return withTx(ctx, u.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, id)
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
	})
}

func (u *UserStore) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM user_invitations WHERE token = $1;`

	return withTx(ctx, u.db, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx, query, token)
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
	})

}

func (u *UserStore) CreateWithInvitation(ctx context.Context, m *model.User, expiry time.Duration) (*string, error) {
	token, err := withReturningTx(ctx, u.db, func(tx *sql.Tx) (*string, error) {
		err := u.create(ctx, tx, m)
		if err != nil {
			// _ = u.Delete(ctx, m.ID)
			return nil, err
		}

		token, err := u.createUserInvitation(ctx, tx, expiry, m.ID)
		if err != nil {
			return nil, err
		}

		return &token, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (u *UserStore) ActivateByInvitationToken(ctx context.Context, token string) error {
	return withTx(ctx, u.db, func(tx *sql.Tx) error {
		//  Get user associated with token

		user, err := u.getByInvitationToken(ctx, tx, token)
		if err != nil {
			return err
		}

		active := true
		var payload payload.UpdateUserPayload
		payload.ID = user.ID
		payload.IsActive = &active

		// update user's active state
		if err := u.Update(ctx, payload); err != nil {
			return err
		}

		// clean the invitation
		if err := u.DeleteToken(ctx, token); err != nil {
			return err
		}

		return nil

	})

}

// PRIVATE METHODS

func (u *UserStore) getByInvitationToken(ctx context.Context, tx *sql.Tx, token string) (*model.User, error) {
	query := `SELECT u.id, u.email, u.first_name, u.last_name, u.username, u.password
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token = $1 AND ui.expiry > NOW()`

	var m model.User
	err := tx.QueryRowContext(ctx, query, token).Scan(
		&m.ID,
		&m.Email,
		&m.FirstName,
		&m.LastName,
		&m.Username,
		&m.Password,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return &m, nil

}

func (u *UserStore) create(ctx context.Context, tx *sql.Tx, m *model.User) error {
	query := `INSERT INTO users (email, first_name, last_name, username, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := tx.QueryRowContext(ctx, query,
		m.Email,
		m.FirstName,
		m.LastName,
		m.Username,
		m.Password,
	).Scan(
		&m.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, expiry time.Duration, userID string) (string, error) {
	query := `INSERT INTO user_invitations (user_id, expiry) VALUES ($1, $2) RETURNING token`
	var token string
	if err := tx.QueryRowContext(ctx, query, userID, time.Now().Add(expiry)).Scan(&token); err != nil {
		return "", err
	}

	return token, nil

}
