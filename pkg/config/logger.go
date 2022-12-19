package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger(
	logLevel string,
) *Logger {
	l := logrus.New()
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{}

	switch logLevel {
	case "DEBUG":
		l.Level = logrus.DebugLevel
	default:
		l.Level = logrus.ErrorLevel
	}

	return &Logger{
		log: l,
	}
}

func (l *Logger) Log(
	lvl logrus.Level,
	msg string,
) {

	fields := logrus.Fields{
		"instrumentation.provider": "newrelic-kubernetes-endpoint-scraper",
	}

	switch lvl {
	case logrus.ErrorLevel:
		l.log.WithFields(fields).Error(msg)
	default:
		l.log.WithFields(fields).Debug(msg)
	}
}

func (l *Logger) LogWithFields(
	lvl logrus.Level,
	msg string,
	attributes map[string]string,
) {

	fields := logrus.Fields{
		"instrumentation.provider": "newrelic-kubernetes-endpoint-scraper",
	}

	for key, val := range attributes {
		fields[key] = val
	}

	switch lvl {
	case logrus.ErrorLevel:
		l.log.WithFields(fields).Error(msg)
	default:
		l.log.WithFields(fields).Debug(msg)
	}
}
