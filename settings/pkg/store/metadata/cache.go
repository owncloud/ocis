package store

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
)

var (
	cachettl = 0
	// these need to be global instances for now as the `Service` (and therefore the `Store`) are instantiated twice (for grpc and http)
	// therefore caches need to cover both instances
	dircache   = initCache(cachettl)
	filescache = initCache(cachettl)
)

// CachedMDC is cache for the metadataclient
type CachedMDC struct {
	next MetadataClient

	files *ttlcache.Cache
	dirs  *ttlcache.Cache
}

// SimpleDownload caches the answer from SimpleDownload or returns the cached one
func (c *CachedMDC) SimpleDownload(ctx context.Context, id string) ([]byte, error) {
	if b, err := c.files.Get(id); err == nil {
		return b.([]byte), nil
	}
	b, err := c.next.SimpleDownload(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = c.files.Set(id, b)
	return b, nil
}

// SimpleUpload caches the answer from SimpleUpload and invalidates the cache
func (c *CachedMDC) SimpleUpload(ctx context.Context, id string, content []byte) error {
	b, err := c.files.Get(id)
	if err == nil && string(b.([]byte)) == string(content) {
		// no need to bug mdc
		return nil
	}

	err = c.next.SimpleUpload(ctx, id, content)
	if err != nil {
		return err
	}

	// invalidate caches
	_ = c.dirs.Remove(path.Dir(id))
	_ = c.files.Set(id, content)
	return nil
}

// Delete invalidates the cache when operation was successful
func (c *CachedMDC) Delete(ctx context.Context, id string) error {
	if err := c.next.Delete(ctx, id); err != nil {
		return err
	}

	// invalidate caches
	_ = removePrefix(c.files, id)
	_ = removePrefix(c.dirs, id)
	return nil
}

// ReadDir caches the response from ReadDir or returnes the cached one
func (c *CachedMDC) ReadDir(ctx context.Context, id string) ([]string, error) {
	i, err := c.dirs.Get(id)
	if err == nil {
		return i.([]string), nil
	}

	fmt.Println("readdir calling metadataservice", id)
	s, err := c.next.ReadDir(ctx, id)
	fmt.Println("readdir calling metadataservice result", s, err)
	if err != nil {
		return nil, err
	}

	return s, c.dirs.Set(id, s)
}

// MakeDirIfNotExist invalidates the cache
func (c *CachedMDC) MakeDirIfNotExist(ctx context.Context, id string) error {
	err := c.next.MakeDirIfNotExist(ctx, id)
	if err != nil {
		return err
	}

	// invalidate caches
	_ = c.dirs.Remove(path.Dir(id))
	return nil
}

// Init instantiates the caches
func (c *CachedMDC) Init(ctx context.Context, id string) error {
	c.dirs = dircache
	c.files = filescache
	return c.next.Init(ctx, id)
}

func initCache(ttlSeconds int) *ttlcache.Cache {
	cache := ttlcache.NewCache()
	_ = cache.SetTTL(time.Duration(ttlSeconds) * time.Second)
	cache.SkipTTLExtensionOnHit(true)
	return cache
}

func removePrefix(cache *ttlcache.Cache, prefix string) error {
	for _, k := range cache.GetKeys() {
		if strings.HasPrefix(k, prefix) {
			if err := cache.Remove(k); err != nil {
				return err
			}
		}
	}
	return nil
}
