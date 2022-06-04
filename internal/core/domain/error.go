package domain

import "errors"

var (
	ErrValue             = errors.New("invalid argument")
	ErrExists            = errors.New("the key already exists")
	ErrKey               = errors.New("the key does not exist")
	ErrExpired           = errors.New("the key has expired")
	ErrSessionTerminated = errors.New("the session had been terminated")
)
