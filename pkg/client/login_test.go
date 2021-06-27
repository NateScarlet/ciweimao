package client

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/NateScarlet/snapshot/pkg/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	var ctx = context.Background()
	var username = os.Getenv("TEST_CIWEIMAO_USERNAME")
	var password = os.Getenv("TEST_CIWEIMAO_PASSWORD")
	if username == "" {
		t.Skip("test user not configured")
	}

	var c = new(Client)
	c.ApplyDefaultConfig()
	res, err := c.Login(ctx, username, password)
	require.NoError(t, err)
	if snapshot.DefaultUpdate {
		var doc = new(interface{})
		err = json.Unmarshal([]byte(res.Raw), doc)
		require.NoError(t, err)
		snapshot.MatchJSON(t, *doc, snapshot.OptionExt(".local.response.json"))
	}
	assert.NotEmpty(t, c.Account)
	assert.NotEmpty(t, c.LoginToken)
}
