package domain

import "errors"

var (
	ErrNotSubscribed     = errors.New("not subscribed")
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrAddressNotFound   = errors.New("address not found")
)
