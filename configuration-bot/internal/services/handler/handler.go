package handler

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	commoncfg "common/config"
	"common/convert"
	"common/data"
	"common/data/model"
	"common/data/store"
	commonerrors "common/errors"
	"configuration-bot/internal/config"
	"configuration-bot/internal/utils"
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

	msg := tgbotapi.NewMessage(incomingMsg.Chat.ID, "")

	user, err := h.dataProvider.UsersProvider().ByUsername(incomingMsg.From.UserName).Get(ctx)
	if err != nil {
		return &msg, errors.Wrapf(err, "failed to get user by username: %s", incomingMsg.From.UserName)
	}

	if convert.FromPtr(user.Role) != model.RoleAdmin {
		return &msg, commonerrors.ErrAccessDenied
	}

	switch command {
	case utils.StartCommand:
		msg.Text = fmt.Sprintf("Hi, %s", incomingMsg.From.UserName)
	case utils.WhitelistCommand:
		h.interactions[incomingMsg.Chat.ID] = command
		msg.Text = "Please, provide username of a user to whitelist"
	case utils.UnlistCommand:
		h.interactions[incomingMsg.Chat.ID] = command
		msg.Text = "Please, provide username of a user to unlist"
	case utils.TokenCommand:
		token := uuid.New()
		msg.Text = fmt.Sprintf("One time Token: %s", token.String())
		if err := h.handleToken(ctx, token); err != nil {
			return &msg, errors.Wrap(err, "failed to generate one time token")
		}

	default:
		msg.Text = "Sorry... I don't know that command, yet!"
	}

	return &msg, nil
}

func (h handler) handleWhitelist(ctx context.Context, incomingMsg *tgbotapi.Message) error {
	_, err := h.dataProvider.UsersProvider().Insert(ctx, model.User{
		Username: convert.ToPtr(incomingMsg.Text),
		Platform: convert.ToPtr("telegram"),
		Role:     convert.ToPtr(model.RoleReader),
	})
	if err != nil {
		if !errors.Is(err, data.ErrDuplicateRecord) {
			return errors.Wrap(err, "failed to insert new user record")
		}
	}

	_, err = h.dataProvider.WhitelistProvider().Insert(ctx, model.Whitelist{Username: convert.ToPtr(incomingMsg.Text)})
	if err != nil {
		return errors.Wrap(err, "failed to insert new whitelist record")
	}
	return nil
}

func (h handler) handleUnlist(ctx context.Context, username string) error {
	err := h.dataProvider.WhitelistProvider().ByUsername(username).Remove(ctx, model.Whitelist{})
	if err != nil {
		if !errors.Is(err, data.ErrNotFound) {
			return errors.Wrapf(err, "failed to remove user from whitelist: %s", username)
		}
	}
	return nil
}

func (h handler) handleInteraction(ctx context.Context, command utils.Command, incomingMsg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(incomingMsg.Chat.ID, "")
	switch command {
	case utils.WhitelistCommand:
		username := incomingMsg.Text
		msg.Text = fmt.Sprintf("Whitelisted: %s", username)
		if err := h.handleWhitelist(ctx, incomingMsg); err != nil {
			if !errors.Is(err, data.ErrDuplicateRecord) {
				return &msg, errors.Wrap(err, "failed to handle whitelist")
			} else {
				msg.Text = "This user was already whitelisted!"
			}
		}
	case utils.UnlistCommand:
		username := incomingMsg.Text
		msg.Text = fmt.Sprintf("Unlisted: %s", username)
		if err := h.handleUnlist(ctx, username); err != nil {
			return &msg, errors.Wrap(err, "failed to handle unlist")
		}
	default:
		msg.Text = "Sorry... I don't know that command, yet!"
	}
	return &msg, nil
}

func (h handler) handleToken(ctx context.Context, token uuid.UUID) error {
	_, err := h.dataProvider.WhitelistProvider().Insert(ctx, model.Whitelist{Token: convert.ToPtr(token)})
	if err != nil {
		return errors.Wrap(err, "failed to insert new whitelist record")
	}
	return nil
}
