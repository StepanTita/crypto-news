package worker

import (
	"fmt"
	"time"

	"crypto-news/internal/services/crawler"
)

type Task struct {
	seq        int // task number seq
	do         crawler.CrawlFunc
	Body       []crawler.ParsedBody
	Err        error
	StatusCode int
	duration   time.Duration
	handleBy   string // worker name
}

func (t Task) String() string {
	return fmt.Sprintf("seq: %d, handleBy: %s, duration: %dms, StatusCode: %d, Err: %v ...", t.seq, t.handleBy, t.duration, t.StatusCode, t.Err)
}
