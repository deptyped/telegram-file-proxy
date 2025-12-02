package main

import (
	"errors"
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
)

type Cache struct {
	filePathCache *ttlcache.Cache[string, string]
}

func newCache() *Cache {
	c := &Cache{
		filePathCache: ttlcache.New(
			ttlcache.WithTTL[string, string](59*time.Minute),
			ttlcache.WithCapacity[string, string](100_000),
		),
	}

	// Start a goroutine to clean up expired items
	go c.filePathCache.Start()

	return c
}

func (c *Cache) cacheFilePath(fileId, filePath string) {
	c.filePathCache.Set(fileId, filePath, ttlcache.DefaultTTL)
}

func (c *Cache) getFilePath(fileId string) (string, error) {
	item := c.filePathCache.Get(fileId)
	if item == nil || item.Value() == "" {
		return "", errors.New("file path not found in cache")
	}
	return item.Value(), nil
}

func (c *Cache) invalidate(fileId string) {
	c.filePathCache.Delete(fileId)
}
