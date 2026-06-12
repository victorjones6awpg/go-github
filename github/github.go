package github

import (
	"context"
	"net/http"
	"time"
)

type Client struct {
	Issues       *IssuesService
	PullRequests *PullRequestsService
	Users        *UsersService
}

func NewClient(httpClient *http.Client) *Client {
	c := &Client{}
	c.Issues = &IssuesService{client: c}
	c.PullRequests = &PullRequestsService{client: c}
	c.Users = &UsersService{client: c}
	return c
}

type Response struct {
	*http.Response
}

type User struct {
	Login *string `json:"login,omitempty"`
}

func (u *User) GetLogin() string {
	if u == nil || u.Login == nil {
		return ""
	}
	return *u.Login
}

type UsersService struct {
	client *Client
}

func (s *UsersService) Get(ctx context.Context, user string) (*User, *Response, error) {
	login := "test-user"
	return &User{Login: &login}, &Response{Response: &http.Response{StatusCode: 200}}, nil
}

type DeduplicateOptions struct {
	Enabled bool
	Window  time.Duration
}

type deduplicateKey struct{}

func WithDeduplication(ctx context.Context, opts DeduplicateOptions) context.Context {
	return context.WithValue(ctx, deduplicateKey{}, opts)
}

func getDeduplicateConfig(ctx context.Context) (bool, time.Duration) {
	if opts, ok := ctx.Value(deduplicateKey{}).(DeduplicateOptions); ok {
		return opts.Enabled, opts.Window
	}
	return true, 2 * time.Minute
}
