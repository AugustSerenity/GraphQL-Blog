package model

import "time"

type Post struct {
	ID              string
	AuthorID        string
	Content         string
	CommentsEnabled bool
	CreatedAt       time.Time
}
