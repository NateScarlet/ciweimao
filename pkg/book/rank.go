package book

import (
	"context"
	"net/url"
	"strconv"

	"github.com/NateScarlet/ciweimao/pkg/client"
	"github.com/tidwall/gjson"
)

type RankPeriod string

var (
	RPWeek  RankPeriod = "week"
	RPMonth RankPeriod = "month"
)

type RankType string

var (
	// 点击
	RTClick RankType = "no_vip_click"
	// 畅销
	RTSales RankType = "fans_value"
	// 月票
	RTMonthlyTicket RankType = "yp"
	// 新书
	RTNewBook RankType = "yp_new"
	// 收藏
	RTBookmark RankType = "favor"
	// 推荐
	RTRecommend RankType = "recommend"
	// 刀片
	RTBlade RankType = "blade"
	// 更新
	RTWordCount RankType = "word_count"
	// 吐槽
	RTTsukkomi RankType = "tsukkomi"
	// 完本
	RTFinished RankType = "complet"
	// 追读
	RTWatching RankType = "track_read"
)

type RankResult struct {
	JSON gjson.Result
}

func (res RankResult) Books() []Book {
	var bookList = res.JSON.Get("data.book_list")
	var size = bookList.Get("#").Int()
	var ret = make([]Book, 0, size)
	bookList.ForEach(func(k, v gjson.Result) bool {
		var book = Book{}
		book.unmarshalBookInfo(v)
		ret = append(ret, book)
		return true
	})
	return ret
}

type RankOptions struct {
	pageSize  int
	pageIndex int
	category  Category
}

type RankOption = func(opts *RankOptions)

// RankOptionPageSize, defaults to 10.
// max page size is 20 (2021-06-28)
func RankOptionPageSize(n int) RankOption {
	return func(opts *RankOptions) {
		opts.pageSize = n
	}
}

func RankOptionPageIndex(n int) RankOption {
	return func(opts *RankOptions) {
		opts.pageIndex = n
	}
}

func RankOptionCategory(category Category) RankOption {
	return func(opts *RankOptions) {
		opts.category = category
	}
}

func Rank(ctx context.Context, rankType RankType, period RankPeriod, opts ...RankOption) (ret RankResult, err error) {
	var args = new(RankOptions)
	args.pageSize = 10
	for _, i := range opts {
		i(args)
	}

	var c = client.For(ctx)
	ret.JSON, err = c.Call(ctx, "/bookcity/get_rank_book_list", url.Values{
		"order":          []string{string(rankType)},
		"time_type":      []string{string(period)},
		"category_index": []string{string(args.category)},
		"count":          []string{strconv.Itoa(args.pageSize)},
		"page":           []string{strconv.Itoa(args.pageIndex)},
	})
	return
}
