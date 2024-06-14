package db

import (
	"fmt"
	"planpilot/internal/config"
	"planpilot/internal/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	cfg := config.New()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB_HOST,
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_NAME,
		cfg.DB_PORT,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("can not connect to postgresql db")
	}

	if err = DB.AutoMigrate(&UserTask{}, &User{}); err != nil {
		panic("can not aply migrations in db")
	}
}

func CreateTask(schedule string, timestamp time.Time, prompt string, active bool) (*UserTask, error) {
	task := UserTask{
		Schedule:  schedule,
		Timestamp: timestamp,
		Prompt:    prompt,
		Active:    active,
	}
	result := DB.Create(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func DeleteTask(task *UserTask) error {
	result := DB.Delete(task)
	return result.Error
}

func GetActiveTasks() ([]UserTask, error) {
	var tasks []UserTask
	result := DB.Where("active = ?", true).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func CreateUser(telegramID uint) (*User, error) {
	user := User{
		TelegramID: telegramID,
	}
	result := DB.Create(&user)
	logger.Info("created user: ", user)
	return &user, result.Error
}

func GetUser(telegramID uint) (*User, error) {
	var user User
	result := DB.Where("telegram_id = ?", telegramID).Find(&user)
	return &user, result.Error
}
