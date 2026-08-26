package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
	"github.com/google/uuid"
)

const MaxCommentLength = 2000

type Service struct {
	repo   repository.Repository
	broker *CommentBroker
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo:   repo,
		broker: NewCommentBroker(),
	}
}

func (s *Service) CreatePost(
	ctx context.Context,
	authorID string,
	content string,
) (*model.Post, error) {
	if strings.TrimSpace(content) == "" {
		return nil, ErrInvalidContent
	}

	post := &model.Post{
		ID:              uuid.NewString(),
		AuthorID:        authorID,
		Content:         content,
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) CreateComment(
	ctx context.Context,
	authorID string,
	postID string,
	parentID *string,
	content string,
) (*model.Comment, error) {
	fmt.Println("========================================")
	fmt.Println("=== CREATE COMMENT SERVICE START ===")
	fmt.Println("postID:", postID)
	fmt.Println("content:", content)

	content = strings.TrimSpace(content)

	if content == "" {
		fmt.Println("=== INVALID CONTENT ===")
		return nil, ErrInvalidContent
	}

	if len([]rune(content)) > MaxCommentLength {
		fmt.Println("=== COMMENT TOO LONG ===")
		return nil, ErrCommentTooLong
	}

	fmt.Println("=== GET POST ===")

	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		fmt.Println("=== GET POST ERROR ===")
		fmt.Println("error:", err)
		return nil, err
	}

	fmt.Println("=== POST FOUND ===")
	fmt.Println("postID:", post.ID)
	fmt.Println("commentsEnabled:", post.CommentsEnabled)

	if !post.CommentsEnabled {
		fmt.Println("=== COMMENTS DISABLED ===")
		return nil, ErrCommentsDisabled
	}

	if parentID != nil {
		fmt.Println("=== GET PARENT COMMENT ===")
		fmt.Println("parentID:", *parentID)

		parent, err := s.repo.GetComment(ctx, *parentID)
		if err != nil {
			fmt.Println("=== GET PARENT COMMENT ERROR ===")
			fmt.Println("error:", err)
			return nil, err
		}

		if parent.PostID != postID {
			return nil, errors.New("parent comment belongs to another post")
		}
	}

	comment := &model.Comment{
		ID:        uuid.NewString(),
		PostID:    postID,
		AuthorID:  authorID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	fmt.Println("=== COMMENT CREATED IN MEMORY ===")
	fmt.Println("commentID:", comment.ID)

	fmt.Println("=== BEFORE REPO CREATE COMMENT ===")

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		fmt.Println("=== REPO CREATE COMMENT ERROR ===")
		fmt.Println("error:", err)
		return nil, err
	}

	fmt.Println("=== COMMENT SAVED TO REPO ===")
	fmt.Println("commentID:", comment.ID)
	fmt.Println("postID:", comment.PostID)

	fmt.Println("=== BEFORE BROKER PUBLISH ===")
	fmt.Println("postID:", postID)
	fmt.Println("commentID:", comment.ID)

	s.broker.Publish(postID, comment)

	fmt.Println("=== AFTER BROKER PUBLISH ===")
	fmt.Println("postID:", postID)
	fmt.Println("commentID:", comment.ID)

	fmt.Println("=== CREATE COMMENT SERVICE END ===")
	fmt.Println("========================================")

	return comment, nil
}

func (s *Service) SetCommentsEnabled(
	ctx context.Context,
	userID string,
	postID string,
	enabled bool,
) (*model.Post, error) {
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != userID {
		return nil, ErrNotPostAuthor
	}

	if err := s.repo.SetCommentsEnabled(ctx, postID, enabled); err != nil {
		return nil, err
	}

	post.CommentsEnabled = enabled

	return post, nil
}

func (s *Service) ListPosts(
	ctx context.Context,
) ([]*model.Post, error) {
	return s.repo.ListPosts(ctx)
}

func (s *Service) GetPost(
	ctx context.Context,
	id string,
) (*model.Post, error) {
	return s.repo.GetPost(ctx, id)
}

func (s *Service) GetComments(
	ctx context.Context,
	postID string,
	limit int,
	cursor *string,
) (*repository.CommentPage, error) {
	return s.repo.GetComments(ctx, repository.CommentListParams{
		PostID: postID,
		Limit:  limit,
		Cursor: cursor,
	})
}

func (s *Service) GetComment(
	ctx context.Context,
	id string,
) (*model.Comment, error) {
	return s.repo.GetComment(ctx, id)
}

func (s *Service) GetCommentChildren(
	ctx context.Context,
	parentID string,
) ([]*model.Comment, error) {
	return s.repo.GetCommentChildren(ctx, parentID)
}

func (s *Service) GetCommentsChildren(
	ctx context.Context,
	parentIDs []string,
) (map[string][]*model.Comment, error) {
	return s.repo.GetCommentsChildren(ctx, parentIDs)
}

func (s *Service) SubscribeComments(
	postID string,
) (<-chan *model.Comment, func(), error) {
	fmt.Println("========================================")
	fmt.Println("=== SERVICE SUBSCRIBE COMMENTS ===")
	fmt.Println("postID:", postID)

	_, err := s.repo.GetPost(context.Background(), postID)
	if err != nil {
		fmt.Println("=== SUBSCRIBE POST NOT FOUND ===")
		fmt.Println("postID:", postID)
		return nil, nil, err
	}

	ch := s.broker.Subscribe(postID)

	cancel := func() {
		s.broker.Unsubscribe(postID, ch)
	}

	fmt.Println("=== SERVICE SUBSCRIBE COMMENTS DONE ===")
	fmt.Println("postID:", postID)

	return ch, cancel, nil
}
