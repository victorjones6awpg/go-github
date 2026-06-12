package github

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

type IssueComment struct {
	ID   *int64  `json:"id,omitempty"`
	Body *string `json:"body,omitempty"`
	User *User   `json:"user,omitempty"`
}

func (c *IssueComment) GetID() int64 {
	if c == nil || c.ID == nil {
		return 0
	}
	return *c.ID
}

func (c *IssueComment) GetBody() string {
	if c == nil || c.Body == nil {
		return ""
	}
	return *c.Body
}

func (c *IssueComment) GetUser() *User {
	if c == nil {
		return nil
	}
	return c.User
}

type IssueListCommentsOptions struct {
	Since *time.Time `url:"since,omitempty"`
}

type IssuesService struct {
	client *Client

	CreateCommentFunc func(ctx context.Context, owner, repo string, number int, comment *IssueComment) (*IssueComment, *Response, error)
	ListCommentsFunc  func(ctx context.Context, owner, repo string, number int, opts *IssueListCommentsOptions) ([]*IssueComment, *Response, error)
}

func (s *IssuesService) CreateComment(ctx context.Context, owner, repo string, number int, comment *IssueComment) (*IssueComment, *Response, error) {
	if s.CreateCommentFunc != nil {
		return s.CreateCommentFunc(ctx, owner, repo, number, comment)
	}
	return nil, nil, errors.New("not implemented")
}

func (s *IssuesService) ListComments(ctx context.Context, owner, repo string, number int, opts *IssueListCommentsOptions) ([]*IssueComment, *Response, error) {
	if s.ListCommentsFunc != nil {
		return s.ListCommentsFunc(ctx, owner, repo, number, opts)
	}
	return nil, nil, nil
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
		return true
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if strings.Contains(strings.ToLower(err.Error()), "timeout") {
		return true
	}
	return false
}

func (s *IssuesService) findDuplicateComment(ctx context.Context, owner, repo string, number int, body string, window time.Duration) (*IssueComment, error) {
	authUser, _, err := s.client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	author := authUser.GetLogin()

	since := time.Now().Add(-window)
	opts := &IssueListCommentsOptions{
		Since: &since,
	}

	comments, _, err := s.ListComments(ctx, owner, repo, number, opts)
	if err != nil {
		return nil, err
	}

	for _, c := range comments {
		if c.GetUser().GetLogin() == author && c.GetBody() == body {
			return c, nil
		}
	}
	return nil, nil
}

func (s *IssuesService) CreateCommentSafe(ctx context.Context, owner, repo string, number int, comment *IssueComment) (*IssueComment, *Response, error) {
	enabled, window := getDeduplicateConfig(ctx)

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 && enabled {
			existingComment, err := s.findDuplicateComment(ctx, owner, repo, number, comment.GetBody(), window)
			if err == nil && existingComment != nil {
				resp := &Response{
					Response: &http.Response{
						StatusCode: http.StatusCreated,
						Status:     "201 Created",
					},
				}
				return existingComment, resp, nil
			}
		}

		result, resp, err := s.CreateComment(ctx, owner, repo, number, comment)
		if err == nil {
			return result, resp, nil
		}

		if !isTimeoutError(err) {
			return nil, resp, err
		}

		lastErr = err
	}
	return nil, nil, lastErr
}
