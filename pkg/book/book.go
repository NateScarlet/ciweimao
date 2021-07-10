package book

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/NateScarlet/ciweimao/pkg/client"
	"github.com/tidwall/gjson"
)

type UpdateStatus string

var (
	UpdateStatusUndefined   UpdateStatus = ""
	UpdateStatusNotFinished UpdateStatus = "0"
	UpdateStatusFinished    UpdateStatus = "1"
)

type Book struct {
	ID           string
	Title        string
	CoverURL     string
	Description  string
	Author       string
	Category     Category
	Created      time.Time
	Updated      time.Time
	UpdateStatus UpdateStatus

	WordCount          uint64
	ChapterCount       uint64
	ClickCount         uint64
	RecommendCount     uint64
	BookmarkCount      uint64
	MonthlyTicketCount uint64
	BladeCount         uint64
	FanValue           uint64
	ReviewCount        uint64
	RewardCount        uint64

	MonthClickCount      uint64
	MonthNoVIPClickCount uint64
	MonthRecommendCount  uint64
	MonthBookmarkCount   uint64
	MonthFanValue        uint64

	WeekClickCount      uint64
	WeekNoVIPClickCount uint64
	WeekRecommendCount  uint64
	WeekBookmarkCount   uint64
	WeekFanValue        uint64

	TotalBladeCount         uint64
	TotalMonthlyTicketCount uint64

	LastChapter Chapter
	Tags        []Tag
}

func (book *Book) unmarshalBookInfo(bookInfo gjson.Result) (err error) {
	book.Author = bookInfo.Get("author_name").String()
	book.BladeCount = bookInfo.Get("current_blade").Uint()
	book.BookmarkCount = bookInfo.Get("total_favor").Uint()
	book.Category = Category(bookInfo.Get("category_index").String())
	book.ChapterCount = bookInfo.Get("chapter_amount").Uint()
	book.ClickCount = bookInfo.Get("total_click").Uint()
	book.CoverURL = bookInfo.Get("cover").String()
	book.Description = bookInfo.Get("description").String()
	book.FanValue = bookInfo.Get("total_fans_value").Uint()
	book.ID = bookInfo.Get("book_id").String()
	book.LastChapter.BookID = bookInfo.Get("last_chapter_info.book_id").String()
	book.LastChapter.ID = bookInfo.Get("last_chapter_info.chapter_id").String()
	book.LastChapter.Index = bookInfo.Get("last_chapter_info.index").Uint()
	book.LastChapter.Title = bookInfo.Get("last_chapter_info.chapter_title").String()
	book.MonthBookmarkCount = bookInfo.Get("month_favor").Uint()
	book.MonthClickCount = bookInfo.Get("month_click").Uint()
	book.MonthFanValue = bookInfo.Get("month_fans_value").Uint()
	book.MonthlyTicketCount = bookInfo.Get("current_yp").Uint()
	book.MonthNoVIPClickCount = bookInfo.Get("month_no_vip_click").Uint()
	book.MonthRecommendCount = bookInfo.Get("month_recommend").Uint()
	book.RecommendCount = bookInfo.Get("total_recommend").Uint()
	book.ReviewCount = bookInfo.Get("review_amount").Uint()
	book.RewardCount = bookInfo.Get("reward_amount").Uint()
	book.Title = bookInfo.Get("book_name").String()
	book.TotalBladeCount = bookInfo.Get("total_blade").Uint()
	book.TotalMonthlyTicketCount = bookInfo.Get("total_yp").Uint()
	book.UpdateStatus = UpdateStatus(bookInfo.Get("up_status").String())
	book.WeekBookmarkCount = bookInfo.Get("week_favor").Uint()
	book.WeekClickCount = bookInfo.Get("week_click").Uint()
	book.WeekFanValue = bookInfo.Get("week_fans_value").Uint()
	book.WeekNoVIPClickCount = bookInfo.Get("week_no_vip_click").Uint()
	book.WeekRecommendCount = bookInfo.Get("week_recommend").Uint()
	book.WordCount = bookInfo.Get("total_word_count").Uint()

	var tagList = bookInfo.Get("tag_list")
	book.Tags = make([]Tag, 0, tagList.Get("#").Int())
	tagList.ForEach(func(_, v gjson.Result) bool {
		var tag = Tag{
			ID:   v.Get("tag_id").String(),
			Name: v.Get("tag_name").String(),
			Type: v.Get("tag_type").String(),
		}
		book.Tags = append(book.Tags, tag)
		return true
	})

	book.LastChapter.Created, err = client.ParseTime(bookInfo.Get("last_chapter_info.uptime").String())
	if err != nil {
		return
	}
	book.LastChapter.Updated, err = client.ParseTime(bookInfo.Get("last_chapter_info.mtime").String())
	if err != nil {
		return
	}
	if v := bookInfo.Get("newtime").String(); v != "" {
		book.Created, err = client.ParseTime(v)
		if err != nil {
			return
		}
	}
	book.Updated, err = client.ParseTime(bookInfo.Get("uptime").String())
	if err != nil {
		return
	}
	return
}

// Fetch update book in-place and returns json response result.
func (book *Book) Fetch(ctx context.Context) (ret gjson.Result, err error) {
	if book.ID == "" {
		err = fmt.Errorf("ciweimao: book: Book.Fetch: id is empty")
		return
	}

	var c = client.For(ctx)
	var data = url.Values{
		"book_id": []string{book.ID},
	}
	ret, err = c.Call(ctx, "/book/get_info_by_id", data)
	if err != nil {
		return
	}
	err = book.unmarshalBookInfo(ret.Get("data.book_info"))
	return
}
