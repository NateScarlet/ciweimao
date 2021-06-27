package book

type Book struct {
	ID          string
	Title       string
	CoverURL    string
	Description string

	WordCount          uint64
	ClickCount         uint64
	RecommendCount     uint64
	BookmarkCount      uint64
	MonthlyTicketCount uint64
	BladeCount         uint64
	FanValue           uint64

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

	LastChapter Chapter
	Tags        []Tag
}
