package config

import (
	"fmt"
	"os"
	"path/filepath"
	"planpilot/internal/logger"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env.local into the system, .env.local has the same structure as .env file
	envFilename := ".env.local"
	currentDir, err := os.Getwd()

	if err != nil {
		logger.Fatal("error to getting the current filepath")
	}

	var envFilepath string

	if strings.HasSuffix(currentDir, "cmd") {
		envFilepath = filepath.Join(currentDir, fmt.Sprintf("../%s", envFilename))
	} else {
		envFilepath = filepath.Join(currentDir, fmt.Sprintf("./%s", envFilename))
	}

	if err := godotenv.Load(envFilepath); err != nil {
		logger.Error(fmt.Sprintf("No %s file found", envFilename))
	}
}

type Config struct {
	// Tokens
	TELEGRAM_TOKEN string
	OPENAI_TOKEN   string

	// Database
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSLMODE  string

	// Email
	EMAIL_USERNAME          string // your email address
	EMAIL_PASSWORD          string //
	EMAIL_IMAP_ADDRESS      string // google imap 'imap.google.com'
	EMAIL_IMAP_ADDRESS_PORT int    // default port is - 993
}

func New() *Config {
	imapPortString := getEnv("EMAIL_IMAP_ADDRESS_PORT", "")
	imapPort, err := strconv.Atoi(imapPortString)
	if err != nil {
		imapPort = 993
	}

	return &Config{
		TELEGRAM_TOKEN:          getEnv("TELEGRAM_TOKEN", ""),
		OPENAI_TOKEN:            getEnv("OPENAI_TOKEN", ""),
		DB_HOST:                 getEnv("DB_HOST", ""),
		DB_PORT:                 getEnv("DB_PORT", ""),
		DB_USER:                 getEnv("DB_USER", ""),
		DB_PASSWORD:             getEnv("DB_PASSWORD", ""),
		DB_NAME:                 getEnv("DB_NAME", ""),
		DB_SSLMODE:              getEnv("DB_SSLMODE", ""),
		EMAIL_USERNAME:          getEnv("EMAIL_USERNAME", ""),
		EMAIL_PASSWORD:          getEnv("EMAIL_PASSWORD", ""),
		EMAIL_IMAP_ADDRESS:      getEnv("EMAIL_IMAP_ADDRESS", ""),
		EMAIL_IMAP_ADDRESS_PORT: imapPort,
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Time in project
const (
	REMIND_UNTIL_EVENT_START = 5

	TIME_LAYOUT = time.RFC3339

	EMAIL_TEMPLATE_RESPONSE = "" +
		"__From__ : `%s`\n" +
		"__Text__ : %s\n\n"
)

// template string to represent event in calendar_service.EventsToString
const (
	EVENT_TEMPLATE = "" +
		"%d. `%s` - `%s` %s\n"

	TIME_ZONE = "Europe/Moscow"
)
