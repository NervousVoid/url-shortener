package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
	AliasExists    = errors.New("alias exists")
)

type URLStorage interface {
	SaveURL(urlToSave string, alias string) error
	GetURL(alias string) (string, error)
	GetAlias(url string) (string, error)
}
