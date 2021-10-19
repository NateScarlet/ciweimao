package client

import (
	"context"
	"net/url"

	"github.com/tidwall/gjson"
)

func (c *Client) Login(ctx context.Context, username, password string) (ret gjson.Result, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var data = url.Values{}
	c.SetDefaultAPIAuthData(data)
	data.Set("login_name", username)
	data.Set("passwd", password)
	data.Del("account")
	data.Del("login_token")

	ret, err = c.Call(ctx, "/signup/login", data)
	if err != nil {
		return
	}
	c.LoginToken = ret.Get("data.login_token").String()
	c.Account = ret.Get("data.reader_info.account").String()
	return
}
