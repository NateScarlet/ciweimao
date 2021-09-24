package client

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NateScarlet/block-cipher/pkg/block_cipher"
	"github.com/tidwall/gjson"
)

// reference https://github.com/zsakvo/Cirno-go/blob/af26c03718a75a86b7198cbdfae0126740e69b55/util/decode.go

func (c Client) apiAESKey() (ret []byte) {
	var v = sha256.Sum256([]byte(c.APIKey))
	return v[:]
}

func (c Client) DecryptAPIResponse(r io.Reader) (ret io.Reader, err error) {
	var decoder = base64.NewDecoder(base64.StdEncoding, r)

	block, err := aes.NewCipher(c.apiAESKey())
	if err != nil {
		return
	}
	var blockMode = cipher.NewCBCDecrypter(block, c.APIInitialVector)

	ret = block_cipher.NewDecrypter(
		decoder,
		blockMode,
		block_cipher.PKCS7Padding{BlockSize: blockMode.BlockSize()},
	)
	return
}

func (c Client) SetDefaultAPIAuthData(data url.Values) {
	if len(data["login_token"]) == 0 {
		data.Set("login_token", c.LoginToken)
	}
	if len(data["device_token"]) == 0 {
		data.Set("device_token", c.DeviceToken)
	}
	if len(data["account"]) == 0 {
		data.Set("account", c.Account)
	}
	if len(data["app_version"]) == 0 {
		data.Set("app_version", c.AppVersion)
	}
}

func (c Client) Call(ctx context.Context, endpoint string, data url.Values) (ret gjson.Result, err error) {
	if data == nil {
		data = url.Values{}
	}
	c.SetDefaultAPIAuthData(data)

	req, err := http.NewRequestWithContext(ctx, "POST", c.EndpointURL(endpoint, nil).String(), strings.NewReader(data.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("ciweimao: api: %s %d", endpoint, resp.StatusCode)
		return
	}

	r, err := c.DecryptAPIResponse(resp.Body)
	if err != nil {
		return
	}

	respData, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	ret = gjson.ParseBytes(respData)
	if code := ret.Get("code").String(); code != "100000" {
		var msg = ret.Raw
		if tip := ret.Get("tip").String(); tip != "" {
			msg = fmt.Sprintf("%s: %s", code, tip)
		}
		err = fmt.Errorf("ciweimao: api: %s: %s", endpoint, msg)
	}

	return
}

func ParseTime(v string) (ret time.Time, err error) {
	return time.ParseInLocation("2006-01-02 15:04:05", v, time.FixedZone("UTC+8", 8*60))
}
