package browse_ai_crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/store"
	"common/iteration"
	"parser/internal/config"
	"parser/internal/services/connector"
	"parser/internal/services/crawler"
)

const BrowseAI = "browse_ai"

type BrowseAICrawler struct {
	log *logrus.Entry

	authToken string
	url       string
	robotID   string
	conn      connector.Connector

	dataProvider store.DataProvider
}

func NewCrawler(cfg config.Config, robotID string) crawler.Crawler {
	return &BrowseAICrawler{
		log: cfg.Logging().WithField("service", "[BROWSE-AI-CRAWLER]"),

		authToken: cfg.Credentials(BrowseAI, "auth_token"),
		url:       cfg.Credentials(BrowseAI, "url"),
		robotID:   robotID,

		dataProvider: store.New(cfg),

		conn: connector.New(cfg),
	}
}

func (c BrowseAICrawler) Crawl(ctx context.Context) (any, int, error) {
	taskBody, statusCode, err := c.conn.Post(ctx, connector.RequestParams{
		Url:  c.url,
		Path: fmt.Sprintf("/v2/robots/%s/tasks", c.robotID),
		Body: nil,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("Bearer %s", c.authToken)},
		},
	})

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to request browse-ai API")
	}

	var taskRespBody rawBody
	if err := json.NewDecoder(taskBody).Decode(&taskRespBody); err != nil {
		return nil, statusCode, errors.Wrap(err, "failed to decode create task response body")
	}

	if statusCode != http.StatusOK {
		return nil, statusCode, nil
	}

	var rawPollBody io.Reader
	pollRespBody := taskRespBody
	for statusCode == http.StatusOK && pollRespBody.Result.Status == TaskStatusInProgress {
		rawPollBody, statusCode, err = c.conn.Poll(ctx, connector.PollParams{
			Url:  c.url,
			Path: fmt.Sprintf("/v2/robots/%s/tasks/%s", c.robotID, taskRespBody.Result.Id),
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", c.authToken)},
			},
			Params: url.Values{
				"cointelegraph_press_releases_limit": []string{"10"},
			},
		})

		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to poll browse-ai API")
		}

		if err := json.NewDecoder(rawPollBody).Decode(&pollRespBody); err != nil {
			return nil, statusCode, errors.Wrap(err, "failed to decode response body")
		}

		if pollRespBody.Result.Status == TaskStatusInProgress {
			time.Sleep(15 * time.Second)
		}
	}

	if statusCode != http.StatusOK {
		return pollRespBody, statusCode, nil
	}

	if pollRespBody.Result.Status == TaskStatusFailed {
		return nil, statusCode, errors.Wrapf(errors.New("task run failed"), "%v", pollRespBody.Result.UserFriendlyError)
	}

	return iteration.Map(pollRespBody.Result.CapturedLists.Releases, toModel), statusCode, nil
}

func toModel(b body) crawler.ParsedBody {
	return b
}
