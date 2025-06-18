package logger

import (
	"log"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

func LogInfo(msg string) {
	log.Printf("%s[INFO]%s %s", green, reset, msg)
}

func LogWarn(msg string) {
	log.Printf("%s[WARN]%s %s", yellow, reset, msg)
}

func LogError(msg string) {
	log.Printf("%s[ERROR]%s %s", red, reset, msg)
}
