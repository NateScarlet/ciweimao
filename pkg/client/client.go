package client

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	APIInitialVector []byte
	APIKey           string
	ServerURL        string
	LoginToken       string
	DeviceToken      string
	Account          string
	AppVersion       string

	http.Client
}

// EndpointURL returns url for server endpint.
func (c Client) EndpointURL(path string, values *url.Values) *url.URL {
	s := c.ServerURL
	if s == "" {
		s = DefaultServerURL
	}

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	u.Path = path
	if values != nil {
		u.RawQuery = values.Encode()
	}
	return u
}

// Default client
var Default = new(Client)

func getenvBase64(name string) []byte {
	var encoded = os.Getenv("name")
	if encoded == "" {
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil
	}
	return decoded
}

var DefaultAPIInitialVector []byte = getenvBase64("CIWEIMAO_API_INITIAL_VECTOR")
var DefaultAPIKey string = os.Getenv("CIWEIMAO_API_KEY")
var DefaultServerURL string = os.Getenv("CIWEIMAO_SERVER_URL")
var DefaultAccount string = os.Getenv("CIWEIMAO_ACCOUNT")
var DefaultLoginToken string = os.Getenv("CIWEIMAO_LOGIN_TOKEN")
var DefaultDeviceToken string = os.Getenv("CIWEIMAO_DEVICE_TOKEN")
var DefaultAppVersion string = os.Getenv("CIWEIMAO_APP_VERSION")

// DefaultUserAgent for new clients
var DefaultUserAgent = os.Getenv("CIWEIMAO_USER_AGENT")

func init() {
	if DefaultServerURL == "" {
		DefaultServerURL = "https://app.hbooker.com"
	}

	// default config from https://github.com/zsakvo/Cirno-go/blob/af26c03718a75a86b7198cbdfae0126740e69b55/config/read.go
	if DefaultUserAgent == "" {
		DefaultUserAgent = "Android com.kuangxiangciweimao.novel"
	}
	if len(DefaultAPIInitialVector) != 16 {
		DefaultAPIInitialVector = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	if DefaultAPIKey == "" {
		DefaultAPIKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"
	}
	if DefaultAppVersion == "" {
		DefaultAppVersion = "2.6.011"
	}

	Default.APIInitialVector = DefaultAPIInitialVector
	Default.APIKey = DefaultAPIKey
	Default.Account = DefaultAccount
	Default.LoginToken = DefaultLoginToken
	Default.DeviceToken = DefaultDeviceToken
	Default.AppVersion = DefaultAppVersion
	Default.SetDefaultHeader("User-Agent", DefaultUserAgent)
}
