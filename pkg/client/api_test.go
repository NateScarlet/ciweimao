package client

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/NateScarlet/snapshot/pkg/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecryptAPIResponse(t *testing.T) {

	var client = Default
	var resp, err = client.PostForm(
		client.EndpointURL("/setting/get_version", nil).String(),
		url.Values{
			"app_version": []string{"2.6.011"},
		},
	)
	require.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Log(string(respData))

	res, err := client.DecryptAPIResponse(bytes.NewBuffer(respData))
	require.NoError(t, err)
	resText, err := ioutil.ReadAll(res)
	require.NoError(t, err)
	t.Log(string(resText))

	var doc = new(interface{})
	err = json.Unmarshal(resText, doc)
	require.NoError(t, err)
	snapshot.MatchJSON(t, *doc, snapshot.OptionTransform(snapshot.TransformSchema))
}

// encryptAPIResponse 使用客户端的 AES-CBC 密钥加密 JSON 响应，返回 base64 编码的密文。
func encryptAPIResponse(c *Client, plaintext string) string {
	// PKCS7 填充
	blockSize := 16
	data := []byte(plaintext)
	padLen := blockSize - len(data)%blockSize
	for i := 0; i < padLen; i++ {
		data = append(data, byte(padLen))
	}

	block, err := aes.NewCipher(c.apiAESKey())
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, c.APIInitialVector)
	mode.CryptBlocks(ciphertext, data)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

// TestCallShouldRefreshTokenOn200001 验证当 API 返回错误码 200001 且 LoginToken 为空时，
// 客户端会调用 TokenRefresher 刷新令牌并重试请求。
func TestCallShouldRefreshTokenOn200001(t *testing.T) {
	var client = new(Client)
	// 使用确定性密钥，让加密/解密一致
	client.APIInitialVector = make([]byte, 16)
	client.APIKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"
	client.LoginToken = "" // 模拟无 LOGIN_TOKEN 场景

	var callCount int32

	// 构造两个响应：
	// 第一次：200001（缺少登录必需参数）
	// 第二次：100000（成功）
	resp200001 := encryptAPIResponse(client, `{"code":"200001","tip":"缺少登录必需参数"}`)
	resp100000 := encryptAPIResponse(client, `{"code":"100000","data":{}}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		w.WriteHeader(http.StatusOK)
		if n == 1 {
			w.Write([]byte(resp200001))
		} else {
			w.Write([]byte(resp100000))
		}
	}))
	defer server.Close()

	client.ServerURL = server.URL

	// 记录 TokenRefresher 是否被调用
	var refreshCalled bool
	client.TokenRefresher = TokenRefreshFunc(func(ctx context.Context, c *Client) error {
		refreshCalled = true
		// 模拟刷新成功：设置 LoginToken
		c.LoginToken = "new_token_from_refresh"
		return nil
	})

	// Act
	ctx := context.Background()
	_, err := client.Call(ctx, "/bookcity/get_rank_book_list", nil)

	// Assert: 不应有错误（重试后成功）
	require.NoError(t, err)
	// TokenRefresher 应被调用一次
	assert.True(t, refreshCalled, "TokenRefresher should have been called for 200001 error")
	// 总共应有 2 次 HTTP 调用（第一次 200001，第二次重试成功）
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount), "should retry after token refresh")
}

// TestCallShouldNotRefreshTokenOn200001IfLoginTokenExists 验证当 LoginToken 已存在时，
// 即使收到 200001 错误也不触发刷新（因为这不是 token 缺失问题）。
func TestCallShouldNotRefreshTokenOn200001IfLoginTokenExists(t *testing.T) {
	var client = new(Client)
	client.APIInitialVector = make([]byte, 16)
	client.APIKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"
	client.LoginToken = "existing_token" // 已有 token

	resp200001 := encryptAPIResponse(client, `{"code":"200001","tip":"缺少登录必需参数"}`)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp200001))
	}))
	defer server.Close()

	client.ServerURL = server.URL

	var refreshCalled bool
	client.TokenRefresher = TokenRefreshFunc(func(ctx context.Context, c *Client) error {
		refreshCalled = true
		return nil
	})

	// Act
	ctx := context.Background()
	_, err := client.Call(ctx, "/bookcity/get_rank_book_list", nil)

	// Assert: 应返回错误（不应重试）
	require.Error(t, err)
	assert.Contains(t, err.Error(), "200001")
	assert.False(t, refreshCalled, "TokenRefresher should NOT be called when LoginToken is already set")
}