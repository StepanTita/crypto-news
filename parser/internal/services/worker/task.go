package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"parser/internal/services/crawler"
)

type Task struct {
	seq        int // task number seq
	do         crawler.CrawlFunc
	Body       []crawler.ParsedBody
	Err        error
	StatusCode int
	StatusBody map[string]json.RawMessage
	duration   time.Duration
	handleBy   string // worker name
}

func (t Task) String() string {
	return fmt.Sprintf("seq: %d, handleBy: %s, duration: %dms, StatusCode: %d, Err: %v ...", t.seq, t.handleBy, t.duration, t.StatusCode, t.Err)
}
