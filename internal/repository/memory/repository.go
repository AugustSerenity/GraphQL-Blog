package memory

import (
	"sync"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
)

var _ repository.Repository = (*Repository)(nil)

type Repository struct {
	mu sync.RWMutex

	posts    map[string]*model.Post
	comments map[string]*model.Comment
	users    map[string]*model.User

	commentsByPost   map[string][]string
	commentsByParent map[string][]string
}

func New() *Repository {
	return &Repository{
		posts:            make(map[string]*model.Post),
		comments:         make(map[string]*model.Comment),
		users:            make(map[string]*model.User),
		commentsByPost:   make(map[string][]string),
		commentsByParent: make(map[string][]string),
	}
}
