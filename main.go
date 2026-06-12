package main

import (
	"context"
	"fmt"
	"time"

	"github.com/victorjones6awpg/go-github/github"
)

func main() {
	fmt.Println("Demonstrating Deduplication Helper...")
	client := github.NewClient(nil)

	commentBody := "Test comment"
	client.Issues.CreateCommentFunc = func(ctx context.Context, owner, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error) {
		fmt.Println("CreateComment called!")
		return nil, nil, fmt.Errorf("request timed out (mock timeout)")
	}

	client.Issues.ListCommentsFunc = func(ctx context.Context, owner, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, *github.Response, error) {
		fmt.Println("ListComments called to check for duplicates!")
		id := int64(1001)
		login := "test-user"
		return []*github.IssueComment{
			{
				ID:   &id,
				Body: &commentBody,
				User: &github.User{Login: &login},
			},
		}, nil, nil
	}

	ctx := github.WithDeduplication(context.Background(), github.DeduplicateOptions{
		Enabled: true,
		Window:  2 * time.Minute,
	})

	comment := &github.IssueComment{Body: &commentBody}
	res, resp, err := client.Issues.CreateCommentSafe(ctx, "owner", "repo", 1, comment)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Successfully handled timeout! Returned Comment ID: %d, Status Code: %d\n", res.GetID(), resp.StatusCode)
}
