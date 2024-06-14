package bot

import (
	"fmt"
	"io"
	"os"
	"planpilot/internal/db"
	"planpilot/internal/logger"
	"planpilot/internal/openai"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// function to check user in db
func createUserIfNotExists(message *tgbotapi.Message) {
	_, err := db.GetUser(uint(message.From.ID))
	if err != nil {
		_, err = db.CreateUser(uint(message.From.ID))

		if err != nil {
			logger.Error("error during creating new user", err)
		} else {
			logger.Info("user was create successfully")
		}
	}
}

// handler for commang /start
func startHandler(message *tgbotapi.Message) (string, error) {
	createUserIfNotExists(message)
	greetingText := "Hi, %s!\n" +
		"I am your personal assistant, you can ask me about everything.\n\n" +
		"Moreover, i can to make summary from your email box. Ask me like this:\n`Check the last 10 emails in my mail box.`\n\n" +
		"Also, i can add new events in your schedule (please send me date and time).\n`Add watching football with a friend to my schedule on June 1 at 21:45`"

	return fmt.Sprintf(greetingText, message.From.UserName), nil
}

// handler for command /help
func helpHandler(message *tgbotapi.Message) string {
	createUserIfNotExists(message)

	helpText := `
	You can also paste any text and I'll send it to a caht GPT and show you the result
	`
	return helpText
}

func userTaskHandler(message *tgbotapi.Message) (string, error) {
	userInput := message.Text
	prompt := parsePrompt()

	result := openai.MakeOpenAICall(fmt.Sprintf("%s\nUser input: %s", string(prompt), userInput))

	if strings.Contains(result, "EMAIL_SCHEDULED_TASK") {
		return emailsScheduledTaskHandler(message, result)
	} else if strings.Contains(result, "PLANNED_TASK") {
		return plannedTaskHandler(message, result)
	} else if strings.Contains(result, "SCHEDULED_TASK") {
		return scheduledTaskHandler(message, result)
	} else if strings.Contains(result, "EMAIL_TASK") {
		return emailsTaskHandler(message, result)
	} else if strings.Contains(result, "RESPONSE") {
		return openai.MakeOpenAICall(userInput), nil
	} else if strings.Contains(result, "GET_CALENDAR") {
		return handlerGetCalendar(message, result)
	} else if strings.Contains(result, "HELP_ADD_EVENT") {
		return handlerHelpAddEvent(message, result)
	} else if strings.Contains(result, "ADD_EVENT") {
		return handlerAddEvent(message, result)
		// } else if strings.Contains(result, "ADD_UKNOWN_EVENT") {
		// return handlerAddUnkownEvent(message, result)
	} else {
		return openai.MakeOpenAICall(userInput), nil
	}
}

func parsePrompt() []byte {
	file, err := os.Open("../prompts/determine")

	if err != nil {
		logger.Error("Error open prompt file: ", err)
		return []byte("")
	}

	prompt, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Error reading prompt file: ", err)
		return []byte("")
	}
	return prompt
}
