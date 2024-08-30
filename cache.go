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
	filePathVal := c.filePathCache.Get(fileId)

	if filePathVal != nil {
		filePath := filePathVal.Value()
		if len(filePath) != 0 {
			return filePath, nil
		}
	}

	return "", errors.New("no cached value")
}
