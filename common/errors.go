package common

import "errors"

var (
	ErrInvalidMethodFormat  = errors.New("invalid method format")
	ErrMethodNotImplemented = errors.New("method not implemented")
)
