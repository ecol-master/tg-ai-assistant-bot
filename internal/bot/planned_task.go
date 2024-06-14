package bot

import (
	"planpilot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func plannedTaskHandler(message *tgbotapi.Message, openaiAns string) (string, error) {
	logger.Info("openai result to planned task:\n", openaiAns)

	return "I'm planned task", nil
}
