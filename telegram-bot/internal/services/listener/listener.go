package listener

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	commoncfg "common/config"
	"common/convert"
	"common/data/store"
	commonerrors "common/errors"
	"telegram-bot/internal/config"
	"telegram-bot/internal/services/handler"
	"telegram-bot/internal/utils"
)

type Listener interface {
	Listen(ctx context.Context) error
}

type listener struct {
	cfg config.Config
	log *logrus.Entry

	handler handler.Handler

	dataProvider store.DataProvider

	bot *tgbotapi.BotAPI
}

func New(cfg config.Config, bot *tgbotapi.BotAPI) Listener {
	return &listener{
		cfg: cfg,
		log: cfg.Logging().WithField("service", "[LISTENER]"),

		handler:      handler.New(cfg),
		dataProvider: store.New(cfg),

		bot: bot,
	}
}

func (l listener) Listen(ctx context.Context) error {
	l.log.Infof("Authorized on account %s", l.bot.Self.UserName)

	buttonsConfig, updatesChan, err := l.configUpdates(l.bot)
	if err != nil {
		return errors.Wrap(err, "failed to configure bot")
	}

	for update := range updatesChan {
		log := l.log.WithField("update_id", update.UpdateID)
		log.Debug("reading updates...")

		var msg *tgbotapi.MessageConfig

		if update.FromChat().IsGroup() || update.FromChat().IsSuperGroup() || update.FromChat().IsPrivate() && update.Message != nil {
			log.Debugf("update from chat: %d, with message: %s", update.Message.Chat.ID, update.Message.Text)

			if update.Message.IsCommand() {
				log.Debug("handling command...")
				msg, err = l.handler.HandleCommand(ctx, update.Message)
				if err != nil {
					if !errors.Is(err, commonerrors.ErrAccessDenied) {
						errRef := uuid.NewString()
						log.WithField("error-ref", errRef).WithError(err).Error("failed to handle")
						msg.Text = fmt.Sprintf("Please retry. Something went wrong...\nError reference is: %s", errRef)
					} else {
						msg.Text = "You are not allowed to perform this action! Please refer to @Vladyslavpv for information."
					}
				}
				log.Debugf("command handled with output: %s", msg.Text)

				msg.ReplyToMessageID = update.Message.MessageID
				msg.ReplyMarkup = buttonsConfig
			}
		} else if update.FromChat().IsChannel() && update.ChannelPost != nil {
			log.Debugf("update from channel: %d, with message: %s", update.ChannelPost.Chat.ID, update.ChannelPost.Text)

			log.Debug("handling command...")
			msg, err = l.handler.HandleCommand(ctx, update.ChannelPost)
			if err != nil {
				if !errors.Is(err, commonerrors.ErrAccessDenied) {
					errRef := uuid.NewString()
					log.WithField("error-ref", errRef).WithError(err).Error("failed to handle")
					msg.Text = fmt.Sprintf("Please retry. Something went wrong...\nError reference is: %s", errRef)
				} else {
					msg.Text = "You are not allowed to perform this action! Please refer to @Vladyslavpv for information."
				}
			}
			log.Debugf("command handled with output: %s", msg.Text)
		}

		if msg == nil {
			continue
		}

		if _, err := l.bot.Send(msg); err != nil {
			log.WithError(err).Error("failed to send message to bot API")
		}
	}
	return nil
}

func (l listener) configUpdates(bot *tgbotapi.BotAPI) (tgbotapi.ReplyKeyboardMarkup, tgbotapi.UpdatesChannel, error) {
	if err := tgbotapi.SetLogger(l.log.WithField("[BOT]", bot.Self.UserName)); err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, nil, errors.Wrap(err, "failed to set bog logger")
	}
	bot.Debug = slices.Contains([]string{commoncfg.EnvironmentDev, commoncfg.EnvironmentLocal, commoncfg.EnvironmentStaging}, l.cfg.Environment())

	commandsConfig := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     utils.StartCommand.String(),
			Description: "start the interaction",
		},
		tgbotapi.BotCommand{
			Command:     utils.SubscribeCommand.String(),
			Description: "subscribe current channel to the news",
		},
	)

	commandsConfig.Scope = convert.ToPtr(tgbotapi.NewBotCommandScopeAllChatAdministrators())

	buttonsConfig := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.StartCommand.Command()),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(utils.SubscribeCommand.Command()),
		),
	)

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	_, err := bot.Request(commandsConfig)
	if err != nil {
		return tgbotapi.ReplyKeyboardMarkup{}, nil, errors.Wrap(err, "failed to send commands config")
	}

	// Start polling Telegram for updates.
	return buttonsConfig, bot.GetUpdatesChan(updateConfig), err
}
