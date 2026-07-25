package coro

import "errors"

var (
	ErrAlreadyActivated = errors.New("already activated")
	ErrNotActivated     = errors.New("not activated")
	ErrNoResult         = errors.New("no result")
	ErrNextTooSoon      = errors.New("called next too soon")
)
