package bot

import (
	"errors"
	"fmt"
	"planpilot/internal/db"
	"planpilot/internal/logger"
	"planpilot/internal/scheduler"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type scheduledTask struct {
	Schedule string
	Prompt   string
}

func scheduledTaskHandler(message *tgbotapi.Message, openaiAns string) (string, error) {
	logger.Info("openai answer on scheduled task: \n", openaiAns)

	task, err := parseScheduleTask(openaiAns)
	if err != nil {
		logger.Error("error during extracting schedule task: ", err)
		return "I can not to extracting your scheduled task", nil
	}

	scheduler.AddTask(db.UserTask{
		UserID:    uint(message.From.ID),
		Prompt:    task.Prompt,
		Timestamp: time.Now().Add(time.Minute),
		Schedule:  task.Schedule,
	}, bot)

	return fmt.Sprintf("I will remind you about: `%s` ", task.Prompt), nil
}

// ```
// SCHEDULED_TASK
// cron: "30 9 * * *"
// prompt: "Reminder: Say good morning to your mom."
// ```

func parseScheduleTask(openaiAns string) (*scheduledTask, error) {
	var task scheduledTask
	lines := strings.Split(strings.Trim(openaiAns, " \n\r"), "\n")
	if len(lines) != 5 {
		return &task, errors.New("scheculed task answer is not, correct:\n" + openaiAns)
	}
	lines = lines[1:4]

	cron_data := strings.Split(strings.Trim(lines[1], " \n\r"), ": ")
	if len(cron_data) != 2 {
		return &task, errors.New("cron string format is not correct" + lines[1])
	}

	task.Schedule = strings.Trim(cron_data[1], `"`)

	prompt_data := strings.SplitN(strings.Trim(lines[2], " \n\r"), ": ", 2)
	if len(prompt_data) != 2 {
		return &task, errors.New("reminder string format is not correct" + lines[2])
	}

	task.Prompt = strings.Trim(prompt_data[1], `"`)
	return &task, nil
}
