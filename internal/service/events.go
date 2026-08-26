package service

import (
	"fmt"
	"sync"

	"github.com/AugustSerenity/GraphQL-Blog/internal/model"
)

type CommentBroker struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan *model.Comment]struct{}
}

func NewCommentBroker() *CommentBroker {
	return &CommentBroker{
		subscribers: make(map[string]map[chan *model.Comment]struct{}),
	}
}

func (b *CommentBroker) Subscribe(postID string) chan *model.Comment {
	ch := make(chan *model.Comment, 10)

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.subscribers[postID] == nil {
		b.subscribers[postID] = make(map[chan *model.Comment]struct{})
	}

	b.subscribers[postID][ch] = struct{}{}

	fmt.Println("=== SUBSCRIBE ===")
	fmt.Println("postID:", postID)
	fmt.Println("subscribers:", len(b.subscribers[postID]))

	return ch
}

func (b *CommentBroker) Unsubscribe(
	postID string,
	ch chan *model.Comment,
) {
	b.mu.Lock()
	defer b.mu.Unlock()

	subscribers, ok := b.subscribers[postID]
	if !ok {
		return
	}

	if _, exists := subscribers[ch]; !exists {
		return
	}

	delete(subscribers, ch)

	if len(subscribers) == 0 {
		delete(b.subscribers, postID)
	}

	close(ch)
}

func (b *CommentBroker) Publish(
	postID string,
	comment *model.Comment,
) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	fmt.Println("=== PUBLISH ===")
	fmt.Println("postID:", postID)
	fmt.Println("commentID:", comment.ID)
	fmt.Println("subscribers:", len(b.subscribers[postID]))

	for ch := range b.subscribers[postID] {
		fmt.Println("sending to subscriber")

		select {
		case ch <- comment:
			fmt.Println("sent")
		default:
			fmt.Println("channel full")
		}
	}
}
