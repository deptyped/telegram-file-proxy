package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	router "github.com/julienschmidt/httprouter"
)

type File struct {
	FileUniqueId string `json:"file_unique_id"`
	FilePath     string `json:"file_path"`
}

type GetFileResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	Result      File   `json:"result"`
}

func fetchFile(apiRoot, botToken, fileId string) (GetFileResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/bot%s/getFile?file_id=%s", apiRoot, botToken, fileId))
	if err != nil {
		return GetFileResponse{}, err
	}
	defer resp.Body.Close()

	var fileInfo GetFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return GetFileResponse{}, err
	}

	return fileInfo, nil
}

func modifyHeaders(headers *http.Header) {
	headers.Del("Server")
	headers.Del("Content-Type")                              // Remove default content type
	headers.Set("Content-Disposition", "inline")             // Display media inline
	headers.Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
}

func ServeFile(config *Config, cache *Cache) router.Handle {
	remote, err := url.Parse(config.ApiRoot)
	if err != nil {
		log.Fatalf("Invalid API root URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = func(resp *http.Response) error {
		modifyHeaders(&resp.Header)
		return nil
	}

	return router.Handle(func(res http.ResponseWriter, req *http.Request, params router.Params) {
		fileId := params.ByName("fileId")

		filePath, err := cache.getFilePath(fileId)
		if err != nil { // Cache miss, fetch from API
			fileInfo, err := fetchFile(config.ApiRoot, config.BotToken, fileId)
			if err != nil {
				log.Printf("Error fetching file: %v", err)
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			cache.cacheFilePath(fileId, fileInfo.Result.FileUniqueId, fileInfo.Result.FilePath)
			filePath = fileInfo.Result.FilePath
		}

		if config.IsApiLocal {
			headers := res.Header()
			modifyHeaders(&headers)
			http.ServeFile(res, req, filePath)
		} else {
			req.URL, _ = url.Parse(fmt.Sprintf("%s/file/bot%s/%s", config.ApiRoot, config.BotToken, filePath))
			req.Host = req.URL.Host
			proxy.ServeHTTP(res, req)
		}
	})
}

func main() {
	config := newConfig()
	cache := newCache()

	router := router.New()
	router.GET("/:fileId", ServeFile(config, cache))

	log.Printf("Server is running at %s\n", config.ServerAddr)
	log.Fatal(http.ListenAndServe(config.ServerAddr, router))
}
