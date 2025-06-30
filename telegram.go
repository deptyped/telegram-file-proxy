package main

import (
	"encoding/json"
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

// GetFile fetches file metadata from the Telegram API.
func (c *Client) GetFile(fileId string) (*getFileResponse, error) {
	url := fmt.Sprintf("%s/bot%s/getFile?file_id=%s", c.apiRoot, c.botToken, fileId)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to telegram: %w", err)
	}
	defer resp.Body.Close()

	var fileInfo getFileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return nil, fmt.Errorf("failed to decode telegram response: %w", err)
	}

	return &fileInfo, nil
}
