package book

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/NateScarlet/snapshot/pkg/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookFetch(t *testing.T) {
	var ctx = context.Background()

	var book = new(Book)
	book.ID = "100011781"

	res, err := book.Fetch(ctx)
	require.NoError(t, err)
	if snapshot.DefaultUpdate {
		var doc = new(interface{})
		err = json.Unmarshal([]byte(res.Raw), doc)
		require.NoError(t, err)
		snapshot.MatchJSON(t, *doc, snapshot.OptionExt(".response.json"))
		snapshot.MatchJSON(t, book, snapshot.OptionExt(".book.json"))
	}

	assert.Equal(t, "100011781", book.ID)
	assert.NotEmpty(t, book.Author)
	assert.NotEmpty(t, book.Description)
	assert.NotEmpty(t, book.Title)
	assert.NotEmpty(t, book.ChapterCount)
	assert.NotEmpty(t, book.Category)
	assert.NotEmpty(t, book.CoverURL)
	assert.NotEmpty(t, book.Created)
	assert.NotEmpty(t, book.Updated)
	assert.NotEmpty(t, book.UpdateStatus)
	assert.NotEmpty(t, book.WordCount)
	assert.NotEmpty(t, book.Tags)
	assert.NotEmpty(t, book.LastChapter)
}
