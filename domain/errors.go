package domain

import "errors"

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidLogin            = errors.New("invalid login")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrCantParseJSON           = errors.New("can't parse JSON")
	ErrInternalServerError     = errors.New("internal server error 500")
	ErrLoginIsBusy             = errors.New("login is busy")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserUnauthorized        = errors.New("user unauthorized")
	ErrNameIsEmpty             = errors.New("name is empty")
	ErrMimeIsEmpty             = errors.New("mime is empty")
	ErrInvalidMetaData         = errors.New("invalid meta data")
	ErrFileNotFound            = errors.New("file not found")
	ErrFileIsTooLarge          = errors.New("file is too large")
	ErrParameterIsEmpty        = errors.New("parameter is empty")
	ErrDocumentNotFound        = errors.New("document not found")
	ErrFileIsDamagedOrNotFound = errors.New("file is damaged or not found")
	ErrFileThisNameIsAlready   = errors.New("file this name is already")
)
