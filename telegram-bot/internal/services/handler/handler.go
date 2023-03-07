package handler

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"common/data/model"
	"common/data/store"
	"telegram-bot/internal/config"
	"telegram-bot/internal/utils"
)

type Handler interface {
	HandleCommand(ctx context.Context, incomingMsg *tgbotapi.Message) (*tgbotapi.MessageConfig, error)
}

type handler struct {
	log       *logrus.Entry
	templator config.Templator

	dataProvider store.DataProvider
}

func New(cfg config.Config) Handler {
	return &handler{
		log:       cfg.Logging().WithField("service", "[HANDLER]"),
		templator: cfg,

		dataProvider: store.New(cfg),
	}
}

func (h handler) HandleCommand(ctx context.Context, incomingMsg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(incomingMsg.Chat.ID, "")

	switch incomingMsg.Command() {
	case utils.StartCommand:
		msg.Text = fmt.Sprintf(h.templator.Template(utils.StartCommand), incomingMsg.From.UserName)
	case utils.SubscribeCommand:
		msg.Text = fmt.Sprintf(h.templator.Template(utils.SubscribeCommand))
		if err := h.handleSubscribe(ctx, incomingMsg.Chat.ID); err != nil {
			return nil, errors.Wrap(err, "failed to subscribe channel")
		}
	default:
		msg.Text = "Sorry... I don't know that command, yet!"
	}

	return &msg, nil
}

func (h handler) handleSubscribe(ctx context.Context, channelID int64) error {
	_, err := h.dataProvider.ChannelsProvider().Insert(ctx, model.Channel{ChannelID: channelID})
	if err != nil {
		return errors.Wrap(err, "failed to insert new channel id")
	}
	return nil
}
