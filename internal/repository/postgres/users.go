package postgres

import (
	"context"
	"database/sql"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

func (r *Repository) GetUser(
	ctx context.Context,
	id string,
) (*model.User, error) {
	user := &model.User{}

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id
		FROM users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return user, nil
}
