package storage

import "image"

func NewInMemoryStorage() InMemoryStorage {
	return InMemoryStorage{
		store: make(map[string]image.Image),
	}
}

type InMemoryStorage struct {
	store map[string]image.Image
}

func (fsc InMemoryStorage) Get(key string) image.Image {
	return fsc.store[key]
}

func (fsc InMemoryStorage) Set(key string, thumbnail image.Image) (image.Image, error) {
	fsc.store[key] = thumbnail
	return thumbnail, nil
}
