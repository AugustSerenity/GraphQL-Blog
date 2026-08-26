package service

import (
	"context"
	"testing"
	"time"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository/memory"
)

func newTestService() *Service {
	repo := memory.New()
	return NewService(repo)
}

func TestCreatePost(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	if post == nil {
		t.Fatal("CreatePost() returned nil post")
	}

	if post.ID == "" {
		t.Error("post ID is empty")
	}

	if post.Content != "Test post" {
		t.Errorf(
			"post content = %q, want %q",
			post.Content,
			"Test post",
		)
	}

	if post.AuthorID != "user-1" {
		t.Errorf(
			"post authorID = %q, want %q",
			post.AuthorID,
			"user-1",
		)
	}

	if !post.CommentsEnabled {
		t.Error("comments should be enabled by default")
	}

	if post.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestCreatePostEmptyContent(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.CreatePost(
		ctx,
		"user-1",
		"   ",
	)

	if err != ErrInvalidContent {
		t.Fatalf(
			"CreatePost() error = %v, want %v",
			err,
			ErrInvalidContent,
		)
	}
}

func TestCreateComment(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	comment, err := svc.CreateComment(
		ctx,
		"user-2",
		post.ID,
		nil,
		"Test comment",
	)
	if err != nil {
		t.Fatalf("CreateComment() error = %v", err)
	}

	if comment == nil {
		t.Fatal("CreateComment() returned nil comment")
	}

	if comment.ID == "" {
		t.Error("comment ID is empty")
	}

	if comment.PostID != post.ID {
		t.Errorf(
			"comment postID = %q, want %q",
			comment.PostID,
			post.ID,
		)
	}

	if comment.AuthorID != "user-2" {
		t.Errorf(
			"comment authorID = %q, want %q",
			comment.AuthorID,
			"user-2",
		)
	}

	if comment.Content != "Test comment" {
		t.Errorf(
			"comment content = %q, want %q",
			comment.Content,
			"Test comment",
		)
	}

	if comment.ParentID != nil {
		t.Error("top-level comment should have nil ParentID")
	}

	if comment.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestCreateCommentPostNotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.CreateComment(
		ctx,
		"user-1",
		"does-not-exist",
		nil,
		"Test comment",
	)

	if err == nil {
		t.Fatal("CreateComment() expected error, got nil")
	}

	if err != memory.ErrNotFound {
		t.Fatalf(
			"CreateComment() error = %v, want %v",
			err,
			memory.ErrNotFound,
		)
	}
}

func TestCreateCommentEmptyContent(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	_, err = svc.CreateComment(
		ctx,
		"user-2",
		post.ID,
		nil,
		"   ",
	)

	if err != ErrInvalidContent {
		t.Fatalf(
			"CreateComment() error = %v, want %v",
			err,
			ErrInvalidContent,
		)
	}
}

func TestCreateCommentTooLong(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	content := make([]rune, MaxCommentLength+1)

	for i := range content {
		content[i] = 'a'
	}

	_, err = svc.CreateComment(
		ctx,
		"user-2",
		post.ID,
		nil,
		string(content),
	)

	if err != ErrCommentTooLong {
		t.Fatalf(
			"CreateComment() error = %v, want %v",
			err,
			ErrCommentTooLong,
		)
	}
}

func TestCreateCommentCommentsDisabled(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	_, err = svc.SetCommentsEnabled(
		ctx,
		"user-1",
		post.ID,
		false,
	)
	if err != nil {
		t.Fatalf("SetCommentsEnabled() error = %v", err)
	}

	_, err = svc.CreateComment(
		ctx,
		"user-2",
		post.ID,
		nil,
		"Test comment",
	)

	if err != ErrCommentsDisabled {
		t.Fatalf(
			"CreateComment() error = %v, want %v",
			err,
			ErrCommentsDisabled,
		)
	}
}

func TestCreateNestedComment(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	parent, err := svc.CreateComment(
		ctx,
		"user-2",
		post.ID,
		nil,
		"Parent comment",
	)
	if err != nil {
		t.Fatalf("CreateComment(parent) error = %v", err)
	}

	child, err := svc.CreateComment(
		ctx,
		"user-3",
		post.ID,
		&parent.ID,
		"Child comment",
	)
	if err != nil {
		t.Fatalf("CreateComment(child) error = %v", err)
	}

	if child.ParentID == nil {
		t.Fatal("child ParentID is nil")
	}

	if *child.ParentID != parent.ID {
		t.Errorf(
			"child ParentID = %q, want %q",
			*child.ParentID,
			parent.ID,
		)
	}
}

func TestCreateNestedCommentWrongPost(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post1, err := svc.CreatePost(
		ctx,
		"user-1",
		"Post 1",
	)
	if err != nil {
		t.Fatalf("CreatePost(post1) error = %v", err)
	}

	post2, err := svc.CreatePost(
		ctx,
		"user-1",
		"Post 2",
	)
	if err != nil {
		t.Fatalf("CreatePost(post2) error = %v", err)
	}

	parent, err := svc.CreateComment(
		ctx,
		"user-2",
		post1.ID,
		nil,
		"Parent",
	)
	if err != nil {
		t.Fatalf("CreateComment(parent) error = %v", err)
	}

	_, err = svc.CreateComment(
		ctx,
		"user-3",
		post2.ID,
		&parent.ID,
		"Child",
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "parent comment belongs to another post" {
		t.Fatalf(
			"error = %q, want %q",
			err.Error(),
			"parent comment belongs to another post",
		)
	}
}

func TestSetCommentsEnabled(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	if !post.CommentsEnabled {
		t.Fatal("comments should initially be enabled")
	}

	updated, err := svc.SetCommentsEnabled(
		ctx,
		"user-1",
		post.ID,
		false,
	)
	if err != nil {
		t.Fatalf("SetCommentsEnabled() error = %v", err)
	}

	if updated.CommentsEnabled {
		t.Error("comments should be disabled")
	}

	updated, err = svc.SetCommentsEnabled(
		ctx,
		"user-1",
		post.ID,
		true,
	)
	if err != nil {
		t.Fatalf("SetCommentsEnabled() error = %v", err)
	}

	if !updated.CommentsEnabled {
		t.Error("comments should be enabled")
	}
}

func TestSetCommentsEnabledWrongAuthor(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	_, err = svc.SetCommentsEnabled(
		ctx,
		"user-2",
		post.ID,
		false,
	)

	if err != ErrNotPostAuthor {
		t.Fatalf(
			"SetCommentsEnabled() error = %v, want %v",
			err,
			ErrNotPostAuthor,
		)
	}
}

func TestSubscribeCommentsPostNotFound(t *testing.T) {
	svc := newTestService()

	_, _, err := svc.SubscribeComments("does-not-exist")

	if err != memory.ErrNotFound {
		t.Fatalf(
			"SubscribeComments() error = %v, want %v",
			err,
			memory.ErrNotFound,
		)
	}
}

func TestSubscribeCommentsReceivesComment(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	post, err := svc.CreatePost(
		ctx,
		"user-1",
		"Test post",
	)
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	events, cancel, err := svc.SubscribeComments(post.ID)
	if err != nil {
		t.Fatalf("SubscribeComments() error = %v", err)
	}
	defer cancel()

	_, err = svc.CreateComment(
		ctx,
		"user-2",
		post.ID,
		nil,
		"Subscription test",
	)
	if err != nil {
		t.Fatalf("CreateComment() error = %v", err)
	}

	select {
	case comment := <-events:
		if comment == nil {
			t.Fatal("received nil comment")
		}

		if comment.PostID != post.ID {
			t.Errorf(
				"comment postID = %q, want %q",
				comment.PostID,
				post.ID,
			)
		}

		if comment.Content != "Subscription test" {
			t.Errorf(
				"comment content = %q, want %q",
				comment.Content,
				"Subscription test",
			)
		}

	case <-time.After(time.Second):
		t.Fatal("timed out waiting for subscription event")
	}
}

var _ *model.Comment
