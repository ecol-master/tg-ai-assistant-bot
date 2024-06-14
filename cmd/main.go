package main

import (
	"fmt"
	"os"
	"planpilot/internal/bot"
	"planpilot/internal/calendar_service"
	"planpilot/internal/config"
	"planpilot/internal/logger"
	"planpilot/internal/scheduler"
	"time"
)

func main() {
	args := os.Args[1:]
	fmt.Println(args)

	testMode := false
	for _, arg := range args {
		if arg == "-test" {
			testMode = false
		}
	}

	if !testMode {
		err := bot.Run()
		scheduler.Start()
		logger.Error("app stopped with error: ", err)
	} else {
		calendarSvc, err := calendar_service.New()
		if err != nil {
			logger.Error("error during creating the calendar service, err: ", err)
		}

		calendarSvc.GetEventsByDate(time.Now().Format(config.TIME_LAYOUT))
	}
}
