package book

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NateScarlet/snapshot/pkg/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRank(t *testing.T) {
	var ctx = context.Background()

	res, err := Rank(ctx, RTClick, RPWeek)
	require.NoError(t, err)

	if snapshot.DefaultUpdate {
		var doc = new(interface{})
		err = json.Unmarshal([]byte(res.JSON.Raw), doc)
		require.NoError(t, err)
		snapshot.MatchJSON(t, *doc, snapshot.OptionExt(".response.json"))
	}

	var books = res.Books()
	assert.Len(t, books, 10)
	for _, book := range books {
		assert.NotEmpty(t, book.ID)
		assert.NotEmpty(t, book.Title)
		assert.NotEmpty(t, book.Description)
	}
	if snapshot.DefaultUpdate {
		require.NoError(t, err)
		snapshot.MatchJSON(t, books, snapshot.OptionExt(".books.json"))
	}
}
