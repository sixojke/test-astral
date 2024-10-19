package domain

import "errors"

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidLogin        = errors.New("invalid login")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrCantParseJSON       = errors.New("can't parse JSON")
	ErrInternalServerError = errors.New("internal server error 500")
	ErrLoginIsBusy         = errors.New("login is busy")
)
