package graph

import (
	"context"
	"fmt"

	"github.com/AugustSerenity/GraphQL-Blog/internal/graph/model"
)

func (r *mutationResolver) CreatePost(
	ctx context.Context,
	input model.CreatePostInput,
) (*model.Post, error) {
	fmt.Println("=== CREATE POST RESOLVER CALLED ===")
	fmt.Println("content:", input.Content)

	post, err := r.Service.CreatePost(
		ctx,
		"user-1",
		input.Content,
	)
	if err != nil {
		fmt.Println("=== CREATE POST ERROR ===", err)
		return nil, err
	}

	fmt.Println("=== POST CREATED ===")
	fmt.Println("postID:", post.ID)
	fmt.Println("authorID:", post.AuthorID)

	return postToGraphQL(post), nil
}

func (r *mutationResolver) CreateComment(
	ctx context.Context,
	input model.CreateCommentInput,
) (*model.Comment, error) {
	fmt.Println("=== CREATE COMMENT RESOLVER CALLED ===")
	fmt.Println("postID:", input.PostID)
	fmt.Println("content:", input.Content)

	comment, err := r.Service.CreateComment(
		ctx,
		"user-1",
		input.PostID,
		input.ParentID,
		input.Content,
	)
	if err != nil {
		fmt.Println("=== CREATE COMMENT ERROR ===", err)
		return nil, err
	}

	fmt.Println("=== COMMENT CREATED ===")
	fmt.Println("commentID:", comment.ID)
	fmt.Println("postID:", comment.PostID)
	fmt.Println("authorID:", comment.AuthorID)

	return commentToGraphQL(comment), nil
}

func (r *mutationResolver) SetCommentsEnabled(
	ctx context.Context,
	postID string,
	enabled bool,
) (*model.Post, error) {
	fmt.Println("=== SET COMMENTS ENABLED ===")
	fmt.Println("postID:", postID)
	fmt.Println("enabled:", enabled)

	post, err := r.Service.SetCommentsEnabled(
		ctx,
		"user-1",
		postID,
		enabled,
	)
	if err != nil {
		fmt.Println("=== SET COMMENTS ERROR ===", err)
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
	fmt.Println("========================================")
	fmt.Println("=== COMMENT ADDED RESOLVER CALLED ===")
	fmt.Println("postID:", postID)
	fmt.Println("========================================")

	events, cancel := r.Service.SubscribeComments(postID)

	fmt.Println("=== SUBSCRIBED TO COMMENT BROKER ===")
	fmt.Println("postID:", postID)

	go func() {
		<-ctx.Done()

		fmt.Println("=== SUBSCRIPTION CONTEXT DONE ===")
		fmt.Println("postID:", postID)

		cancel()
	}()

	result := make(chan *model.Comment)

	go func() {
		defer close(result)

		fmt.Println("=== SUBSCRIPTION EVENT LOOP STARTED ===")
		fmt.Println("postID:", postID)

		for comment := range events {
			fmt.Println("=== EVENT RECEIVED BY RESOLVER ===")
			fmt.Println("commentID:", comment.ID)
			fmt.Println("postID:", comment.PostID)
			fmt.Println("content:", comment.Content)

			select {
			case result <- commentToGraphQL(comment):
				fmt.Println("=== EVENT SENT TO GRAPHQL ===")
				fmt.Println("commentID:", comment.ID)

			case <-ctx.Done():
				fmt.Println("=== CONTEXT DONE WHILE SENDING EVENT ===")
				return
			}
		}

		fmt.Println("=== SUBSCRIPTION EVENT CHANNEL CLOSED ===")
		fmt.Println("postID:", postID)
	}()

	return result, nil
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
