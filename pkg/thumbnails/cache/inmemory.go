package cache

import "image"

func NewInMemoryCache() InMemoryCache {
	return InMemoryCache{
		store: make(map[string]image.Image),
	}
}

type InMemoryCache struct {
	store map[string]image.Image
}

func (fsc InMemoryCache) Get(key string) image.Image {
	return fsc.store[key]
}

func (fsc InMemoryCache) Set(key string, thumbnail image.Image) (image.Image, error) {
	fsc.store[key] = thumbnail
	return thumbnail, nil
}
