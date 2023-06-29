package handler

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	commoncfg "common/config"
	"common/data"
	"common/data/model"
	"common/data/store"
	commonerrors "common/errors"
	"telegram-bot/internal/config"
	"telegram-bot/internal/utils"
)

type Handler interface {
	HandleCommand(ctx context.Context, incomingMsg *tgbotapi.Message) (*tgbotapi.MessageConfig, error)
}

type handler struct {
	log       *logrus.Entry
	templator commoncfg.Templator

	dataProvider store.DataProvider

	interactions map[int64]utils.Command
}

func New(cfg config.Config) Handler {
	return &handler{
		log:       cfg.Logging().WithField("service", "[HANDLER]"),
		templator: cfg,

		dataProvider: store.New(cfg),

		interactions: make(map[int64]utils.Command),
	}
}

func (h handler) HandleCommand(ctx context.Context, incomingMsg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(incomingMsg.Chat.ID, "")

	command := utils.Command(incomingMsg.Command())

	if inter, ok := h.interactions[incomingMsg.Chat.ID]; ok {
		command = inter

		msg, err := h.handleInteraction(ctx, inter, incomingMsg)
		if err != nil {
			return msg, errors.Wrap(err, "failed to handle interaction")
		}

		delete(h.interactions, incomingMsg.Chat.ID)
		return msg, nil
	}

	switch command {
	case utils.StartCommand:
		msg.Text = fmt.Sprintf(h.templator.Template(utils.StartCommand.String()), incomingMsg.From.UserName)
	case utils.SubscribeCommand:
		if incomingMsg.From == nil {
			h.interactions[incomingMsg.Chat.ID] = command
			msg.Text = "Please, provide a one time token."
			return &msg, nil
		}

		msg.Text = fmt.Sprintf(h.templator.Template(utils.SubscribeCommand.String()))
		if err := h.handleSubscribeUsername(ctx, incomingMsg.From.UserName, incomingMsg.Chat.ID); err != nil {
			if !errors.Is(err, data.ErrDuplicateRecord) {
				return &msg, errors.Wrap(err, "failed to subscribe channel")
			}
			msg.Text = "This channel was already registered!"
		}
	default:
		msg.Text = "Sorry... I don't know that command, yet!"
	}

	return &msg, nil
}

func (h handler) handleSubscribeUsername(ctx context.Context, username string, channelID int64) error {
	_, err := h.dataProvider.WhitelistProvider().ByUsername(username).Get(ctx)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return commonerrors.ErrAccessDenied
		}
		return errors.Wrap(err, "failed to get whitelist record by username")
	}

	return h.handleSubscribe(ctx, channelID)
}

func (h handler) handleSubscribeToken(ctx context.Context, token string, channelID int64) error {
	err := h.dataProvider.WhitelistProvider().ExtractToken(ctx, uuid.MustParse(token))
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return commonerrors.ErrAccessDenied
		}
		return errors.Wrap(err, "failed to get whitelist record by token")
	}

	return h.handleSubscribe(ctx, channelID)
}

func (h handler) handleSubscribe(ctx context.Context, channelID int64) error {
	_, err := h.dataProvider.ChannelsProvider().Insert(ctx, model.Channel{ChannelID: channelID})
	if err != nil {
		return errors.Wrap(err, "failed to insert new channel id")
	}
	return nil
}

func (h handler) handleInteraction(ctx context.Context, command utils.Command, incomingMsg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(incomingMsg.Chat.ID, "")
	switch command {
	case utils.SubscribeCommand:
		token := incomingMsg.Text
		msg.Text = fmt.Sprintf("Subscribed: %s", token)
		if err := h.handleSubscribeToken(ctx, token, incomingMsg.Chat.ID); err != nil {
			if !errors.Is(err, data.ErrDuplicateRecord) {
				return &msg, errors.Wrap(err, "failed to handle whitelist")
			} else {
				msg.Text = "This user was already whitelisted!"
			}
		}
	default:
		msg.Text = "Sorry... I don't know that command, yet!"
	}
	return &msg, nil
}
