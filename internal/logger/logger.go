package logger

import (
	"log"
	"os"
)

type _level int

const (
	Debug _level = iota
	Info
	Warn
	Error
	Fatal
)

func (level _level) Name() string {
	switch level {
	case 0:
		return "DEBUG"
	case 1:
		return "INFO"
	case 2:
		return "WARN"
	case 3:
		return "ERROR"
	case 4:
		return "FATAL"
	default:
		panic("Unreachable")
	}
}

var (
	LoggerLevel = Debug
	_stdout     = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	_stderr     = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
)

type Logger struct {
	prefix string
}

func GetLogger(prefix string) *Logger {
	return &Logger{prefix}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.print(Debug, format, a)
}

func (logger *Logger) Info(format string, a ...interface{}) {
	logger.print(Info, format, a)
}

func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.print(Warn, format, a)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.print(Error, format, a)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	_stderr.Fatalf(format, a...)
}

func (logger *Logger) print(level _level, format string, a []interface{}) {
	if level >= LoggerLevel {
		p := append([]interface{}{level.Name(), logger.prefix}, a...)
		if level >= Warn {
			_stderr.Printf("%s %s "+format+"\n", p...)
		} else {
			_stdout.Printf("%s %s "+format+"\n", p...)
		}
	}
}
