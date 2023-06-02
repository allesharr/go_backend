package logger

import (
	"fmt"
	"time"
)

const (
	INFO = "INFO"
	WARN = "WARN"
	ERR  = "ERR"
)

type Logger struct{}

func (l Logger) Log(log_level string, message string) {
	fmt.Printf("[%s][%s]: %s\n", log_level, time.Now().Format("02.01.2006 15:04:05"), message)
}
