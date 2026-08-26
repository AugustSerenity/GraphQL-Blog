package service

import (
	"testing"
	"time"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

func TestCommentBrokerSubscribeAndPublish(t *testing.T) {
	broker := NewCommentBroker()

	ch := broker.Subscribe("post-1")

	comment := &model.Comment{
		ID:      "comment-1",
		PostID:  "post-1",
		Content: "hello",
	}

	broker.Publish("post-1", comment)

	select {
	case got := <-ch:
		if got.ID != comment.ID {
			t.Fatalf(
				"comment ID = %q, want %q",
				got.ID,
				comment.ID,
			)
		}

		if got.Content != comment.Content {
			t.Fatalf(
				"comment content = %q, want %q",
				got.Content,
				comment.Content,
			)
		}

	case <-time.After(time.Second):
		t.Fatal("timeout waiting for comment event")
	}
}

func TestCommentBrokerDoesNotPublishToAnotherPost(t *testing.T) {
	broker := NewCommentBroker()

	ch := broker.Subscribe("post-1")

	comment := &model.Comment{
		ID:      "comment-1",
		PostID:  "post-2",
		Content: "hello",
	}

	broker.Publish("post-2", comment)

	select {
	case got := <-ch:
		t.Fatalf(
			"received unexpected comment: %v",
			got.ID,
		)

	case <-time.After(100 * time.Millisecond):
		// expected: no event
	}
}

func TestCommentBrokerMultipleSubscribers(t *testing.T) {
	broker := NewCommentBroker()

	ch1 := broker.Subscribe("post-1")
	ch2 := broker.Subscribe("post-1")

	comment := &model.Comment{
		ID:      "comment-1",
		PostID:  "post-1",
		Content: "hello",
	}

	broker.Publish("post-1", comment)

	for i, ch := range []chan *model.Comment{ch1, ch2} {
		select {
		case got := <-ch:
			if got.ID != comment.ID {
				t.Fatalf(
					"subscriber %d: comment ID = %q, want %q",
					i+1,
					got.ID,
					comment.ID,
				)
			}

		case <-time.After(time.Second):
			t.Fatalf(
				"subscriber %d: timeout waiting for comment",
				i+1,
			)
		}
	}
}

func TestCommentBrokerUnsubscribe(t *testing.T) {
	broker := NewCommentBroker()

	ch := broker.Subscribe("post-1")

	broker.Unsubscribe("post-1", ch)

	comment := &model.Comment{
		ID:      "comment-1",
		PostID:  "post-1",
		Content: "hello",
	}

	broker.Publish("post-1", comment)

	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("received event after unsubscribe")
		}

	case <-time.After(100 * time.Millisecond):
		t.Fatal("channel was not closed after unsubscribe")
	}
}
