package storage

import "errors"

var (
	// ErrURLNotFound throw this in case of selecting non-contained url from db
	ErrURLNotFound = errors.New("url not found")

	// ErrURLExists throw this in case of inserting new url to db, where the url already exists
	ErrURLExists = errors.New("url exists")
)
