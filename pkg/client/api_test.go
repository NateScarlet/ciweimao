package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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
