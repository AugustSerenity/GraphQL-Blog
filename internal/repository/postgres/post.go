package postgres

import (
	"context"
	"database/sql"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

func (r *Repository) CreatePost(
	ctx context.Context,
	post *model.Post,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO posts (
			id,
			author_id,
			content,
			comments_enabled,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		post.ID,
		post.AuthorID,
		post.Content,
		post.CommentsEnabled,
		post.CreatedAt,
	)

	return err
}

func (r *Repository) GetPost(
	ctx context.Context,
	id string,
) (*model.Post, error) {
	post := &model.Post{}

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			author_id,
			content,
			comments_enabled,
			created_at
		FROM posts
		WHERE id = $1
		`,
		id,
	).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Content,
		&post.CommentsEnabled,
		&post.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return post, nil
}

func (r *Repository) ListPosts(
	ctx context.Context,
) ([]*model.Post, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id,
			author_id,
			content,
			comments_enabled,
			created_at
		FROM posts
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := make([]*model.Post, 0)

	for rows.Next() {
		post := &model.Post{}

		if err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Content,
			&post.CommentsEnabled,
			&post.CreatedAt,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) SetCommentsEnabled(
	ctx context.Context,
	postID string,
	enabled bool,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`
		UPDATE posts
		SET comments_enabled = $1
		WHERE id = $2
		`,
		enabled,
		postID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
