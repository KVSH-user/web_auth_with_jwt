package errors

import "errors"

var (
	ErrUsernameTaken  = errors.New("username already taken")
	ErrIncorrectEmail = errors.New("incorrect email")
)
