package repository

import (
	"context"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

type Repository interface {
	PostRepository
	CommentRepository
	UserRepository
}

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPost(ctx context.Context, id string) (*model.Post, error)
	ListPosts(ctx context.Context) ([]*model.Post, error)
	SetCommentsEnabled(ctx context.Context, postID string, enabled bool) error
}

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetComment(ctx context.Context, id string) (*model.Comment, error)
	GetComments(ctx context.Context, params CommentListParams) (*CommentPage, error)
	GetCommentChildren(ctx context.Context, parentID string) ([]*model.Comment, error)
	GetCommentsChildren(ctx context.Context, parentIDs []string) (map[string][]*model.Comment, error)
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
}

type CommentListParams struct {
	PostID   string
	ParentID *string
	Limit    int
	Cursor   *string
}

type CommentPage struct {
	Items       []*model.Comment
	HasNextPage bool
	EndCursor   *string
}
