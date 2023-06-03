package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	chat_bot "github.com/StepanTita/go-EdgeGPT/chat-bot"
	"github.com/google/uuid"
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

	common.RunEveryWithBackoff(1*time.Hour, 15*time.Second, 15*time.Minute, func() error {
		s.log.Debug("Generating digest...")

		err := bot.Init(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to initialize bot")
		}

		parsedResponsesChan, err := bot.Ask(ctx, s.cfg.GPTConfig().InitialPrompt(), s.cfg.GPTConfig().Style(), true)
		if err != nil {
			return errors.Wrap(err, "failed to ask bot")
		}

		digestResponse := s.readResponses(parsedResponsesChan)
		if err != nil {
			return errors.Wrap(err, "failed to generate digest")
		}

		resourcesList := make([]model.NewsMediaResource, 0, len(digestResponse.linksSet))
		for link := range digestResponse.linksSet {
			resourcesList = append(resourcesList, model.NewsMediaResource{
				Type: convert.ToPtr(model.ResourceTypeSource),
				URL:  convert.ToPtr(link),
				Meta: nil,
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

func (s service) readResponses(responsesChan <-chan chat_bot.ParsedFrame) generationsResponse {
	response := generationsResponse{
		linksSet: make(map[string]bool),
	}
	for msg := range responsesChan {
		if msg.Skip {
			continue
		}

		response.content = fmt.Sprintf("%s\n\n%s", msg.Text, strings.TrimPrefix(msg.AdaptiveCards, "\n"))
		for _, link := range msg.Links {
			response.linksSet[link] = true
		}
	}
	re := regexp.MustCompile(`\<coins\>\[([A-Z\,]+)\]\<\/coins\>`)
	if match := re.FindStringSubmatch(response.content); len(match) > 0 {
		response.coins = strings.Split(match[1], ",")
	}
	return response
}

func toNewsChannelsBatch(news *model.News, channels []model.Channel) []model.NewsChannel {
	newsChannels := make([]model.NewsChannel, len(channels))
	for i, c := range channels {
		newsChannels[i] = model.NewsChannel{
			ChannelID: c.ChannelID,
			NewsID:    news.ID,
		}
	}
	return newsChannels
}

func createCoinsNewsCoinsBatch(newsID uuid.UUID, codes []string) ([]model.Coin, []model.NewsCoin) {
	coins := make([]model.Coin, len(codes))
	newsCoins := make([]model.NewsCoin, len(codes))
	for i, code := range codes {
		newsCoins[i] = model.NewsCoin{
			Code:   code,
			NewsID: newsID,
		}

		coins[i] = model.Coin{
			Code: code,
			Slug: code,
		}
	}
	return coins, newsCoins
}
