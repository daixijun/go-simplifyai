package simplifyai

import (
	"net/http"
	"time"
)

type Option func(opts *Client)

// WithHttpClient 自定义 http.Client
func WithHttpClient(httpClient *http.Client) Option {
	return func(opts *Client) {
		opts.httpClient = httpClient
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout *time.Duration) Option {
	return func(opts *Client) {
		opts.httpClient.Timeout = *timeout
	}
}

// WithBaseUrl 设置请求地址
func WithBaseUrl(baseUrl string) Option {
	return func(opts *Client) {
		opts.baseUrl = baseUrl
	}
}
