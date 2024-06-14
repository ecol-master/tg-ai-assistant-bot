package bot

import (
	"errors"
	"fmt"
	"planpilot/internal/calendar_service"
	"planpilot/internal/config"
	"planpilot/internal/db"
	"planpilot/internal/logger"
	"planpilot/internal/openai"
	"planpilot/internal/scheduler"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// handler get calendar
func handlerGetCalendar(message *tgbotapi.Message, openaiAns string) (string, error) {
	task, err := openai.ParseGetCalendarTask(openaiAns)
	if err != nil {
		return "", errors.New("can not parse the GetCalendar task data")
	}
	calendarSvc, err := calendar_service.New()
	if err != nil {
		return "", errors.New("can not create the calendar_service")
	}

	scheduleDate, err := time.Parse(config.TIME_LAYOUT, task.Timestamp)
	if err != nil {
		logger.Error("error during parsing the date in handlerGetCalendar, err: ", err)
		return "", err
	}

	events, err := calendarSvc.GetEventsByDate(task.Timestamp)
	if err != nil {
		return "", errors.New("can not get the events from claendar")
	}
	if len(events.Items) == 0 {
		return fmt.Sprintf("You don't have any events on `%s`", formatDate(scheduleDate)), nil
	}

	eventsStr, err := calendarSvc.EventsToString(events)
	if err != nil {
		return eventsStr, err
	}

	return fmt.Sprintf("Your schedule on `%s` 📆:\n\n%s", formatDate(scheduleDate), eventsStr), nil
}

// handler help add event
func handlerHelpAddEvent(message *tgbotapi.Message, openaiAns string) (string, error) {
	calendarSvc, err := calendar_service.New()
	if err != nil {
		logger.Error("error during creating calendarService at handlerHelpAddEvent, err: ", err)
		return "I can not help you to add new event.", err
	}

	helpAddEventTask, err := openai.ParseHelpAddEventTask(openaiAns)
	if err != nil {
		logger.Error("error during parsing openaiAns in helpAddEvent handler, err: ", err)
		return "I can not help you to add new event.", err
	}
	events, err := calendarSvc.GetEventsByDate(helpAddEventTask.Timestamp)

	if err != nil {
		logger.Error("error during getting events by date in handlerHelpAddEvent, err: ", err)
		return "I can not help you to add new event.", err
	}
	if len(events.Items) == 0 {
		promptTemplate := "Answer me as you my personal assistant. I want to planning event `%s` at `%s`. Help me to add this in my schedule"
		prompt := fmt.Sprintf(promptTemplate, helpAddEventTask.Reminder, helpAddEventTask.Timestamp)
		return openai.MakeOpenAICall(prompt), nil
	}

	eventsStr, err := calendarSvc.EventsToString(events)
	if err != nil {
		logger.Error("error during converting events to string in handlerHelpAddEvent, err: ", err)
		return "I can not help you to add new event.", err
	}

	openaiPromptTemplate := "Answer me as you my personal assistant. I have a schedule: %s\n Help me to insert event `%s` in my schedule"
	prompt := fmt.Sprintf(openaiPromptTemplate, eventsStr, helpAddEventTask.Reminder)
	return openai.MakeOpenAICall(prompt), nil
}

// hanlder add event
func handlerAddEvent(message *tgbotapi.Message, openaiAns string) (string, error) {
	failedText := "I can not add this task to your schedule"

	calendarSvc, err := calendar_service.New()
	if err != nil {
		logger.Error("error during the creating new calendar_service, err", err)
		return failedText, err
	}

	event, err := calendarSvc.AddEvent(openaiAns)
	if err != nil {
		logger.Error("error during adding event to schedule, err", err)
		return failedText, err
	}

	events, err := calendarSvc.GetEventsByDate(event.Start.DateTime)
	if err != nil {
		logger.Error("error during getting events by date in handlerAddEvent, err: ", err)
		return failedText, err
	}

	eventsStr, err := calendarSvc.EventsToString(events)
	if err != nil {
		logger.Error("error during the converting events to string, err:", err)
		return failedText, err
	}

	scheduleDate, err := time.ParseInLocation(config.TIME_LAYOUT, event.Start.DateTime, time.UTC)
	if err != nil {
		logger.Error("error during parsing the date in handlerAddEvent, err: ", err)
		return eventsStr, err
	}

	logger.Info("LOCAL: ", time.Now().Local().UTC())
	result := fmt.Sprintf("I add new event `%s`.\nYour new schedule on `%s` 📆:\n\n%s", event.Summary, formatDate(scheduleDate), eventsStr)

	err = scheduler.AddTask(db.UserTask{
		UserID:    uint(message.From.ID),
		Schedule:  "",
		Timestamp: scheduleDate,
		Prompt:    event.Summary,
		Active:    true,
	}, bot)

	if err == nil {
		result += "\nI will remind you about these event."
	}

	return result, nil
}

// add unkonwn event
func handlerAddUnkownEvent(message *tgbotapi.Message, openaiAns string) (string, error) {
	return "This is uknown event:\n" + openaiAns, nil
}

func formatDate(datetime time.Time) string {
	const template = "January 2, 2006"
	return datetime.Format(template)
}
