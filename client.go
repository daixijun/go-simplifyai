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
const (
	baseURL            = "https://translate.simplifyai.cn/api/v1"
	defaultContextType = "application/json"
	defaultTimeout     = 10 * time.Second
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	timeout    time.Duration
	baseUrl    string
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{apiKey: apiKey, baseUrl: baseURL, timeout: defaultTimeout}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.timeout,
		}
	}

	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, contentType string) ([]byte, error) {
	if len(contentType) == 0 {
		contentType = defaultContextType
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseUrl+path, body)
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
func (c *Client) CreateTranslationTask(ctx context.Context, req *CreateTranslationTaskRequest) (*CreateTranslationTaskResponse, error) {
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
	data, err := c.doRequest(ctx, http.MethodPost, "/translations", bodyBuffer, bodyWriter.FormDataContentType())
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
func (c *Client) QueryTranslationTask(ctx context.Context, taskId string) (*QueryTranslationTaskResponse, error) {
	data, err := c.doRequest(ctx, http.MethodGet, "/translations/"+taskId, nil, "")
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
func (c *Client) DeleteTranslationTask(ctx context.Context, taskId string) error {
	_, err := c.doRequest(ctx, http.MethodDelete, "/translations/"+taskId, nil, "")
	return err
}

// StartTranslationTask 启动翻译任务
func (c *Client) StartTranslationTask(ctx context.Context, taskId string) (*StartTranslationTaskResponse, error) {
	data, err := c.doRequest(ctx, http.MethodPut, "/translations/"+taskId, nil, "")
	if err != nil {
		return nil, err
	}

	var task StartTranslationTaskResponse
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// ListAvailableLanguages 列出支持翻译的语言
func (c *Client) ListAvailableLanguages(ctx context.Context) ([]string, error) {
	data, err := c.doRequest(ctx, http.MethodGet, "/languages", nil, "")
	if err != nil {
		return nil, err
	}

	var languages []string
	if err := json.Unmarshal(data, &languages); err != nil {
		return nil, err
	}

	return languages, nil
}
