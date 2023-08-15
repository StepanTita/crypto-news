package connector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"parser/internal/config"
)

type Connector interface {
	Post(ctx context.Context, r RequestParams) (io.Reader, int, error)
	Poll(ctx context.Context, r PollParams) (io.Reader, int, error)
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

func (c connector) Poll(ctx context.Context, r PollParams) (io.Reader, int, error) {
	c.log.WithField("params", r.Params.Encode()).Debugf("Requesting, %s%s...", r.Url, r.Path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s?%s", r.Url, r.Path, r.Params.Encode()), nil)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to create new get request")
	}

	req.Header = r.Headers
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to do get client request")
	}

	return resp.Body, resp.StatusCode, nil
}

func (c connector) Post(ctx context.Context, r RequestParams) (io.Reader, int, error) {
	//c.log.WithField("body", string(r.Body)).Debugf("Requesting, %s%s...", r.Url, r.Path)
	//
	//req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", r.Url, r.Path), bytes.NewReader(r.Body))
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "failed to create new post request")
	//}
	//
	//req.Header = r.Headers
	//resp, err := c.client.Do(req)
	//if err != nil {
	//	return nil, 0, errors.Wrap(err, "failed to do post client request")
	//}
	//
	//return resp.Body, resp.StatusCode, nil
	return strings.NewReader(`{
"result": {"id": "eff43fe4-3f23-4052-bd4c-9bc51dc43c20", "status": "in-progress"}
}`), http.StatusOK, nil
}
