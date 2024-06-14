package bot

import (
	"log"
	"planpilot/internal/config"
	"planpilot/internal/logger"
	"planpilot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
)

func init() {
	b, err := tgbotapi.NewBotAPI(config.New().TELEGRAM_TOKEN)
	if err != nil {
		log.Panic("Error during start the bot: ", err)
	}
	bot = b
}

func Run() error {

	// true, if you want to see all the requests in terminal
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	scheduler.Start()

	var (
		answer string
		err    error
	)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			switch update.Message.Command() {
			case "start":
				answer, err = startHandler(update.Message)
				if err != nil {
					logger.Error("error in /start handler", err)
				}
			case "help":
				answer = helpHandler(update.Message)
			case "status":
				answer = "I'm ok."
			default:
				answer, err = userTaskHandler(update.Message)
				if err != nil {
					logger.Error("error in scheduledTask, err: ", err)
				}
			}
			msg.Text = answer
			msg.ParseMode = "Markdown"
			logger.Info("bot answer: ", answer)
			bot.Send(msg)
		}
	}
	return nil
}
