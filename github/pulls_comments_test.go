package github

import (
	"context"
	"errors"
	"testing"
)

func TestPullRequestsService_CreateCommentSafe_TimeoutAndDeduplicate(t *testing.T) {
	client := NewClient(nil)
	commentBody := "Hello PR"
	comment := &PullRequestComment{Body: &commentBody}

	calls := 0
	client.PullRequests.CreateCommentFunc = func(ctx context.Context, owner, repo string, number int, c *PullRequestComment) (*PullRequestComment, *Response, error) {
		calls++
		if calls == 1 {
			return nil, nil, mockTimeoutError{errors.New("request timed out")}
		}
		t.Errorf("CreateComment should not be called a second time")
		return nil, nil, errors.New("should not be called")
	}

	client.PullRequests.ListCommentsFunc = func(ctx context.Context, owner, repo string, number int, opts *PullRequestListCommentsOptions) ([]*PullRequestComment, *Response, error) {
		id := int64(456)
		login := "test-user"
		return []*PullRequestComment{
			{
				ID:   &id,
				Body: &commentBody,
				User: &User{Login: &login},
			},
		}, nil, nil
	}

	res, resp, err := client.PullRequests.CreateCommentSafe(context.Background(), "owner", "repo", 1, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GetID() != 456 {
		t.Errorf("expected ID 456, got %d", res.GetID())
	}
	if resp == nil || resp.StatusCode != 201 {
		t.Errorf("expected simulated 201 response, got %v", resp)
	}
	if calls != 1 {
		t.Errorf("expected exactly 1 call to CreateComment, got %d", calls)
	}
}

func TestPullRequestsService_CreateReviewCommentSafe_TimeoutAndDeduplicate(t *testing.T) {
	client := NewClient(nil)
	commentBody := "Hello Review"
	comment := &PullRequestComment{Body: &commentBody}

	calls := 0
	client.PullRequests.CreateReviewCommentFunc = func(ctx context.Context, owner, repo string, number int, c *PullRequestComment) (*PullRequestComment, *Response, error) {
		calls++
		if calls == 1 {
			return nil, nil, mockTimeoutError{errors.New("request timed out")}
		}
		t.Errorf("CreateReviewComment should not be called a second time")
		return nil, nil, errors.New("should not be called")
	}

	client.PullRequests.ListCommentsFunc = func(ctx context.Context, owner, repo string, number int, opts *PullRequestListCommentsOptions) ([]*PullRequestComment, *Response, error) {
		id := int64(789)
		login := "test-user"
		return []*PullRequestComment{
			{
				ID:   &id,
				Body: &commentBody,
				User: &User{Login: &login},
			},
		}, nil, nil
	}

	res, resp, err := client.PullRequests.CreateReviewCommentSafe(context.Background(), "owner", "repo", 1, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GetID() != 789 {
		t.Errorf("expected ID 789, got %d", res.GetID())
	}
	if resp == nil || resp.StatusCode != 201 {
		t.Errorf("expected simulated 201 response, got %v", resp)
	}
	if calls != 1 {
		t.Errorf("expected exactly 1 call to CreateReviewComment, got %d", calls)
	}
}
