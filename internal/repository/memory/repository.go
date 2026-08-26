package memory

import (
	"sync"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

type Repository struct {
	mu sync.RWMutex

	posts    map[string]*model.Post
	comments map[string]*model.Comment
	users    map[string]*model.User
}

func New() *Repository {
	return &Repository{
		posts:    make(map[string]*model.Post),
		comments: make(map[string]*model.Comment),
		users:    make(map[string]*model.User),
	}
}
