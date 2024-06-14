package calendar_service

import (
	"fmt"
	"planpilot/internal/logger"
	"testing"
)

func TestCalendar(t *testing.T) {
	calendarSvc, err := New()
	if err != nil {
		t.Error("error during the creating calendar service, err: ", err)
	}
	events, err := calendarSvc.GetEventsByDate("2024-05-26T13:45:00Z")
	if err != nil {
		t.Error("error, err: ", err)
	}
	result, err := calendarSvc.EventsToString(events)
	logger.Info("Schedule:")
	fmt.Println(result)
}
