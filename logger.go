package dagflow

import (
	"log"
	"os"
)


// A Logger is a minimalistic interface for the library to log messages to. Should
// be used to provide custom logging writers for the library to use.


type LogPriority int
const (
	DEBUG LogPriority = iota
	INFO
	WARN
	ERROR
)


type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

type defaultLogger struct {
	logger *log.Logger
	priority LogPriority
}


func NewDefaultLogger(priority LogPriority, prefix string) *defaultLogger {
	return &defaultLogger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
		priority: priority,
	}
}

func (dl *defaultLogger) log(priority LogPriority, args ...interface{}) {
	if dl.priority <= priority {
		dl.logger.Print(args...)
	}
}

func (dl *defaultLogger) logf(priority LogPriority, format string, args ...interface{}) {
	if dl.priority <= priority {
		dl.logger.Printf(format, args...)
	}
}


func (dl *defaultLogger) Debug(args ...interface{}) {
	dl.log(DEBUG, args...)
}

func (dl *defaultLogger) Info(args ...interface{}) {
	dl.log(INFO, args...)
}

func (dl *defaultLogger) Warn(args ...interface{}) {
	dl.log(WARN, args...)
}

func (dl *defaultLogger) Error(args ...interface{}) {
	dl.log(ERROR, args...)
}

func (dl *defaultLogger) Debugf(format string, args ...interface{}) {
	dl.logf(DEBUG, format, args...)
}

func (dl *defaultLogger) Infof(format string, args ...interface{}) {
	dl.logf(INFO, format, args...)
}

func (dl *defaultLogger) Warnf(format string, args ...interface{}) {
	dl.logf(WARN, format, args...)
}

func (dl *defaultLogger) Errorf(format string, args ...interface{}) {
	dl.logf(ERROR, format, args...)
}

