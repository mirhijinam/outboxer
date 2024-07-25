package logger

import (
	"log"
	"os"
)

type Logger struct {
	filepath string
	info     *log.Logger
	warn     *log.Logger
	err      *log.Logger
}

func New(filepath string) (*Logger, error) {
	openFlags := os.O_APPEND | os.O_CREATE | os.O_WRONLY

	fileInfo, err := os.OpenFile(filepath, openFlags, 0666)
	if err != nil {
		return nil, err
	}

	fileWarn, err := os.OpenFile(filepath, openFlags, 0666)
	if err != nil {
		return nil, err
	}

	fileErr, err := os.OpenFile(filepath, openFlags, 0666)
	if err != nil {
		return nil, err
	}

	logFlags := log.LstdFlags

	logInfo := log.New(fileInfo, "INFO:\t", logFlags)
	logWarn := log.New(fileWarn, "WARN:\t", logFlags)
	logErr := log.New(fileErr, "ERR:\t", logFlags)

	return &Logger{
		filepath: filepath,
		info:     logInfo,
		warn:     logWarn,
		err:      logErr,
	}, nil
}

func (l *Logger) Info(text ...any) {
	l.info.Println(text...)
}

func (l *Logger) Warn(text ...any) {
	l.warn.Println(text...)
}

func (l *Logger) Err(text ...any) {
	l.err.Println(text...)
}
