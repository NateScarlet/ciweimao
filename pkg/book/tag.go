package book

import (
	"context"

	"github.com/NateScarlet/ciweimao/pkg/client"
	"github.com/tidwall/gjson"
)

type Tag struct {
	ID   string
	Name string
	Type string
}

type OfficialTagsResult struct {
	JSON gjson.Result
}

func (res OfficialTagsResult) Tags() (ret []Tag) {
	var tagList = res.JSON.Get("data.official_tag_list")
	ret = make([]Tag, 0, tagList.Get("#").Int())
	tagList.ForEach(func(key, value gjson.Result) bool {
		ret = append(ret, Tag{
			Name: value.Get("tag_name").String(),
			Type: value.Get("tag_type").String(),
		})
		return true
	})
	return
}

func OfficialTags(ctx context.Context) (ret OfficialTagsResult, err error) {
	var c = client.For(ctx)
	ret.JSON, err = c.Call(ctx, "/book/get_official_tag_list", nil)
	return
}
