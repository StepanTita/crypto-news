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
	"telegram-bot/internal/config"
	"telegram-bot/internal/utils"
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

	newsIDs = utils.Unique(newsIDs)
	news, err := p.dataProvider.NewsProvider().ByIDs(newsIDs).Select(ctx)
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return 0, errors.Wrap(err, "failed to select pending news")
		}
	}

	count := 0
	successfulIDs := make([]uuid.UUID, 0, 10)
	for _, n := range news {
		for _, newsChannel := range newsChannelsMapping[n.ID] {
			msg, err := p.buildMessage(newsChannel.ChannelID, n)
			if err != nil {
				p.log.WithError(err).Error("failed to build message")
				continue
			}

			// TODO: use MessageEntity
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

			successfulIDs = append(successfulIDs, newsChannel.ID)
			count++
		}
	}

	if err := p.dataProvider.NewsChannelsProvider().BySources(p.cfg.Sources()).ByIDs(successfulIDs).Remove(ctx, model.NewsChannel{}); err != nil {
		return count, errors.Wrap(err, "failed to remove entity")
	}

	return count, nil
}

func (p poster) buildMessage(channelID int64, news model.News) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(channelID, "")
	msg.ParseMode = tgbotapi.ModeHTML

	references := strings.Builder{}

	for _, resource := range news.Media.Resources {
		if resource.Meta == nil {
			continue
		}
		var metaLinks model.MetaLinksData
		if err := json.Unmarshal(resource.Meta, &metaLinks); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal media meta")
		}
		references.WriteString(fmt.Sprintf("%s <a href=\"%s\">%s</a>.\n", metaLinks.ID, metaLinks.URL, metaLinks.Title))
	}

	msg.Text = fmt.Sprintf(p.cfg.Template(commonutils.NewsPost),
		escapeKeepingHTML(convert.FromPtr(news.Media.Title)),
		escapeKeepingHTML(convert.FromPtr(news.Media.Text)),
		escapeKeepingHTML(references.String()),
		escapeKeepingHTML(convert.FromPtr(news.Source)),
	)

	return &msg, nil
}

func (p poster) sendMultiple(msg *tgbotapi.MessageConfig) error {
	partText := ""
	parts := strings.Split(msg.Text, "\n\n")

	for _, part := range parts {
		if len(partText+part) < telegramMaxMessageLen {
			if partText == "" {
				partText = part
			} else {
				partText += fmt.Sprintf("\n\n%s", part)
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

func escapeKeepingHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
	)

	return replacer.Replace(text)
}
