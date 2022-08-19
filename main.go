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

func fetchFile(apiRoot string, botToken string, fileId string) (GetFileResponse, error) {
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

func modifyHeaders(h *http.Header) {
	h.Del("Server")
	// remove application/octet-stream mime type
	h.Del("Content-Type")
	// display media instead of downloading
	h.Set("Content-Disposition", "inline")
	// cache media for 1 year
	h.Set("Cache-Control", "public, max-age=31536000")
}

func ServeFile(config *Config) router.Handle {
	remote, err := url.Parse(config.ApiRoot)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = func(r *http.Response) error {
		modifyHeaders(&r.Header)

		return nil
	}

	return router.Handle(func(res http.ResponseWriter, req *http.Request, params router.Params) {
		fileId := params.ByName("fileId")

		filePath, err := getFilePath(fileId)
		if err != nil {
			fileInfo, err := fetchFile(config.ApiRoot, config.BotToken, fileId)
			if err != nil {
				log.Println(err)
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !fileInfo.Ok {
				res.WriteHeader(fileInfo.ErrorCode)
				return
			}

			cacheFilePath(fileId, fileInfo.Result.FileUniqueId, fileInfo.Result.FilePath)

			filePath = fileInfo.Result.FilePath
		}

		if config.IsApiLocal {
			headers := res.Header()
			modifyHeaders(&headers)

			http.ServeFile(res, req, filePath)
		} else {
			url, _ := url.Parse(fmt.Sprintf("%s/file/bot%s/%s", config.ApiRoot, config.BotToken, filePath))

			req.Host = url.Host
			req.URL.Path = url.Path

			proxy.ServeHTTP(res, req)
		}

	})
}

func main() {
	config := newConfig()

	router := router.New()
	router.GET("/:fileId", ServeFile(config))

	log.Printf("Server is running at %s\n", config.ServerAddr)
	log.Fatal(http.ListenAndServe(config.ServerAddr, router))
}
