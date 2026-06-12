package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Hello, Bounty Hunter!")

	// 演示使用幂等键发送评论 POST 请求
	// 实际使用时，请替换为真实的 GitHub API 评论 URL 和 Token
	client := &http.Client{Timeout: 10 * time.Second}
	commentBody := map[string]string{"body": "This is a test comment with idempotency key."}
	url := "https://api.github.com/repos/owner/repo/issues/1/comments"

	resp, err := PostComment(client, url, commentBody)
	if err != nil {
		log.Printf("Failed to post comment: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Comment posted, status: %s\n", resp.Status)
}
