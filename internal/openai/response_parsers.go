package openai

import (
	"errors"
	"planpilot/internal/logger"
	"strconv"
	"strings"
	"time"
)

// Здесь находятся структуры, которые описывают ответы от GPT.
// Структуры полностью копипастят структуру ответов из файла prompts/determine

type EmailScheduledTask struct {
	Cron        string
	EmailsCount int
}

type EmailsTask struct {
	EmailsCount int
}

type ScheduledTask struct {
	Cron         string
	ReminderText string
}

type PlannedTask struct {
	Timestamp time.Time
	Reminder  string
}

type GetCalendarTask struct {
	Timestamp string
}

type HelpAddEventTask struct {
	Timestamp string
	Reminder  string
}

type AddEventTask struct {
	TimestampStart string
	TimestampEnd   string
	Reminder       string
}

func ParseEmailScheduledTask(openaiResponse string) (*EmailScheduledTask, error) {
	var task EmailScheduledTask

	lines := strings.Split(strings.Trim(openaiResponse, " \n\r"), "\n")
	if len(lines) != 5 {
		return &task, errors.New("openaiAns is not correct, ans: " + openaiResponse)
	}

	lines = lines[2:4]
	cron_data := strings.Split(lines[0], ": ")
	if len(cron_data) != 2 {
		return &task, errors.New("cron_data is not correct, cron_data: " + lines[0])
	}
	task.Cron = strings.Trim(cron_data[1], `"`)

	emailsData := strings.Split(lines[1], ": ")
	if len(emailsData) != 2 {
		return &task, errors.New("emails_data is not correct, emails_data: " + lines[1])
	}

	emailsCnt, err := strconv.Atoi(strings.Trim(emailsData[1], " \n\r"))
	if err != nil {
		return &task, errors.New("count of emails is not integer, emails_count: " + emailsData[1])
	}
	task.EmailsCount = emailsCnt
	return &task, nil
}

func ParseEmailsTask(openaiResponse string) (*EmailsTask, error) {
	var task EmailsTask
	lines := strings.Split(openaiResponse, "\n")
	if len(lines) != 4 {
		logger.Error("can not parse openaiResponse with text: ", openaiResponse)
		return &task, errors.New("error during parsing openaiAns")
	}
	emailsData := strings.Split(strings.Trim(lines[2], " \n\r"), ": ")
	if len(emailsData) != 2 {
		logger.Error("can not parse openaiResponse with text: ", openaiResponse)
		return &task, errors.New("error during parsing openaiAns")
	}

	emailsCount, err := strconv.Atoi(emailsData[1])
	if err != nil {
		logger.Error("error during converting emailsCount, err: ", err)
		return &task, errors.New("error during parsing openaiAns")
	}
	task.EmailsCount = emailsCount
	return &task, nil
}

// GET_CALENDAR task parser
func ParseGetCalendarTask(openaiResponse string) (*GetCalendarTask, error) {
	var task GetCalendarTask
	responseData := strings.Split(strings.Trim(openaiResponse, " \n\r"), "\n")
	if len(responseData) != 4 {
		return &task, errors.New("response on GetCalendarTask is not correct")
	}
	timestampData := strings.Split(responseData[2], ": ")
	if len(timestampData) != 2 {
		return &task, errors.New("the timestamp line in GetCalendarTask is not correct")
	}
	task.Timestamp = strings.Trim(timestampData[1], ` \r\n"`)

	return &task, nil
}

// help add new event
func ParseHelpAddEventTask(openaiResponse string) (*HelpAddEventTask, error) {
	var task HelpAddEventTask
	responseData := strings.Split(strings.Trim(openaiResponse, " \n\r"), "\n")
	if len(responseData) != 5 {
		return &task, errors.New("HelpAddEventTaskData is not correct, response:" + openaiResponse)
	}
	timestampData := strings.Split(responseData[2], ": ")
	if len(timestampData) != 2 {
		return &task, errors.New("timestampStartData is not correct, data: " + responseData[2])
	}

	task.Timestamp = strings.Trim(timestampData[1], ` \r"`)

	reminderData := strings.Split(responseData[3], ": ")
	if len(reminderData) != 2 {
		return &task, errors.New("reminderData is not correct, reminderData: " + responseData[3])
	}

	task.Reminder = strings.Trim(reminderData[1], ` \r"`)
	return &task, nil
}

// add event task
func ParseAddEventTask(openaiResponse string) (*AddEventTask, error) {
	var task AddEventTask
	responseData := strings.Split(strings.Trim(openaiResponse, " \n\r"), "\n")
	if len(responseData) != 6 {
		return &task, errors.New("AddEventTaskData is not correct, length:" + string(len(responseData)) + " response:" + openaiResponse)
	}
	// parsing "timestamp_start" field
	timestampStartData := strings.Split(responseData[2], ": ")
	if len(timestampStartData) != 2 {
		return &task, errors.New("timestampStartData is not correct, data: " + responseData[2])
	}
	task.TimestampStart = strings.Trim(timestampStartData[1], ` \r"`)

	// parsing "timestamp_end" field
	timestampEndData := strings.Split(responseData[3], ": ")
	if len(timestampEndData) != 2 {
		return &task, errors.New("timestampEndData is not correct, data: " + responseData[3])
	}
	task.TimestampEnd = strings.Trim(timestampEndData[1], `\r"`)

	// parsing "reminder" field
	reminderData := strings.Split(responseData[4], ": ")
	if len(reminderData) != 2 {
		return &task, errors.New("reminderData is not correct, reminderData: " + responseData[3])
	}

	task.Reminder = strings.Trim(reminderData[1], ` \r"`)
	return &task, nil
}
