package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type fileResult struct {
	FileUniqueId string `json:"file_unique_id"`
	FilePath     string `json:"file_path"`
}

type getFileResponse struct {
	Ok          bool       `json:"ok"`
	ErrorCode   int        `json:"error_code"`
	Description string     `json:"description"`
	Result      fileResult `json:"result"`
}

type Client struct {
	apiRoot  string
	botToken string
	client   *http.Client
}

func NewClient(apiRoot, botToken string) *Client {
	return &Client{
		apiRoot:  apiRoot,
		botToken: botToken,
		client:   &http.Client{},
	}
}

// GetFile fetches file path from the Telegram API for the given file_id.
// Returns an error if the request fails, the response is invalid, or the API returns an error.
func (c *Client) GetFile(fileId string) (string, error) {
	url := fmt.Sprintf("%s/bot%s/getFile?file_id=%s", c.apiRoot, c.botToken, fileId)
	resp, err := c.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to execute request to telegram: %w", err)
	}
	defer resp.Body.Close()

	var fileInfo getFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return "", fmt.Errorf("failed to decode telegram response: %w", err)
	}

	if !fileInfo.Ok {
		return "", fmt.Errorf("telegram API error (%d): %s", fileInfo.ErrorCode, fileInfo.Description)
	}

	if fileInfo.Result.FilePath == "" {
		return "", errors.New("telegram API returned an empty file path")
	}

	return fileInfo.Result.FilePath, nil
}
