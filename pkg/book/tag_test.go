package book

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NateScarlet/snapshot/pkg/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOfficialTags(t *testing.T) {
	var ctx = context.Background()

	res, err := OfficialTags(ctx)
	require.NoError(t, err)
	if snapshot.DefaultUpdate {
		var doc = new(interface{})
		err = json.Unmarshal([]byte(res.JSON.Raw), doc)
		require.NoError(t, err)
		snapshot.MatchJSON(t, *doc, snapshot.OptionExt(".response.json"))
	}

	tags := res.Tags()
	assert.NotEmpty(t, tags)
	for _, tag := range tags {
		assert.NotEmpty(t, tag.Name)
		assert.NotEmpty(t, tag.Type)
	}
	if snapshot.DefaultUpdate {
		require.NoError(t, err)
		snapshot.MatchJSON(t, tags, snapshot.OptionExt(".tags.json"))
	}
}
