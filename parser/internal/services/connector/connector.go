package connector

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"crypto-news/internal/config"
)

type Connector interface {
	Request(ctx context.Context, r RequestParams) (io.ReadCloser, int, error)
}

type connector struct {
	log *logrus.Entry

	client http.Client
}

func New(cfg config.Config) Connector {
	return &connector{
		log: cfg.Logging().WithField("service", "[CONN]"),
		client: http.Client{
			Timeout: 0,
		},
	}
}

func (c connector) Request(ctx context.Context, r RequestParams) (io.ReadCloser, int, error) {
	c.log.WithField("params", r.Params.Encode()).Debugf("Requesting, %s%s...", r.Url, r.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s?%s", r.Url, r.Path, r.Params.Encode()), nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create new request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to do client request")
	}

	return resp.Body, resp.StatusCode, nil
}
