package memory

import (
	"context"
	"testing"
	"time"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
	"github.com/AugustSerenity/GraphQL-Blog/internal/repository"
)

func TestCreateAndGetPost(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	got, err := repo.GetPost(ctx, post.ID)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	if got.ID != post.ID {
		t.Errorf("ID = %q, want %q", got.ID, post.ID)
	}

	if got.AuthorID != post.AuthorID {
		t.Errorf(
			"AuthorID = %q, want %q",
			got.AuthorID,
			post.AuthorID,
		)
	}

	if got.Content != post.Content {
		t.Errorf(
			"Content = %q, want %q",
			got.Content,
			post.Content,
		)
	}

	if got.CommentsEnabled != post.CommentsEnabled {
		t.Errorf(
			"CommentsEnabled = %v, want %v",
			got.CommentsEnabled,
			post.CommentsEnabled,
		)
	}
}

func TestGetPostNotFound(t *testing.T) {
	repo := New()
	ctx := context.Background()

	_, err := repo.GetPost(ctx, "does-not-exist")

	if err != ErrNotFound {
		t.Fatalf(
			"GetPost() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}

func TestListPosts(t *testing.T) {
	repo := New()
	ctx := context.Background()

	posts := []*model.Post{
		{
			ID:              "post-1",
			AuthorID:        "user-1",
			Content:         "Post 1",
			CommentsEnabled: true,
			CreatedAt:       time.Now(),
		},
		{
			ID:              "post-2",
			AuthorID:        "user-1",
			Content:         "Post 2",
			CommentsEnabled: true,
			CreatedAt:       time.Now(),
		},
	}

	for _, post := range posts {
		if err := repo.CreatePost(ctx, post); err != nil {
			t.Fatalf("CreatePost() error = %v", err)
		}
	}

	got, err := repo.ListPosts(ctx)
	if err != nil {
		t.Fatalf("ListPosts() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf(
			"ListPosts() returned %d posts, want 2",
			len(got),
		)
	}
}

func TestSetCommentsEnabled(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	if err := repo.SetCommentsEnabled(ctx, post.ID, false); err != nil {
		t.Fatalf(
			"SetCommentsEnabled() error = %v",
			err,
		)
	}

	got, err := repo.GetPost(ctx, post.ID)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	if got.CommentsEnabled {
		t.Error("CommentsEnabled = true, want false")
	}

	if err := repo.SetCommentsEnabled(ctx, post.ID, true); err != nil {
		t.Fatalf(
			"SetCommentsEnabled() error = %v",
			err,
		)
	}

	got, err = repo.GetPost(ctx, post.ID)
	if err != nil {
		t.Fatalf("GetPost() error = %v", err)
	}

	if !got.CommentsEnabled {
		t.Error("CommentsEnabled = false, want true")
	}
}

func TestSetCommentsEnabledNotFound(t *testing.T) {
	repo := New()
	ctx := context.Background()

	err := repo.SetCommentsEnabled(
		ctx,
		"does-not-exist",
		false,
	)

	if err != ErrNotFound {
		t.Fatalf(
			"SetCommentsEnabled() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}

func TestCreateAndGetComment(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	comment := &model.Comment{
		ID:        "comment-1",
		PostID:    post.ID,
		AuthorID:  "user-2",
		Content:   "Test comment",
		CreatedAt: time.Now(),
	}

	if err := repo.CreateComment(ctx, comment); err != nil {
		t.Fatalf(
			"CreateComment() error = %v",
			err,
		)
	}

	got, err := repo.GetComment(ctx, comment.ID)
	if err != nil {
		t.Fatalf("GetComment() error = %v", err)
	}

	if got.ID != comment.ID {
		t.Errorf(
			"ID = %q, want %q",
			got.ID,
			comment.ID,
		)
	}

	if got.PostID != comment.PostID {
		t.Errorf(
			"PostID = %q, want %q",
			got.PostID,
			comment.PostID,
		)
	}

	if got.Content != comment.Content {
		t.Errorf(
			"Content = %q, want %q",
			got.Content,
			comment.Content,
		)
	}
}

func TestCreateCommentPostNotFound(t *testing.T) {
	repo := New()
	ctx := context.Background()

	comment := &model.Comment{
		ID:        "comment-1",
		PostID:    "does-not-exist",
		AuthorID:  "user-1",
		Content:   "Test comment",
		CreatedAt: time.Now(),
	}

	err := repo.CreateComment(ctx, comment)

	if err != ErrNotFound {
		t.Fatalf(
			"CreateComment() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}

func TestGetCommentNotFound(t *testing.T) {
	repo := New()
	ctx := context.Background()

	_, err := repo.GetComment(
		ctx,
		"does-not-exist",
	)

	if err != ErrNotFound {
		t.Fatalf(
			"GetComment() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}

func TestNestedComments(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	parent := &model.Comment{
		ID:        "comment-parent",
		PostID:    post.ID,
		AuthorID:  "user-1",
		Content:   "Parent",
		CreatedAt: time.Now(),
	}

	if err := repo.CreateComment(ctx, parent); err != nil {
		t.Fatalf(
			"CreateComment(parent) error = %v",
			err,
		)
	}

	child := &model.Comment{
		ID:        "comment-child",
		PostID:    post.ID,
		AuthorID:  "user-2",
		ParentID:  &parent.ID,
		Content:   "Child",
		CreatedAt: time.Now(),
	}

	if err := repo.CreateComment(ctx, child); err != nil {
		t.Fatalf(
			"CreateComment(child) error = %v",
			err,
		)
	}

	children, err := repo.GetCommentChildren(
		ctx,
		parent.ID,
	)
	if err != nil {
		t.Fatalf(
			"GetCommentChildren() error = %v",
			err,
		)
	}

	if len(children) != 1 {
		t.Fatalf(
			"GetCommentChildren() returned %d children, want 1",
			len(children),
		)
	}

	if children[0].ID != child.ID {
		t.Errorf(
			"child ID = %q, want %q",
			children[0].ID,
			child.ID,
		)
	}
}

func TestGetComments(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	for i := 1; i <= 3; i++ {
		comment := &model.Comment{
			ID:        "comment-" + string(rune('0'+i)),
			PostID:    post.ID,
			AuthorID:  "user-1",
			Content:   "Comment",
			CreatedAt: time.Now(),
		}

		if err := repo.CreateComment(ctx, comment); err != nil {
			t.Fatalf(
				"CreateComment() error = %v",
				err,
			)
		}
	}

	page, err := repo.GetComments(
		ctx,
		repository.CommentListParams{
			PostID: post.ID,
			Limit:  2,
		},
	)
	if err != nil {
		t.Fatalf(
			"GetComments() error = %v",
			err,
		)
	}

	if len(page.Items) != 2 {
		t.Fatalf(
			"GetComments() returned %d comments, want 2",
			len(page.Items),
		)
	}

	if !page.HasNextPage {
		t.Error("HasNextPage = false, want true")
	}

	if page.EndCursor == nil {
		t.Fatal("EndCursor is nil")
	}
}

func TestGetCommentsPagination(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	for i := 1; i <= 3; i++ {
		comment := &model.Comment{
			ID:        "comment-" + string(rune('0'+i)),
			PostID:    post.ID,
			AuthorID:  "user-1",
			Content:   "Comment",
			CreatedAt: time.Now(),
		}

		if err := repo.CreateComment(ctx, comment); err != nil {
			t.Fatalf(
				"CreateComment() error = %v",
				err,
			)
		}
	}

	firstPage, err := repo.GetComments(
		ctx,
		repository.CommentListParams{
			PostID: post.ID,
			Limit:  2,
		},
	)
	if err != nil {
		t.Fatalf(
			"GetComments(first page) error = %v",
			err,
		)
	}

	if firstPage.EndCursor == nil {
		t.Fatal("first page EndCursor is nil")
	}

	secondPage, err := repo.GetComments(
		ctx,
		repository.CommentListParams{
			PostID: post.ID,
			Limit:  2,
			Cursor: firstPage.EndCursor,
		},
	)
	if err != nil {
		t.Fatalf(
			"GetComments(second page) error = %v",
			err,
		)
	}

	if len(secondPage.Items) != 1 {
		t.Fatalf(
			"second page returned %d comments, want 1",
			len(secondPage.Items),
		)
	}

	if secondPage.HasNextPage {
		t.Error("second page HasNextPage = true, want false")
	}
}

func TestCreateCommentParentNotFound(t *testing.T) {
	repo := New()
	ctx := context.Background()

	post := &model.Post{
		ID:              "post-1",
		AuthorID:        "user-1",
		Content:         "Test post",
		CommentsEnabled: true,
		CreatedAt:       time.Now(),
	}

	if err := repo.CreatePost(ctx, post); err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}

	parentID := "does-not-exist"

	comment := &model.Comment{
		ID:        "comment-1",
		PostID:    post.ID,
		AuthorID:  "user-2",
		ParentID:  &parentID,
		Content:   "Child",
		CreatedAt: time.Now(),
	}

	err := repo.CreateComment(ctx, comment)

	if err != ErrNotFound {
		t.Fatalf(
			"CreateComment() error = %v, want %v",
			err,
			ErrNotFound,
		)
	}
}
