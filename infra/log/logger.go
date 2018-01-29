package log

import "github.com/juju/loggo"

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Level int

const (
	Leveldebug Level = iota
	LevelInfo
	LevelError
)

func SetLevel(level Level) {
	var config string
	switch level {
	case LevelInfo:
		config = "<root>=INFO"
	case LevelError:
		config = "<root>=ERROR"
	case Leveldebug:
		fallthrough
	default:
		config = "<root>=DEBUG"
	}

	loggo.ConfigureLoggers(config)
}

func NewLogger(packageName string) Logger {
	return loggo.GetLogger(packageName)
}
