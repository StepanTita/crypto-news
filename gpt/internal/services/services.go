package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common"
	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"gpt/internal/config"
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
	bot := chat_bot.New(s.cfg.GPTConfig())

	common.RunEveryWithBackoff(s.cfg.GenerateEvery(), 15*time.Second, 15*time.Minute, func() error {
		s.log.Debug("Generating digest...")

		err := bot.Init(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to initialize bot")
		}

		// we shouldn't create single post longer than 10 minutes, if that happens - probably something went wrong
		deadlineCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Minute))
		defer cancel()

		parsedResponsesChan, err := bot.Ask(deadlineCtx, s.cfg.GPTConfig().InitialPrompt(), s.cfg.GPTConfig().Style(), true)
		if err != nil {
			return errors.Wrap(err, "failed to ask bot")
		}

		digestResponse, err := s.readResponses(deadlineCtx, parsedResponsesChan)
		if err != nil {
			return errors.Wrap(err, "failed to generate digest")
		}

		resourcesList := make([]model.NewsMediaResource, 0, len(digestResponse.links))
		for _, link := range digestResponse.links {
			metaLinks := model.MetaLinksData{
				ID:    link.ID,
				URL:   link.URL,
				Title: link.Title,
			}

			metaLinksBody, err := json.Marshal(metaLinks)
			if err != nil {
				return errors.Wrap(err, "failed to marshal meta links body")
			}
			resourcesList = append(resourcesList, model.NewsMediaResource{
				Type: convert.ToPtr(model.ResourceTypeSource),
				URL:  convert.ToPtr(link.URL),
				Meta: metaLinksBody,
			})
		}

		date := common.CurrentTimestamp()

		createdNews, err := s.dataProvider.NewsProvider().Insert(ctx, model.News{
			Media: &model.NewsMedia{
				Title:     convert.ToPtr(fmt.Sprintf("Digest hour: %d, Day: %d", date.Hour(), date.Day())),
				Text:      convert.ToPtr(digestResponse.content),
				Resources: resourcesList,
			},
			Source: convert.ToPtr("gpt-bing"),
			Status: convert.ToPtr(model.StatusPending),
		})
		if err != nil {
			return errors.Wrap(err, "failed to insert news digest")
		}

		coinsBatch, newsCoinsBatch := createCoinsNewsCoinsBatch(createdNews.ID, digestResponse.coins)
		if err = s.dataProvider.CoinsProvider().UpsertCoinsBatch(ctx, coinsBatch); err != nil {
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

		s.log.WithFields(logrus.Fields{
			"digest-hour": date.Hour(),
			"digest-day":  date.Day(),
		}).Debug("Finished generating")
		return nil
	})

	return nil
}

func (s service) readResponses(ctx context.Context, responsesChan <-chan chat_bot.ParsedFrame) (*generationsResponse, error) {
	response := &generationsResponse{}
	for msg := range responsesChan {
		if msg.Skip {
			continue
		}

		response.content = fmt.Sprintf("%s\n\n%s", msg.Text, strings.TrimPrefix(msg.AdaptiveCards, "\n"))
		response.links = msg.Links
	}

	if ctx.Err() != nil {
		return nil, errors.New("failed to ask bot due to cancelled context")
	}

	coinsSet := make(map[string]bool)
	for _, match := range coinsRegex.FindAllStringSubmatch(response.content, -1) {
		for _, coin := range strings.Split(match[1], ",") {
			coinsSet[strings.TrimSpace(coin)] = true
		}
	}

	for k := range coinsSet {
		response.coins = append(response.coins, k)
	}

	response.content = coinsRegex.ReplaceAllString(response.content, "")
	return response, nil
}
