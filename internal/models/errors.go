package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvaildCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplcate email")
)
