package service

import "errors"

var (
	ErrInvalidContent   = errors.New("content is invalid")
	ErrCommentTooLong   = errors.New("comment is too long")
	ErrCommentsDisabled = errors.New("comments are disabled")
	ErrNotPostAuthor    = errors.New("user is not post author")
)
