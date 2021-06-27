package book

import "time"

type Chapter struct {
	ID     string
	BookID string
	Title  string
	Index  uint64

	Uploaded time.Time
	Modified time.Time
}
