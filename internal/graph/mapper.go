package graph

import (
	"github.com/AugustSerenity/GraphQL-Blog/internal/graph/model"
	domain "github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

func postToGraphQL(post *domain.Post) *model.Post {
	if post == nil {
		return nil
	}

	return &model.Post{
		ID:              post.ID,
		Content:         post.Content,
		CommentsEnabled: post.CommentsEnabled,
	}
}

func commentToGraphQL(comment *domain.Comment) *model.Comment {
	if comment == nil {
		return nil
	}

	return &model.Comment{
		ID:      comment.ID,
		Content: comment.Content,
	}
}

func userToGraphQL(user *domain.User) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		ID: user.ID,
	}
}
