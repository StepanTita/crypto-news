package poster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	commonutils "common"
	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"common/iteration"
	"common/locale"
	"common/transform"
	"telegram-bot/internal/config"
)

const telegramMaxMessageLen = 4096

type Poster interface {
	Post(ctx context.Context) (int, error)
}

type poster struct {
	cfg config.Config

	log *logrus.Entry

	dataProvider store.DataProvider

	bot *tgbotapi.BotAPI
}

func New(cfg config.Config, bot *tgbotapi.BotAPI) Poster {
	return &poster{
		cfg: cfg,
		log: cfg.Logging().WithField("service", "[TELEGRAM-POSTER]"),

		dataProvider: store.New(cfg),

		bot: bot,
	}
}

func (p poster) Post(ctx context.Context) (int, error) {
	newsChannels, err := p.dataProvider.NewsChannelsProvider().BySources(p.cfg.Sources()).Select(ctx)
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return 0, errors.Wrap(err, "failed to select pending news-channels")
		}
	}

	if len(newsChannels) == 0 {
		p.log.Debug("No news-channels found. Skipping...")
		return 0, nil
	}

	newsIDs := make([]uuid.UUID, len(newsChannels))

	// newsID : []newsChannels
	newsChannelsMapping := make(map[uuid.UUID][]model.NewsChannel)
	for i, newsChannel := range newsChannels {
		newsIDs[i] = newsChannel.NewsID
		newsChannelsMapping[newsChannel.NewsID] = append(newsChannelsMapping[newsChannel.NewsID], newsChannel)
	}

	newsIDs = iteration.Unique(newsIDs)
	news, err := p.dataProvider.NewsProvider().ByIDs(newsIDs).Select(ctx)
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return 0, errors.Wrap(err, "failed to select pending news")
		}
	}

	count := 0
	successfulIDs := make([]uuid.UUID, 0, 10)
	for _, n := range news {
		// TODO: refactor to get batch instead of querying db in loop
		coins, err := p.dataProvider.CoinsProvider().ByNewsID(n.ID).Select(ctx)
		if err != nil {
			if !errors.Is(err, data.ErrNotFound) {
				p.log.WithError(err).Error("failed to get coins")
			}
		}

		for _, newsChannel := range newsChannelsMapping[n.ID] {
			msg, media, err := p.buildMessage(newsChannel.ChannelID, n, coins)
			if err != nil {
				p.log.WithError(err).Error("failed to build message")
				continue
			}

			if len(msg.Text) > telegramMaxMessageLen {
				if err := p.sendMultiple(msg); err != nil {
					p.log.WithError(err).Error("failed to send multiple posts to bot API")
					continue
				}
			} else {
				if _, err := p.bot.Send(msg); err != nil {
					p.log.WithError(err).Error("failed to send post to bot API")
					continue
				}
			}

			if media != nil {
				if m, ok := media.(tgbotapi.PhotoConfig); ok {
					if _, err := p.bot.Send(m); err != nil {
						p.log.WithError(err).Error("failed to send media to bot API")
						continue
					}
				} else if m, ok := media.(tgbotapi.MediaGroupConfig); ok {
					if _, err := p.bot.SendMediaGroup(m); err != nil {
						p.log.WithError(err).Error("failed to send media to bot API")
						continue
					}
				}
			}

			successfulIDs = append(successfulIDs, newsChannel.ID)
			count++
		}
	}

	if err := p.dataProvider.NewsChannelsProvider().BySources(p.cfg.Sources()).ByIDs(successfulIDs).Remove(ctx, model.NewsChannel{}); err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return count, errors.Wrap(err, "failed to remove entity")
		}
	}

	return count, nil
}

func (p poster) buildMessage(channelID int64, news model.News, coins []model.Coin) (*tgbotapi.MessageConfig, tgbotapi.Chattable, error) {
	msg := tgbotapi.NewMessage(channelID, "")
	msg.ParseMode = tgbotapi.ModeHTML

	references := strings.Builder{}

	body := convert.FromPtr(news.Media.Text)

	images := make([]any, 0, 5)
	for _, resource := range news.Media.Resources {
		if convert.FromPtr(resource.Type) == model.ResourceTypeSource {
			if resource.Meta == nil {
				continue
			}
			var metaLinks model.MetaLinksData
			if err := json.Unmarshal(resource.Meta, &metaLinks); err != nil {
				return nil, nil, errors.Wrap(err, "failed to unmarshal media meta")
			}
			body = strings.ReplaceAll(body, fmt.Sprintf("[^%s^][%s]", metaLinks.ID, metaLinks.ID),
				fmt.Sprintf("<a href=\"%s\">[%s]</a>", metaLinks.URL, metaLinks.ID))
			references.WriteString(fmt.Sprintf("[%s] <a href=\"%s\">%s</a>.\n", metaLinks.ID, metaLinks.URL, metaLinks.Title))
		} else if convert.FromPtr(resource.Type) == model.ResourceTypeImage {
			image, err := downloadFile(convert.FromPtr(resource.URL))
			if err != nil {
				p.log.WithError(err).Errorf("failed to read image from url: %s", convert.FromPtr(resource.URL))
				continue
			}
			images = append(images, tgbotapi.NewInputMediaPhoto(tgbotapi.FileBytes{Bytes: image, Name: uuid.NewString()}))
		}
	}

	coinsHashTags := strings.Builder{}
	for _, coin := range coins {
		coinsHashTags.WriteString(fmt.Sprintf("#%s ", coin.Code))
	}

	rawTemplate := p.cfg.Template(commonutils.NewsPost)

	msg.Text = transform.CleanUnsupportedHTML(fmt.Sprintf(locale.PrepareTemplate(p.cfg, rawTemplate, convert.FromPtr(news.Locale)),
		escapeKeepingHTML(convert.FromPtr(news.Media.Title)),
		escapeKeepingHTML(body),
		escapeKeepingHTML(references.String()),
		escapeKeepingHTML(coinsHashTags.String()),
		escapeKeepingHTML(convert.FromPtr(news.Source)),
	))

	if len(images) == 1 {
		return &msg, tgbotapi.NewPhoto(channelID, images[0].(tgbotapi.InputMediaPhoto).Media), nil
	} else if len(images) > 1 {
		return &msg, tgbotapi.NewMediaGroup(channelID, images), nil
	}

	return &msg, nil, nil
}

func (p poster) sendMultiple(msg *tgbotapi.MessageConfig) error {
	partText := ""
	parts := strings.Split(msg.Text, "\n")

	for _, part := range parts {
		if len(partText+part) < telegramMaxMessageLen {
			if partText == "" {
				partText = part
			} else {
				partText += fmt.Sprintf("\n%s", part)
			}
		} else {
			msg.Text = partText
			if _, err := p.bot.Send(msg); err != nil {
				return errors.Wrap(err, "failed to send message")
			}
			partText = ""
		}
	}
	if partText != "" {
		msg.Text = partText
		if _, err := p.bot.Send(msg); err != nil {
			return errors.Wrap(err, "failed to send message")
		}
	}

	return nil
}
