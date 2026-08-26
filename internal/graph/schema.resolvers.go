package graph

import (
	"context"

	"github.com/AugustSerenity/GraphQL-Blog/internal/graph/model"
)

func (r *mutationResolver) CreatePost(
	ctx context.Context,
	input model.CreatePostInput,
) (*model.Post, error) {
	post, err := r.Service.CreatePost(
		ctx,
		"user-1",
		input.Content,
	)
	if err != nil {
		return nil, err
	}

	return postToGraphQL(post), nil
}

func (r *mutationResolver) CreateComment(
	ctx context.Context,
	input model.CreateCommentInput,
) (*model.Comment, error) {
	comment, err := r.Service.CreateComment(
		ctx,
		"user-1",
		input.PostID,
		input.ParentID,
		input.Content,
	)
	if err != nil {
		return nil, err
	}

	return commentToGraphQL(comment), nil
}

func (r *mutationResolver) SetCommentsEnabled(
	ctx context.Context,
	postID string,
	enabled bool,
) (*model.Post, error) {
	post, err := r.Service.SetCommentsEnabled(
		ctx,
		"user-1",
		postID,
		enabled,
	)
	if err != nil {
		return nil, err
	}

	return postToGraphQL(post), nil
}

func (r *queryResolver) Posts(
	ctx context.Context,
) ([]*model.Post, error) {
	posts, err := r.Service.ListPosts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Post, 0, len(posts))

	for _, post := range posts {
		result = append(result, postToGraphQL(post))
	}

	return result, nil
}

func (r *queryResolver) Post(
	ctx context.Context,
	id string,
) (*model.Post, error) {
	post, err := r.Service.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}

	return postToGraphQL(post), nil
}

func (r *queryResolver) Comments(
	ctx context.Context,
	postID string,
	pagination *model.CommentsPagination,
) (*model.CommentConnection, error) {
	var limit int
	var cursor *string

	if pagination != nil {
		if pagination.Limit != nil {
			limit = *pagination.Limit
		}

		cursor = pagination.Cursor
	}

	page, err := r.Service.GetComments(
		ctx,
		postID,
		limit,
		cursor,
	)
	if err != nil {
		return nil, err
	}

	parentIDs := make([]string, 0, len(page.Items))

	for _, comment := range page.Items {
		parentIDs = append(parentIDs, comment.ID)
	}

	children, err := r.Service.GetCommentsChildren(ctx, parentIDs)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Comment, 0, len(page.Items))

	for _, comment := range page.Items {
		item := commentToGraphQL(comment)

		for _, child := range children[comment.ID] {
			item.Children = append(
				item.Children,
				commentToGraphQL(child),
			)
		}

		items = append(items, item)
	}

	return &model.CommentConnection{
		Items: items,
		PageInfo: &model.PageInfo{
			HasNextPage: page.HasNextPage,
			EndCursor:   page.EndCursor,
		},
	}, nil
}

func (r *subscriptionResolver) CommentAdded(
	ctx context.Context,
	postID string,
) (<-chan *model.Comment, error) {
	return nil, nil
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type (
	mutationResolver     struct{ *Resolver }
	queryResolver        struct{ *Resolver }
	subscriptionResolver struct{ *Resolver }
)
