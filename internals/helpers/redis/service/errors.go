package service

import "errors"

var (
	ErrFailedToGetLock = errors.New("failed to get lock")
	ErrInvalidReply    = errors.New("invalid reply")
)
