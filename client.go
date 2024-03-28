package simplifyai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// DOC: https://translate.simplifyai.cn/developer
const baseURL = "https://translate.simplifyai.cn/api/v1"
const defaultContextType = "application/json"

type client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string, httpClient *http.Client) *client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &client{httpClient: httpClient, apiKey: apiKey}
}

func (c *client) request(ctx context.Context, method, path string, body io.Reader, contentType string) ([]byte, error) {
	if len(contentType) == 0 {
		contentType = defaultContextType
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, ErrorMap[resp.StatusCode]
	}

	return io.ReadAll(resp.Body)
}

// CreateTranslationTask 创建翻译任务
func (c *client) CreateTranslationTask(ctx context.Context, req *CreateTranslationTaskRequest) (*CreateTranslationTaskResponse, error) {
	file, err := os.Open(req.File)
	if err != nil {
		return nil, fmt.Errorf("open file failed: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("file", req.File)

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return nil, fmt.Errorf("copy file to bodyWriter failed: %w", err)
	}
	_ = bodyWriter.WriteField("fromLang", req.FromLang)
	_ = bodyWriter.WriteField("toLang", req.ToLang)
	if len(req.Glossary) > 0 {
		_ = bodyWriter.WriteField("glossary", req.Glossary)
	}
	_ = bodyWriter.Close()
	data, err := c.request(ctx, "POST", "/translations", bodyBuffer, bodyWriter.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var task CreateTranslationTaskResponse
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}

	return &task, nil

}

// QueryTranslationTask 查询翻译任务
func (c *client) QueryTranslationTask(ctx context.Context, taskID string) (*QueryTranslationTaskResponse, error) {
	data, err := c.request(ctx, "GET", "/translations/"+taskID, nil, "")
	if err != nil {
		return nil, err
	}

	var task QueryTranslationTaskResponse
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// DeleteTranslationTask 删除翻译任务
func (c *client) DeleteTranslationTask(ctx context.Context, taskID string) error {
	_, err := c.request(ctx, "DELETE", "/translations/"+taskID, nil, "")
	return err
}

// StartTranslationTask 启动翻译任务
func (c *client) StartTranslationTask(ctx context.Context, taskID string) (*QueryTranslationTaskResponse, error) {
	data, err := c.request(ctx, "PUT", "/translations/"+taskID, nil, "")
	if err != nil {
		return nil, err
	}

	var task QueryTranslationTaskResponse
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// ListAvailableLanguages 列出支持翻译的语言
func (c *client) ListAvailableLanguages(ctx context.Context) ([]string, error) {
	data, err := c.request(ctx, "GET", "/languages", nil, "")
	if err != nil {
		return nil, err
	}

	var languages []string
	if err := json.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	return languages, nil
}
