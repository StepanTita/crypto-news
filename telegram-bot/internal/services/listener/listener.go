package listener

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	commoncfg "common/config"
	"common/data/store"
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
		if update.Message == nil {
			continue
		}

		log.Debugf("update from chat: %d, with message: %s", update.Message.Chat.ID, update.Message.Text)

		var msg *tgbotapi.MessageConfig

		if update.Message.IsCommand() {
			log.Debug("handling command...")
			msg, err = l.handler.HandleCommand(ctx, update.Message)
			if err != nil {
				return errors.Wrap(err, "failed to handle command")
			}
			log.Debugf("command handled with output: %s", msg.Text)
		}

		if msg == nil {
			continue
		}

		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = buttonsConfig

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
