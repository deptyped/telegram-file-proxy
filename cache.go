package main

import (
	"errors"
	"time"

	ttlcache "github.com/jellydator/ttlcache/v3"
)

type Cache struct {
	fileUniqueIdCache *ttlcache.Cache[string, string]
	filePathCache     *ttlcache.Cache[string, string]
}

func newCache() *Cache {
	c := &Cache{
		fileUniqueIdCache: ttlcache.New(
			ttlcache.WithTTL[string, string](24*time.Hour),
			ttlcache.WithCapacity[string, string](100_000),
		),
		filePathCache: ttlcache.New(
			ttlcache.WithTTL[string, string](59*time.Minute),
			ttlcache.WithCapacity[string, string](100_000),
		),
	}

	// Start goroutines to clean up expired items
	go c.fileUniqueIdCache.Start()
	go c.filePathCache.Start()

	return c
}

func (c *Cache) cacheFilePath(fileId, fileUniqueId, filePath string) {
	c.fileUniqueIdCache.Set(fileId, fileUniqueId, ttlcache.DefaultTTL)
	c.filePathCache.Set(fileUniqueId, filePath, ttlcache.DefaultTTL)
}

func (c *Cache) getFilePath(fileId string) (string, error) {
	fileUniqueIdVal := c.fileUniqueIdCache.Get(fileId)
	if fileUniqueIdVal != nil {
		fileUniqueId := fileUniqueIdVal.Value()
		filePathVal := c.filePathCache.Get(fileUniqueId)

		if filePathVal != nil {
			filePath := filePathVal.Value()

			if len(filePath) != 0 {
				return filePath, nil
			}
		}
	}

	return "", errors.New("no cached value")
}
