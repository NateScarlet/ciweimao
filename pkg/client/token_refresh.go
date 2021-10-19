package client

import (
	"context"
)

// TokenRefresher refresh login token on demand.
// use custom logic to store refreshed Client.LoginToken
type TokenRefresher interface {
	RefreshToken(ctx context.Context, c *Client) (err error)
}

type TokenRefreshFunc func(ctx context.Context, c *Client) (err error)

func (fn TokenRefreshFunc) RefreshToken(ctx context.Context, c *Client) error {
	return fn(ctx, c)

}

func NewLoginTokenRefresher(username, password string) TokenRefresher {
	return TokenRefreshFunc(func(ctx context.Context, c *Client) (err error) {
		_, err = c.Login(ctx, username, password)
		if err != nil {
			return
		}
		return
	})
}

type contextKeySkipTokenRefresh struct{}

func WithSkipTokenRefresh(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, contextKeySkipTokenRefresh{}, v)
}

func SkipTokenRefresh(ctx context.Context) bool {
	var v, _ = ctx.Value(contextKeySkipTokenRefresh{}).(bool)
	return v
}
