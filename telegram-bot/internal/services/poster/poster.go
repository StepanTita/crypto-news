package poster

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	commonutils "common"
	commoncfg "common/config"

	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	"telegram-bot/internal/config"
	"telegram-bot/internal/utils"
)

type Poster interface {
	Post(ctx context.Context) (int, error)
}

type poster struct {
	templator commoncfg.Templator
	log       *logrus.Entry

	dataProvider store.DataProvider

	bot *tgbotapi.BotAPI
}

func New(cfg config.Config, bot *tgbotapi.BotAPI) Poster {
	return &poster{
		log:       cfg.Logging().WithField("service", "[TELEGRAM-POSTER]"),
		templator: cfg,

		dataProvider: store.New(cfg),

		bot: bot,
	}
}

func (p poster) Post(ctx context.Context) (int, error) {
	newsChannels, err := p.dataProvider.NewsChannelsProvider().Select(ctx)
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return 0, errors.Wrap(err, "failed to select pending news-channels")
		}
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
	for _, n := range news {
		for _, newsChannel := range newsChannelsMapping[n.ID] {
			msg := tgbotapi.NewMessage(newsChannel.ChannelID, "")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.Text = fmt.Sprintf(p.templator.Template(commonutils.NewsPost),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, convert.FromPtr(n.Media.Title)),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, convert.FromPtr(n.Media.Text)),
				tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, convert.FromPtr(n.Source)))
			if _, err := p.bot.Send(msg); err != nil {
				p.log.WithError(err).Error("failed to send post to bot API")
				continue
			}

			count++
		}
	}

	if err := p.dataProvider.NewsChannelsProvider().Remove(ctx, model.NewsChannel{}); err != nil {
		return count, errors.Wrap(err, "failed to remove entity")
	}

	return count, nil
}
