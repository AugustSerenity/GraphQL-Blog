package memory

import (
	"context"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

func (r *Repository) CreatePost(ctx context.Context, post *model.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.ID] = post

	return nil
}

func (r *Repository) GetPost(ctx context.Context, id string) (*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return nil, ErrNotFound
	}

	return post, nil
}

func (r *Repository) ListPosts(ctx context.Context) ([]*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*model.Post, 0, len(r.posts))

	for _, post := range r.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *Repository) SetCommentsEnabled(
	ctx context.Context,
	postID string,
	enabled bool,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	post, ok := r.posts[postID]
	if !ok {
		return ErrNotFound
	}

	post.CommentsEnabled = enabled

	return nil
}
