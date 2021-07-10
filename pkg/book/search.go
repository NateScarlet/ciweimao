package book

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/NateScarlet/ciweimao/pkg/client"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Order string

var (
	OAverageSale     Order = "average_buy"
	OMonthlyTicket   Order = "total_yp"
	OMonthNoVIPClick Order = "month_no_vip_click"
	OTotalBookmark   Order = "total_favor"
	OTotalClick      Order = "total_click"
	OTotalRecommend  Order = "total_recommend"
	OTotalWordCount  Order = "total_word_count"
	OUpdateTime      Order = "uptime"
	OWeekClick       Order = "week_click"
	OWeekNoVIPClick  Order = "week_no_vip_click"
)

type WordCountRange string

var (
	WordCountUndefined   WordCountRange = ""
	WordCountLt300k      WordCountRange = "1"
	WordCountGt300kLt50k WordCountRange = "2"
	WordCountGt500kLt1m  WordCountRange = "3"
	WordCountGt1mLt2m    WordCountRange = "4"
	WordCountGt2m        WordCountRange = "5"
)

type Expense string

var (
	ExpenseUndefined Expense = ""
	ExpenseFree      Expense = "0"
	ExpensePaid      Expense = "1"
)

type UpdateTimeRange string

var (
	UpdateTimeUndefined   UpdateTimeRange = ""
	UpdateTimeIn3Days     UpdateTimeRange = "1"
	UpdateTimeIn7Days     UpdateTimeRange = "2"
	UpdateTimeInHalfMonth UpdateTimeRange = "3"
	UpdateTimeInMonth     UpdateTimeRange = "4"
)

// SearchOptions contains merged SearchOption result
type SearchOptions struct {
	category     Category
	pageIndex    int
	pageSize     int
	order        Order
	tags         []string
	wordCount    WordCountRange
	updateStatus UpdateStatus
	expense      Expense
	query        string
	updateTime   UpdateTimeRange
}

type SearchOption = func(opts *SearchOptions)

func SearchOptionCategory(category Category) SearchOption {
	return func(opts *SearchOptions) {
		opts.category = category
	}
}

func SearchOptionWordCount(wordCountRange WordCountRange) SearchOption {
	return func(opts *SearchOptions) {
		opts.wordCount = wordCountRange
	}
}

// SearchOptionPageSize specify result count, defaults to 10
func SearchOptionPageSize(n int) SearchOption {
	return func(opts *SearchOptions) {
		opts.pageSize = n
	}
}

func SearchOptionPageIndex(n int) SearchOption {
	return func(opts *SearchOptions) {
		opts.pageIndex = n
	}
}

func SearchOptionExpense(expense Expense) SearchOption {
	return func(opts *SearchOptions) {
		opts.expense = expense
	}
}

func SearchOptionQuery(query string) SearchOption {
	return func(opts *SearchOptions) {
		opts.query = query
	}
}

func SearchOptionTag(tags ...string) SearchOption {
	return func(opts *SearchOptions) {
		opts.tags = tags
	}
}

func marshallSearchTags(tags []string) (ret string, err error) {
	ret = "[]"
	for index, tag := range tags {
		ret, err = sjson.Set(ret, fmt.Sprintf("%d.filter", index), "1")
		if err != nil {
			return
		}

		ret, err = sjson.Set(ret, fmt.Sprintf("%d.tag", index), tag)
		if err != nil {
			return
		}
	}
	return
}

func SearchOptionUpdateStatus(updateStatus UpdateStatus) SearchOption {
	return func(opts *SearchOptions) {
		opts.updateStatus = updateStatus
	}
}

func SearchOptionUpdateTime(updateTimeRange UpdateTimeRange) SearchOption {
	return func(opts *SearchOptions) {
		opts.updateTime = updateTimeRange
	}
}
func SearchOptionOrder(order Order) SearchOption {
	return func(opts *SearchOptions) {
		opts.order = order
	}
}

type SearchResult struct {
	JSON gjson.Result
}

func (res SearchResult) Books() (ret []Book) {
	var bookList = res.JSON.Get("data.book_list")
	ret = make([]Book, 0, bookList.Get("#").Int())
	bookList.ForEach(func(_, v gjson.Result) bool {
		var book = Book{}
		book.unmarshalBookInfo(v)
		ret = append(ret, book)
		return true
	})
	return

}

func Search(ctx context.Context, opts ...SearchOption) (ret SearchResult, err error) {
	var args = new(SearchOptions)
	args.pageSize = 10
	for _, i := range opts {
		i(args)
	}

	tags, err := marshallSearchTags(args.tags)
	if err != nil {
		return
	}

	var data = url.Values{
		"category_index": []string{string(args.category)},
		"count":          []string{strconv.Itoa(args.pageSize)},
		"filter_uptime":  []string{string(args.updateTime)},
		"filter_word":    []string{string(args.wordCount)},
		"is_paid":        []string{string(args.expense)},
		"key":            []string{args.query},
		"order":          []string{string(args.order)},
		"page":           []string{strconv.Itoa(args.pageIndex)},
		"tags":           []string{tags},
		"up_status":      []string{string(args.updateStatus)},
	}

	var c = client.For(ctx)
	ret.JSON, err = c.Call(ctx, "/bookcity/get_filter_search_book_list", data)
	return

}
