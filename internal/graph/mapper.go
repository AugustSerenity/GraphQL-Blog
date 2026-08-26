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

	result := &model.Comment{
		ID:      comment.ID,
		Content: comment.Content,
	}

	if comment.ParentID != nil {
		result.Parent = &model.Comment{
			ID: *comment.ParentID,
		}
	}

	return result
}
