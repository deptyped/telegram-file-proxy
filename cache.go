package main

import (
	"errors"
	"time"

	lru "github.com/hnlq715/golang-lru"
)

// Since a file can have different valid file_ids, we need to cache the file path using file_unique_id.

// (file_id, file_unique_id)
var fileUniqueIdCache, _ = lru.NewARCWithExpire(100_000, 3*time.Hour)

// (file_unique_id, file_path)
var filePathCache, _ = lru.NewWithExpire(
	100_000,
	59*time.Minute, // It is guaranteed that the link will be valid for at least 1 hour.
)

func cacheFilePath(fileId, fileUniqueId, filePath string) {
	fileUniqueIdCache.Add(fileId, fileUniqueId)
	filePathCache.Add(fileUniqueId, filePath)
}

func getFilePath(fileId string) (string, error) {
	fileUniqueIdVal, _ := fileUniqueIdCache.Get(fileId)
	if fileUniqueIdVal != nil {
		fileUniqueId := fileUniqueIdVal.(string)
		filePathVal, _ := filePathCache.Get(fileUniqueId)

		if filePathVal != nil {
			filePath := filePathVal.(string)

			if len(filePath) != 0 {
				return filePath, nil
			}
		}
	}

	return "", errors.New("no cached value")
}
