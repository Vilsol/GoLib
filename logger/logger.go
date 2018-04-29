package logger

import (
	"fmt"
	"time"
)

var (
	Green   = string([]byte{27, 91, 51, 48, 59, 52, 50, 109})
	White   = string([]byte{27, 91, 51, 48, 59, 52, 55, 109})
	Yellow  = string([]byte{27, 91, 51, 48, 59, 52, 51, 109})
	Red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	Blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	Magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	Cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	Reset   = string([]byte{27, 91, 48, 109})
)

type Logger struct {
	Prefix string
}

func New(prefix string) *Logger {
	return &Logger{
		Prefix: prefix,
	}
}

func (logger *Logger) Log(format string, data ...interface{}) {
	message := format

	if data != nil {
		if len(data) > 0 {
			message = fmt.Sprintf(format, data...)
		}
	}

	fmt.Printf("[%s] %v | %s\n",
		logger.Prefix,
		time.Now().Format("2006/01/02 - 15:04:05"),
		message,
	)
}
