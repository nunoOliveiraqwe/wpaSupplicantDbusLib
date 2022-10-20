package wpaSuppDBusLib

import (
	"log"
	"os"
)

type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type wpaDefaultLogger struct {
	srcLogger *log.Logger
}

func newDefaultLogger() Logger {
	logger := wpaDefaultLogger{srcLogger: log.New(os.Stdout, "wpaGoDbus", log.Ldate|log.Ltime|log.Lshortfile)}
	return &logger
}

func (l *wpaDefaultLogger) Info(args ...interface{}) {
	l.srcLogger.Print("[INFO] ", args)
}

func (l *wpaDefaultLogger) Warn(args ...interface{}) {
	l.srcLogger.Print("[WARN] ", args)
}

func (l *wpaDefaultLogger) Error(args ...interface{}) {
	l.srcLogger.Print("[ERROR] ", args)
}
