package calendar_service

import (
	"fmt"
	"planpilot/internal/config"
	"planpilot/internal/logger"
	"planpilot/internal/openai"
	"time"

	"google.golang.org/api/calendar/v3"
)

func (c *CalendarService) GetEventsByDate(timestampt string) (*calendar.Events, error) {
	startDate, err := time.Parse(config.TIME_LAYOUT, timestampt)
	if err != nil {
		logger.Error("error during converting the timestamp in GetEventsByDate, err: ", err)
		return nil, err
	}

	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	logger.Info(startDate.Year(), startDate.Day())
	nextDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC).Add(time.Second * 60 * 60 * 24)
	logger.Info(nextDay)
	events, err := c.srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(startDate.Format(config.TIME_LAYOUT)).TimeMax(nextDay.Format(config.TIME_LAYOUT)).MaxResults(10).OrderBy("startTime").Do()

	return events, err
}

func (c *CalendarService) EventsToString(events *calendar.Events) (string, error) {
	var result string

	for i, event := range events.Items {
		startDate, err := makeDateToString(event.Start.DateTime)
		if err != nil {
			logger.Error("error during parsing the startDate from event, err: ", err)
			return "", err
		}
		endDate, err := makeDateToString(event.End.DateTime)
		if err != nil {
			logger.Error("error during parsing the endDate from event, err: ", err)
			return "", nil
		}

		result += fmt.Sprintf(config.EVENT_TEMPLATE, i+1, startDate, endDate, event.Summary)
	}

	return result, nil
}

// help functoin to convert date to string
func makeDateToString(date string) (string, error) {
	var result string
	dt, err := time.Parse(config.TIME_LAYOUT, date)
	if err != nil {
		return "", err
	}
	result = dt.Format("15:04")

	return result, nil
}

func (c *CalendarService) createEvent(openaiResponse string) (*calendar.Event, error) {
	var event calendar.Event

	addEventTask, err := openai.ParseAddEventTask(openaiResponse)
	if err != nil {
		return &event, err
	}
	eventStartDateTime, err := makeEventDateTime(addEventTask.TimestampStart)
	if err != nil {
		logger.Error("error during the creating calendar.EventDateTime start, err: ", err)
		return nil, err
	}
	eventEndDateTime, err := makeEventDateTime(addEventTask.TimestampEnd)
	if err != nil {
		logger.Error("error during the creating calendar.EventDateTime end, err: ", err)
		return nil, err
	}

	event = calendar.Event{
		Summary: addEventTask.Reminder,
		Start:   eventStartDateTime,
		End:     eventEndDateTime,
	}

	return &event, nil
}

func makeEventDateTime(timestampStart string) (*calendar.EventDateTime, error) {
	var eventDateTime calendar.EventDateTime

	tm, err := time.Parse(config.TIME_LAYOUT, timestampStart)
	if err != nil {
		return nil, err
	}

	// eventDateTime.Date = tm.Format("2006-01-02")
	eventDateTime.DateTime = tm.Format(config.TIME_LAYOUT)
	eventDateTime.TimeZone = config.TIME_ZONE

	return &eventDateTime, nil
}

func (c *CalendarService) AddEvent(openaiResponse string) (*calendar.Event, error) {
	newEvent, err := c.createEvent(openaiResponse)
	if err != nil {
		return nil, err
	}
	return c.srv.Events.Insert("primary", newEvent).Do()
}
