package errors

import "errors"

var (
	ErrEmptyMessage = errors.New("message is empty")
	ErrNotFound     = errors.New("entity not found")
)
