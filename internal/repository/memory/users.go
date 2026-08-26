package memory

import (
	"context"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

func (r *Repository) GetUser(
	ctx context.Context,
	id string,
) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	return user, nil
}
