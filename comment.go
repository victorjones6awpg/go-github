package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// GenerateIdempotencyKey 生成一个随机的幂等键，用于 POST 请求防止重复提交。
// 返回 16 字节的十六进制字符串。
func GenerateIdempotencyKey() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate idempotency key: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// setIdempotencyHeader 向 HTTP 请求中添加幂等键头。
func setIdempotencyHeader(req *http.Request) error {
	key, err := GenerateIdempotencyKey()
	if err != nil {
		return err
	}
	req.Header.Set("X-GitHub-Idempotency-Key", key)
	return nil
}

// PostComment 向指定 URL 发送 POST 请求，自动生成并添加幂等键。
// body 将被序列化为 JSON。返回 HTTP 响应和可能的错误。
func PostComment(client *http.Client, url string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 添加幂等键头
	if err := setIdempotencyHeader(req); err != nil {
		return nil, fmt.Errorf("failed to set idempotency key: %w", err)
	}

	return client.Do(req)
}
