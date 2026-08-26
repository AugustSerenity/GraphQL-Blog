package model

import "time"

type Comment struct {
	ID        string
	PostID    string
	AuthorID  string
	ParentID  *string
	Content   string
	CreatedAt time.Time
}
