package storage

import (
	"fmt"
	"sync"
)

type InMemStorage struct {
	Mu      sync.RWMutex
	Urls    map[string]string
	Aliases map[string]string
}

func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		Mu:      sync.RWMutex{},
		Urls:    make(map[string]string),
		Aliases: make(map[string]string),
	}
}

func (ims *InMemStorage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.inmem.SaveURL"

	ims.Mu.RLock()
	_, ok := ims.Urls[urlToSave]
	ims.Mu.RUnlock()

	if ok {
		return fmt.Errorf("%s: %w", op, ErrURLExists)
	}

	ims.Mu.RLock()
	_, ok = ims.Aliases[alias]
	ims.Mu.RUnlock()

	if ok {
		return fmt.Errorf("%s: %w", op, AliasExists)
	}

	ims.Mu.Lock()
	ims.Aliases[alias] = urlToSave
	ims.Urls[urlToSave] = alias
	ims.Mu.Unlock()

	return nil
}

func (ims *InMemStorage) GetURL(alias string) (string, error) {
	const op = "storage.inmem.GetURL"

	ims.Mu.RLock()
	_, ok := ims.Aliases[alias]
	ims.Mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrURLNotFound)
	}

	var url string
	ims.Mu.RLock()
	url = ims.Aliases[alias]
	ims.Mu.RUnlock()

	return url, nil
}

func (ims *InMemStorage) GetAlias(url string) (string, error) {
	const op = "storage.inmem.GetURL"

	ims.Mu.RLock()
	_, ok := ims.Urls[url]
	ims.Mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("%s: %w", op, ErrURLNotFound)
	}

	var alias string
	ims.Mu.RLock()
	alias = ims.Urls[url]
	ims.Mu.RUnlock()

	return alias, nil
}
