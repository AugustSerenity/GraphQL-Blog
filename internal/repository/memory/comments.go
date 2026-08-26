package memory

import (
	"context"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
)

func (r *Repository) CreateComment(
	ctx context.Context,
	comment *model.Comment,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[comment.PostID]; !ok {
		return ErrNotFound
	}

	if comment.ParentID != nil {
		if _, ok := r.comments[*comment.ParentID]; !ok {
			return ErrNotFound
		}
	}

	r.comments[comment.ID] = comment
	r.commentsByPost[comment.PostID] = append(
		r.commentsByPost[comment.PostID],
		comment.ID,
	)

	if comment.ParentID != nil {
		parentID := *comment.ParentID

		r.commentsByParent[parentID] = append(
			r.commentsByParent[parentID],
			comment.ID,
		)
	}

	return nil
}

func (r *Repository) GetComment(
	ctx context.Context,
	id string,
) (*model.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, ok := r.comments[id]
	if !ok {
		return nil, ErrNotFound
	}

	return comment, nil
}

func (r *Repository) GetComments(
	ctx context.Context,
	params repository.CommentListParams,
) (*repository.CommentPage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.commentsByPost[params.PostID]

	if params.ParentID != nil {
		ids = r.commentsByParent[*params.ParentID]
	}

	start := 0

	if params.Cursor != nil {
		for i, id := range ids {
			if id == *params.Cursor {
				start = i + 1
				break
			}
		}
	}

	limit := params.Limit
	if limit <= 0 {
		limit = 5
	}

	end := start + limit
	if end > len(ids) {
		end = len(ids)
	}

	items := make([]*model.Comment, 0, end-start)

	for _, id := range ids[start:end] {
		items = append(items, r.comments[id])
	}

	hasNextPage := end < len(ids)

	var endCursor *string

	if len(items) > 0 {
		cursor := items[len(items)-1].ID
		endCursor = &cursor
	}

	return &repository.CommentPage{
		Items:       items,
		HasNextPage: hasNextPage,
		EndCursor:   endCursor,
	}, nil
}
