package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"

	"parser/internal/config"
	"parser/internal/services/crawler"
)

type Worker interface {
	Produce(ctx context.Context, cf []crawler.Crawler)
	Work(ctx context.Context) []Task
}

type worker struct {
	log *logrus.Entry

	workersNum   int
	tasksPending *atomic.Int32

	rl          *rate.Limiter
	tasksChan   chan Task
	resultsChan chan Task
}

func New(workersNum int, cfg config.Config) Worker {
	return &worker{
		log: cfg.Logging().WithField("service", "[WORKER]"),

		workersNum: workersNum,

		tasksPending: atomic.NewInt32(0),
		rl:           rate.NewLimiter(rate.Every(1*time.Second), cfg.RateLimit()),

		tasksChan:   make(chan Task, workersNum),
		resultsChan: make(chan Task, workersNum),
	}
}

func (w worker) Work(ctx context.Context) []Task {
	// distribute
	for i := 0; i < w.workersNum; i++ {
		go func(name string) {
			w.work(ctx, name)
		}(fmt.Sprintf("w-%d", i))
	}

	results := make([]Task, 0, 10)
	for t := range w.resultsChan {
		results = append(results, t)
	}

	return results
}

func (w worker) work(ctx context.Context, name string) {
	for task := range w.tasksChan {
		w.log.Debugf("Running task: %d", task.seq)

		start := time.Now()
		body, code, err := task.do(ctx)
		if err != nil {
			task.Err = err
		}

		switch t := body.(type) {
		case []crawler.ParsedBody:
			task.Body = t
		case map[string]json.RawMessage:
			task.StatusBody = t
		}

		task.StatusCode = code
		task.handleBy = name
		task.duration = time.Duration(time.Since(start).Milliseconds())

		w.log.Debugf("Finishing task: %s", task.String())

		w.resultsChan <- task

		if w.tasksPending.Dec() == 0 {
			close(w.resultsChan)
		}
	}
}
