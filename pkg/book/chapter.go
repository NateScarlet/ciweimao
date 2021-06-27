package book

import "time"

type Chapter struct {
	ID     string
	BookID string
	Title  string
	Index  uint64

	Created time.Time
	Updated time.Time
}
