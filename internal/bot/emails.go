package bot

import (
	"planpilot/internal/db"
	"planpilot/internal/gmail_service"
	"planpilot/internal/logger"
	"planpilot/internal/openai"
	"planpilot/internal/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func emailsScheduledTaskHandler(message *tgbotapi.Message, openaiAns string) (string, error) {
	task, err := openai.ParseEmailScheduledTask(openaiAns)
	if err != nil {
		logger.Error("can not parse the emailScheduledTask, err: ", err)
		return "I can not understand your task", err
	}
	scheduler.AddEmailsTask(db.UserTask{
		UserID:   uint(message.From.ID),
		Schedule: task.Cron,
		Prompt:   openaiAns,
		Active:   true,
	}, bot)

	return "I add to planning sending emails", nil
}

// one-time email task
func emailsTaskHandler(message *tgbotapi.Message, openaiAns string) (string, error) {
	task, err := openai.ParseEmailsTask(openaiAns)
	if err != nil {
		logger.Error("error during parsing emailsTask, err: ", err)
		return "I can not to get compress emails from your box", err
	}

	gmailService, err := gmail_service.New()
	if err != nil {
		logger.Error("error during creating gmailService: ", err)
		return "I can not to get compress emails from your box", err
	}

	compressedEmailsText, err := gmailService.MakeCompressedEmailsText(task.EmailsCount)
	if err != nil {
		logger.Error("error during getting compressed text: ", err)
		return "I can not to get compress emails from your box", err
	}

	return compressedEmailsText, nil
}
