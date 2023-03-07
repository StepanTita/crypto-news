package poster

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/convert"
	"common/data/model"
	"common/data/store"
	"telegram-bot/internal/config"
	"telegram-bot/internal/utils"
)

type Poster interface {
	Post(ctx context.Context) (int, error)
}

type poster struct {
	templator config.Templator
	log       *logrus.Entry

	dataProvider store.DataProvider

	bot *tgbotapi.BotAPI
}

func New(cfg config.Config, bot *tgbotapi.BotAPI) Poster {
	return &poster{
		log:       cfg.Logging().WithField("service", "[POSTER]"),
		templator: cfg,

		dataProvider: store.New(cfg),

		bot: bot,
	}
}

func (p poster) Post(ctx context.Context) (int, error) {

	// TODO: remove processed, instead of status change
	newsChannels, err := p.dataProvider.NewsChannelsProvider().ByStatus(ctx, model.StatusPending).Select(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to select pending news-channels")
	}

	newsIDs := make([]uuid.UUID, len(newsChannels))

	newsChannelsMapping := make(map[uuid.UUID][]model.NewsChannel)
	for i, newsChannel := range newsChannels {
		newsIDs[i] = newsChannel.NewsID
		newsChannelsMapping[newsChannel.NewsID] = append(newsChannelsMapping[newsChannel.NewsID], newsChannel)
	}

	newsIDs = utils.Unique(newsIDs)
	news, err := p.dataProvider.NewsProvider().ByIDs(ctx, newsIDs).Select(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to select pending news")
	}

	count := 0
	for _, n := range news {
		for _, newsChannel := range newsChannelsMapping[n.ID] {
			msg := tgbotapi.NewMessage(newsChannel.ChannelID, "")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			msg.Text = fmt.Sprintf(p.templator.Template(utils.NewsPost),
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

	//_, err = p.dataProvider.NewsChannelsProvider().ByStatus(ctx, model.StatusPending).Update(ctx, model.UpdateNewsChannelParams{
	//	Status: convert.ToPtr(model.StatusProcessed),
	//})
	//if err != nil {
	//	return count, errors.Wrap(err, "failed to update news status to processed")
	//}

	return count, nil
}
