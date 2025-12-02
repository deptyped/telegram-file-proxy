package main

import (
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Server struct {
	config   *Config
	cache    *Cache
	telegram *Client
	proxy    *httputil.ReverseProxy
}

func NewServer(config *Config, cache *Cache, telegram *Client) (*Server, error) {
	remote, err := url.Parse(config.ApiRoot)
	if err != nil {
		return nil, fmt.Errorf("invalid API root URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	s := &Server{
		config:   config,
		cache:    cache,
		telegram: telegram,
		proxy:    proxy,
	}

	proxy.ModifyResponse = s.modifyProxyResponse
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Error contacting the upstream server", http.StatusBadGateway)
	}

	return s, nil
}

// resolveFilePath resolves a file_id to a file_path, using cache when available.
// On any error, the cache is automatically invalidated for this file_id.
func (s *Server) resolveFilePath(fileId string) (string, error) {
	// Try cache first
	if filePath, err := s.cache.getFilePath(fileId); err == nil {
		return filePath, nil
	}

	// Not in cache, fetch from Telegram API
	filePath, err := s.telegram.GetFile(fileId)
	if err != nil {
		// Invalidate cache on any error (corrupted response, API error, empty path, etc.)
		s.cache.invalidate(fileId)
		return "", err
	}

	// Cache the new path and return it
	s.cache.cacheFilePath(fileId, filePath)
	return filePath, nil
}

var allowedMIMETypes = map[string]struct{}{
	"image/jpeg":               {},
	"image/gif":                {},
	"image/png":                {},
	"image/webp":               {},
	"application/x-tgsticker":  {},
	"video/webm":               {},
	"application/pdf":          {},
	"application/zip":          {},
	"audio/mpeg":               {},
	"audio/ogg":                {},
	"audio/x-matroska":         {},
	"audio/opus":               {},
	"video/quicktime":          {},
	"video/mp4":                {},
	"video/x-matroska":         {},
	"application/octet-stream": {},
}

func (s *Server) handleFileRequest(w http.ResponseWriter, r *http.Request) {
	fileId := strings.TrimPrefix(r.URL.Path, "/")
	if fileId == "" {
		http.Error(w, "Please provide a file_id in the URL path (e.g., /<file_id>)", http.StatusBadRequest)
		return
	}

	if fileId == "favicon.ico" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	filePath, err := s.resolveFilePath(fileId)
	if err != nil {
		log.Printf("Error resolving file path for file_id %s: %v", fileId, err)
		http.Error(w, "Could not resolve file_id.", http.StatusNotFound)
		return
	}

	contentType := r.URL.Query().Get("content-type")
	fileName := r.URL.Query().Get("filename")

	if s.config.IsApiLocal {
		headers := w.Header()
		s.modifyHeaders(&headers, contentType, fileName)
		http.ServeFile(w, r, filePath)
	} else {
		// Rewrite the request to be proxied
		r.URL.Path = fmt.Sprintf("/file/bot%s/%s", s.config.BotToken, filePath)
		r.RequestURI = ""
		r.Host = r.URL.Host
		r.Header = make(http.Header)

		proxy := *s.proxy
		proxy.ModifyResponse = func(resp *http.Response) error {
			if resp.StatusCode >= 400 {
				log.Printf("Upstream returned status %d for file_id %s. Invalidating cache.", resp.StatusCode, fileId)
				s.cache.invalidate(fileId)
			}
			return s.modifyProxyResponse(resp)
		}
		proxy.ServeHTTP(w, r)
	}
}

func (s *Server) modifyHeaders(h *http.Header, contentType, fileName string) {
	h.Del("Server")
	if _, ok := allowedMIMETypes[contentType]; ok {
		h.Set("Content-Type", contentType)
	} else {
		h.Del("Content-Type")
	}

	disposition := "inline"
	if fileName != "" {
		disposition = mime.FormatMediaType("inline", map[string]string{"filename": fileName})
	}
	h.Set("Content-Disposition", disposition)
	h.Set("Cache-Control", "public, max-age=31536000") // 1 year
}

func (s *Server) modifyProxyResponse(resp *http.Response) error {
	contentType := resp.Request.URL.Query().Get("content-type")
	fileName := resp.Request.URL.Query().Get("filename")
	s.modifyHeaders(&resp.Header, contentType, fileName)
	return nil
}
