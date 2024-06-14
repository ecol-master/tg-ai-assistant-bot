package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID  uint `gorm:"uniqueIndex;not null"`
	UserContext string
}

type UserTask struct {
	gorm.Model
	UserID    uint   `gorm:"index;not null"` // Reference to the User
	Schedule  string `gorm:"index;not null"`
	Timestamp time.Time
	Prompt    string `gorm:"not null"`
	Active    bool   `gorm:"not null"`
}

type UserEmailTask struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null"` // Reference to the User
	Schedule    string `gorm:"index;not null"`
	Timestamp   time.Time
	Prompt      string `gorm:"not null"`
	Active      bool   `gorm:"not null"`
	EmailsCount int
}
