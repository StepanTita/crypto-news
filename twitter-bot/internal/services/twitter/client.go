package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/constraints"
	"golang.org/x/oauth2"

	commonutils "common"
	config2 "common/config"

	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"twitter-bot/internal/config"
)

const (
	appJsonContentType = "application/json; charset=utf-8"
)

type Client interface {
	PostTweet(ctx context.Context, news model.News) error
}

type twitterClient struct {
	log       *logrus.Entry
	templator config2.Templator

	oauthConfig *oauth2.Config

	dataProvider store.DataProvider

	token *oauth2.Token
}

func New(cfg config.Config) Client {
	return &twitterClient{
		log:       cfg.Logging().WithField("service", "[TWITTER-CLIENT]"),
		templator: cfg,

		oauthConfig: cfg.OAuthConfig(),

		dataProvider: store.New(cfg),

		token: &oauth2.Token{},
	}
}

func (t *twitterClient) PostTweet(ctx context.Context, news model.News) error {
	if !t.isFresh() {
		if err := t.refresh(ctx); err != nil {
			return errors.Wrap(err, "failed to refresh token")
		}
	}

	tweet := Tweet{
		// TODO: fix temporary workaround till we fix the markdown issue in twitter
		Text: fmt.Sprintf(t.templator.Template(fmt.Sprintf("%s_%s", commonutils.NewsPost, "twitter")),
			convert.FromPtr(news.Media.Title),
			convert.FromPtr(news.Media.Text),
			convert.FromPtr(news.Source)),
	}
	tweet.Text = tweet.Text[:Min(len(tweet.Text), 260)]

	body, err := json.Marshal(tweet)
	if err != nil {
		return errors.Wrap(err, "failed to marshal tweet")
	}

	resp, err := t.oauthConfig.Client(ctx, t.token).Post("https://api.twitter.com/2/tweets", appJsonContentType, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "failed to send tweet post request")
	}

	if resp.StatusCode != http.StatusCreated {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to read response body, request failed with status code: %d", resp.StatusCode))
		}
		return errors.New(fmt.Sprintf("failed to create tweet, request failed with status code: %d, response body: %s", resp.StatusCode, b))
	}

	return nil
}

func (t *twitterClient) isFresh() bool {
	return t.token.Valid()
}

func (t *twitterClient) refresh(ctx context.Context) error {
	if err := t.dataProvider.KVProvider().GetStruct(ctx, model.ToKey(t.token, false), &t.token); err != nil {
		if errors.Is(err, data.ErrNotFound) {
			t.log.Warn("Auth token not found for twitter API to post!")
			return nil
		}
		return errors.Wrap(err, "failed to get token from redis store")
	}
	return nil
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
