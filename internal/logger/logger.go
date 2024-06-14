package logger

import (
	"log"
	"os"
)

var (
	logInfo *log.Logger
	logErr  *log.Logger
)

func init() {
	logInfo = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	logErr = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func Info(args ...any) {
	logInfo.Println(args...)
}

func Error(args ...any) {
	logErr.Println(args...)
}

func Fatal(args ...any) {
	logErr.Fatal(args...)
}
