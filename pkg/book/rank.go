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

type Rank struct {
	Type     RankType
	Period   RankPeriod
	Category Category
}

type RankFetchResult struct {
	JSON gjson.Result
}

func (res RankFetchResult) Books() (ret []Book) {
	var bookList = res.JSON.Get("data.book_list")
	var size = bookList.Get("#").Int()
	ret = make([]Book, 0, size)
	bookList.ForEach(func(k, v gjson.Result) bool {
		var book = Book{
			ID:                   v.Get("book_id").String(),
			Title:                v.Get("book_name").String(),
			CoverURL:             v.Get("cover").String(),
			Description:          v.Get("description").String(),
			WordCount:            v.Get("total_word_count").Uint(),
			ClickCount:           v.Get("total_click").Uint(),
			RecommendCount:       v.Get("total_recommend").Uint(),
			BookmarkCount:        v.Get("total_favor").Uint(),
			MonthlyTicketCount:   v.Get("total_yp").Uint(),
			BladeCount:           v.Get("total_blade").Uint(),
			FanValue:             v.Get("total_fans_value").Uint(),
			MonthClickCount:      v.Get("month_click").Uint(),
			MonthNoVIPClickCount: v.Get("month_no_vip_click").Uint(),
			MonthRecommendCount:  v.Get("month_recommend").Uint(),
			MonthBookmarkCount:   v.Get("month_favor").Uint(),
			MonthFanValue:        v.Get("month_fans_value").Uint(),
			WeekClickCount:       v.Get("week_click").Uint(),
			WeekNoVIPClickCount:  v.Get("week_no_vip_click").Uint(),
			WeekRecommendCount:   v.Get("week_recommend").Uint(),
			WeekBookmarkCount:    v.Get("week_favor").Uint(),
			WeekFanValue:         v.Get("week_fans_value").Uint(),
		}
		var chapterInfo = v.Get("last_chapter_info")
		book.LastChapter = Chapter{
			ID:     chapterInfo.Get("chapter_id").String(),
			BookID: book.ID,
			Title:  chapterInfo.Get("chapter_title").String(),
			Index:  chapterInfo.Get("chapter_index").Uint(),
		}
		book.LastChapter.Uploaded, _ = client.ParseTime(chapterInfo.Get("uptime").String())
		book.LastChapter.Modified, _ = client.ParseTime(chapterInfo.Get("mtime").String())
		v.Get("tag_list").ForEach(func(key, value gjson.Result) bool {
			book.Tags = append(book.Tags, Tag{
				ID:   value.Get("tag_id").String(),
				Name: value.Get("tag_name").String(),
			})
			return true
		})
		ret = append(ret, book)
		return true
	})
	return
}

// Fetch rank with pagination.
// max page size is 20 (2021-06-28)
func (r Rank) Fetch(ctx context.Context, pageSize int, pageIndex int) (ret RankFetchResult, err error) {
	var c = client.For(ctx)
	res, err := c.Call(ctx, "/bookcity/get_rank_book_list", url.Values{
		"order":          []string{string(r.Type)},
		"time_type":      []string{string(r.Period)},
		"category_index": []string{string(r.Category)},
		"count":          []string{strconv.Itoa(pageSize)},
		"page":           []string{strconv.Itoa(pageIndex)},
	})
	if err != nil {
		return
	}
	ret.JSON = res
	return
}
