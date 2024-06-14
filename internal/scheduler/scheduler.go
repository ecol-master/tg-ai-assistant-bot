package scheduler

import (
	"fmt"
	"log"
	"planpilot/internal/config"
	"planpilot/internal/db"
	"planpilot/internal/gmail_service"
	"planpilot/internal/logger"
	"planpilot/internal/openai"
	"time"

	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var scheduler, _ = gocron.NewScheduler()

func Start() {
	logger.Info("scheduler started")
	scheduler.Start()

	tasks, err := db.GetActiveTasks()
	if err != nil {
		log.Fatal("Error getting previously defined tasks")
		return
	}

	for _, task := range tasks {
		fmt.Println("task: ", task)
		// err := AddTask(task)
		// if err != nil {
		// 	log.Fatalf("Error while creating task %d", task.ID)
		// }
	}
}

func AddTask(task db.UserTask, bot *tgbotapi.BotAPI) error {
	var jobDefenition gocron.JobDefinition
	logger.Info("timestamp start: ", task.Timestamp, " now: ", time.Now())
	if task.Schedule != "" {
		jobDefenition = gocron.CronJob(
			task.Schedule, false,
		)
	} else {
		jobDefenition = gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(task.Timestamp.Local().Add(-1 * time.Minute * config.REMIND_UNTIL_EVENT_START)))
		logger.Info("ONE TIME JOB: ", task.Timestamp.Local().Add(-1*time.Minute*config.REMIND_UNTIL_EVENT_START))
	}
	job, err := scheduler.NewJob(
		jobDefenition,
		gocron.NewTask(
			func(userTask db.UserTask, bot *tgbotapi.BotAPI) {
				logger.Info("Executing UserTask: ", userTask.Prompt)
				msg := tgbotapi.NewMessage(int64(task.UserID), fmt.Sprintf("I remind you about: `%s`", userTask.Prompt))
				msg.ParseMode = "Markdown"
				bot.Send(msg)
				db.DeleteTask(&userTask)
			},
			task,
			bot,
		),
	)

	if err != nil {
		logger.Error("Error creating job", err)
	} else {
		log.Print(job.ID().String() + "started")
	}

	return err
}

func AddEmailsTask(task db.UserTask, bot *tgbotapi.BotAPI) error {
	job, err := scheduler.NewJob(
		gocron.CronJob(
			task.Schedule, false,
		),
		gocron.NewTask(
			func(userTask db.UserTask) {
				logger.Info("Executing UserTask: ", userTask.Prompt)
				var msg tgbotapi.MessageConfig

				responseTask, err := openai.ParseEmailScheduledTask(task.Prompt)
				if err != nil {
					logger.Error("error during parsing scheduled email task, err: ", err)
					msg = tgbotapi.NewMessage(int64(task.UserID), "Sorry, i can't to analyze your mail box")
					bot.Send(msg)
					return
				}

				gmailService, err := gmail_service.New()

				if err != nil {
					logger.Error("error during creating gmailService: ", err)
					msg = tgbotapi.NewMessage(int64(task.UserID), "Sorry, i can't to analyze your mail box")
					bot.Send(msg)
					return
				}

				compressedEmails, err := gmailService.MakeCompressedEmailsText(responseTask.EmailsCount)

				if err != nil {
					logger.Error("error during getting compressed emails: ", err)
					msg = tgbotapi.NewMessage(int64(task.UserID), "Sorry, i can't to analyze your mail box")
				} else {
					msg = tgbotapi.NewMessage(int64(task.UserID), compressedEmails)
				}
				msg.ParseMode = "Markdown"
				bot.Send(msg)
			},
			task,
		),
	)

	if err != nil {
		logger.Error("Error creating job", err)
	} else {
		log.Print(job.ID().String() + "started")
	}

	return err
}
