package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"common"
	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"common/iteration"
	"common/math"
	"gpt/internal/bot"
	"gpt/internal/config"
)

const (
	processingLimit = 10
)

type Service interface {
	Run(ctx context.Context) error
}

type service struct {
	cfg config.Config
	log *logrus.Entry

	dataProvider store.DataProvider
}

func New(cfg config.Config) Service {
	return &service{
		cfg: cfg,
		log: cfg.Logging().WithField("service", "[GPT]"),

		dataProvider: store.New(cfg),
	}
}

func (s service) Run(ctx context.Context) error {
	s.log.Info("Staring gpt generator bot service...")
	summarizationBot := bot.NewOpenAI(s.cfg)

	common.RunEveryWithBackoff(s.cfg.GenerateEvery(), 15*time.Second, 15*time.Minute, func() error {
		s.log.Debug("Generating digest...")

		totalRows, err := s.dataProvider.RawNewsProvider().Count(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to count raw news")
		}

		for i := 0; i < int(totalRows)/processingLimit+1; i++ {
			rawNews, err := s.dataProvider.RawNewsProvider().Order("id", data.OrderAsc).Limit(processingLimit).Offset(uint64(i * processingLimit)).Select(ctx)
			if err != nil {
				if !errors.Is(err, data.ErrNotFound) {
					return errors.Wrap(err, "failed to count raw news")
				}
				s.log.Debug("Done processing pending raw news")
				return nil
			}

			titleIDs := iteration.Map(rawNews, func(t model.RawNews) uuid.UUID {
				return t.TitleID
			})

			titles, err := s.dataProvider.TitlesProvider().ByIDs(titleIDs).Select(ctx)

			aggregatedTextBuf := bytes.NewBuffer(make([]byte, 0, 1000))

			for _, rawNewsPiece := range rawNews {
				aggregatedTextBuf.WriteString(convert.FromPtr(rawNewsPiece.Body))
			}

			for _, locale := range s.cfg.Locales() {
				s.log.WithField("locale", locale).Debug("Generating for locale")

				timestamp := common.CurrentTimestamp()

				news, digestResponse, err := s.generateDigestForLocale(
					ctx,
					summarizationBot,
					s.cfg.QueryContext(), aggregatedTextBuf.String(), locale,
					titles,
					timestamp,
				)
				if err != nil {
					return errors.Wrapf(err, "failed to generate for locale: %s", locale)
				}

				if err := s.addNews(ctx, news, digestResponse); err != nil {
					return errors.Wrap(err, "failed to add news")
				}
			}

			rawNewsIDs := iteration.Map(rawNews, func(t model.RawNews) uuid.UUID {
				return t.ID
			})

			if err := s.dataProvider.RawNewsProvider().ByIDs(rawNewsIDs).Remove(ctx, model.RawNews{}); err != nil {
				return errors.Wrap(err, "failed to remove processed raw news")
			}
		}

		return nil
	})

	s.log.Info("Finishing gpt generator bot service...")

	return nil
}

func (s service) generateDigestForLocale(ctx context.Context,
	bot bot.Bot,
	queryContext, aggregatedText, locale string,
	titles []model.Title,
	timestamp time.Time) (*model.News, []model.Coin, error) {
	// we shouldn't create single post longer than 10 minutes, if that happens - probably something went wrong
	deadlineCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
	defer cancel()

	lang := display.English.Tags().Name(language.Make(locale))

	replyMsg, err := bot.Ask(deadlineCtx, aggregatedText[:math.Min(len(aggregatedText), maxInputChars)], queryContext, lang)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to ask bot")
	}

	content, coins := parseCoins(replyMsg.Text)

	resourcesList := make([]model.NewsMediaResource, 0, len(titles))
	for i, title := range titles {
		metaLinks := model.MetaLinksData{
			ID:    strconv.Itoa(i),
			URL:   convert.FromPtr(title.URL),
			Title: convert.FromPtr(title.Title),
		}

		metaLinksBody, err := json.Marshal(metaLinks)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to marshal meta sources body")
		}

		resourcesList = append(resourcesList, model.NewsMediaResource{
			Type: convert.ToPtr(model.ResourceTypeSource),
			URL:  title.URL,
			Meta: metaLinksBody,
		})
	}

	news := &model.News{
		Locale: convert.ToPtr(locale),
		Media: &model.NewsMedia{
			Title:     convert.ToPtr(fmt.Sprintf("Digest hour: %d, Day: %d", timestamp.Hour(), timestamp.Day())),
			Text:      convert.ToPtr(content),
			Resources: resourcesList,
		},
		Source: convert.ToPtr("gpt-bing"),
		Status: convert.ToPtr(model.StatusPending),
	}

	s.log.WithFields(logrus.Fields{
		"digest-hour": timestamp.Hour(),
		"digest-day":  timestamp.Day(),
	}).Debug("Finished generating")

	return news, coins, nil
}

func (s service) addNews(ctx context.Context, news *model.News, coins []model.Coin) error {
	createdNews, err := s.dataProvider.NewsProvider().Insert(ctx, convert.FromPtr(news))
	if err != nil {
		return errors.Wrap(err, "failed to insert news digest")
	}

	newsCoinsBatch := createCoinsNewsCoinsBatch(createdNews.ID, coins)
	if err = s.dataProvider.CoinsProvider().UpsertCoinsBatch(ctx, coins); err != nil {
		return errors.Wrap(err, "failed to insert batch of coins")
	}

	if err = s.dataProvider.NewsCoinsProvider().InsertBatch(ctx, newsCoinsBatch); err != nil {
		return errors.Wrap(err, "failed to insert batch of news-coins")
	}

	channels, err := s.dataProvider.ChannelsProvider().Select(ctx)
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return errors.Wrap(err, "failed to select channels")
		}
	}

	newsChannelsBatch := toNewsChannelsBatch(createdNews, channels)
	if err = s.dataProvider.NewsChannelsProvider().InsertBatch(ctx, newsChannelsBatch); err != nil {
		return errors.Wrap(err, "failed to insert batch of news-channels")
	}

	return nil
}
