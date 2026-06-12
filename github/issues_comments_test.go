package github

import (
	"context"
	"errors"
	"testing"
)

type mockTimeoutError struct {
	error
}

func (e mockTimeoutError) Timeout() bool   { return true }
func (e mockTimeoutError) Temporary() bool { return true }

func TestIssuesService_CreateCommentSafe_Success(t *testing.T) {
	client := NewClient(nil)
	commentBody := "Hello world"
	comment := &IssueComment{Body: &commentBody}

	client.Issues.CreateCommentFunc = func(ctx context.Context, owner, repo string, number int, c *IssueComment) (*IssueComment, *Response, error) {
		id := int64(123)
		return &IssueComment{ID: &id, Body: c.Body}, nil, nil
	}

	res, _, err := client.Issues.CreateCommentSafe(context.Background(), "owner", "repo", 1, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GetID() != 123 {
		t.Errorf("expected ID 123, got %d", res.GetID())
	}
}

func TestIssuesService_CreateCommentSafe_TimeoutAndDeduplicate(t *testing.T) {
	client := NewClient(nil)
	commentBody := "Hello world"
	comment := &IssueComment{Body: &commentBody}

	calls := 0
	client.Issues.CreateCommentFunc = func(ctx context.Context, owner, repo string, number int, c *IssueComment) (*IssueComment, *Response, error) {
		calls++
		if calls == 1 {
			return nil, nil, mockTimeoutError{errors.New("request timed out")}
		}
		t.Errorf("CreateComment should not be called a second time")
		return nil, nil, errors.New("should not be called")
	}

	client.Issues.ListCommentsFunc = func(ctx context.Context, owner, repo string, number int, opts *IssueListCommentsOptions) ([]*IssueComment, *Response, error) {
		id := int64(123)
		login := "test-user"
		return []*IssueComment{
			{
				ID:   &id,
				Body: &commentBody,
				User: &User{Login: &login},
			},
		}, nil, nil
	}

	res, resp, err := client.Issues.CreateCommentSafe(context.Background(), "owner", "repo", 1, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GetID() != 123 {
		t.Errorf("expected ID 123, got %d", res.GetID())
	}
	if resp == nil || resp.StatusCode != 201 {
		t.Errorf("expected simulated 201 response, got %v", resp)
	}
	if calls != 1 {
		t.Errorf("expected exactly 1 call to CreateComment, got %d", calls)
	}
}

func TestIssuesService_CreateCommentSafe_TimeoutNoDuplicate(t *testing.T) {
	client := NewClient(nil)
	commentBody := "Hello world"
	comment := &IssueComment{Body: &commentBody}

	calls := 0
	client.Issues.CreateCommentFunc = func(ctx context.Context, owner, repo string, number int, c *IssueComment) (*IssueComment, *Response, error) {
		calls++
		if calls == 1 {
			return nil, nil, mockTimeoutError{errors.New("request timed out")}
		}
		id := int64(124)
		return &IssueComment{ID: &id, Body: c.Body}, nil, nil
	}

	client.Issues.ListCommentsFunc = func(ctx context.Context, owner, repo string, number int, opts *IssueListCommentsOptions) ([]*IssueComment, *Response, error) {
		return nil, nil, nil
	}

	res, _, err := client.Issues.CreateCommentSafe(context.Background(), "owner", "repo", 1, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GetID() != 124 {
		t.Errorf("expected ID 124, got %d", res.GetID())
	}
	if calls != 2 {
		t.Errorf("expected 2 calls to CreateComment, got %d", calls)
	}
}
