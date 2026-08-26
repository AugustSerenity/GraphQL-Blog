package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
)

func (r *Repository) CreateComment(
	ctx context.Context,
	comment *model.Comment,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO comments (
			id,
			post_id,
			author_id,
			parent_id,
			content,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		comment.ID,
		comment.PostID,
		comment.AuthorID,
		comment.ParentID,
		comment.Content,
		comment.CreatedAt,
	)

	return err
}

func (r *Repository) GetComment(
	ctx context.Context,
	id string,
) (*model.Comment, error) {
	comment := &model.Comment{}

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			id,
			post_id,
			author_id,
			parent_id,
			content,
			created_at
		FROM comments
		WHERE id = $1
		`,
		id,
	).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.ParentID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return comment, nil
}

func (r *Repository) GetComments(
	ctx context.Context,
	params repository.CommentListParams,
) (*repository.CommentPage, error) {
	limit := params.Limit

	if limit <= 0 {
		limit = 5
	}

	query := `
		SELECT
			id,
			post_id,
			author_id,
			parent_id,
			content,
			created_at
		FROM comments
		WHERE post_id = $1
	`

	args := []any{
		params.PostID,
	}

	argIndex := 2

	if params.ParentID != nil {
		query += `
			AND parent_id = $` + strconv.Itoa(argIndex)

		args = append(args, *params.ParentID)
		argIndex++
	} else {
		query += `
			AND parent_id IS NULL
		`
	}

	if params.Cursor != nil {
		query += `
			AND created_at < (
				SELECT created_at
				FROM comments
				WHERE id = $` + strconv.Itoa(argIndex) + `
			)
		`

		args = append(args, *params.Cursor)
		argIndex++
	}

	query += `
		ORDER BY created_at DESC, id DESC
		LIMIT $` + strconv.Itoa(argIndex)

	args = append(args, limit+1)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]*model.Comment, 0, limit+1)

	for rows.Next() {
		comment := &model.Comment{}

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.ParentID,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	hasNextPage := len(comments) > limit

	if hasNextPage {
		comments = comments[:limit]
	}

	var endCursor *string

	if len(comments) > 0 {
		cursor := comments[len(comments)-1].ID
		endCursor = &cursor
	}

	return &repository.CommentPage{
		Items:       comments,
		HasNextPage: hasNextPage,
		EndCursor:   endCursor,
	}, nil
}

func (r *Repository) GetCommentChildren(
	ctx context.Context,
	parentID string,
) ([]*model.Comment, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id,
			post_id,
			author_id,
			parent_id,
			content,
			created_at
		FROM comments
		WHERE parent_id = $1
		ORDER BY created_at ASC, id ASC
		`,
		parentID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]*model.Comment, 0)

	for rows.Next() {
		comment := &model.Comment{}

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.ParentID,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *Repository) GetCommentsChildren(
	ctx context.Context,
	parentIDs []string,
) (map[string][]*model.Comment, error) {
	result := make(map[string][]*model.Comment, len(parentIDs))

	if len(parentIDs) == 0 {
		return result, nil
	}

	placeholders := make([]string, len(parentIDs))
	args := make([]any, len(parentIDs))

	for i, id := range parentIDs {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		args[i] = id

		result[id] = make([]*model.Comment, 0)
	}

	query := `
		SELECT
			id,
			post_id,
			author_id,
			parent_id,
			content,
			created_at
		FROM comments
		WHERE parent_id IN (` + strings.Join(placeholders, ", ") + `)
		ORDER BY created_at ASC, id ASC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		comment := &model.Comment{}

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.ParentID,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}

		if comment.ParentID != nil {
			result[*comment.ParentID] = append(
				result[*comment.ParentID],
				comment,
			)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
