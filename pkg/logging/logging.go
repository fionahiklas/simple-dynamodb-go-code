//go:generate mockgen -package logging_test -destination=./mock_logging_test.go -source $GOFILE
package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type SimpleLogger interface {
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type config interface {
	LogLevel() string
	EnableJsonLogFormat() bool
}

func NewLogger(config config) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logLevelFromString(log, config.LogLevel()))
	if config.EnableJsonLogFormat() {
		log.SetFormatter(&logrus.JSONFormatter{})
	}
	return log
}

func NewTestLogger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	return log
}

func logLevelFromString(log *logrus.Logger, logLevelString string) logrus.Level {
	resultLevel, err := logrus.ParseLevel(strings.ToLower(logLevelString))
	if err != nil {
		log.Errorf("Couldn't parse log level: %s, defaulting to INFO", logLevelString)
		resultLevel = logrus.InfoLevel
	}
	return resultLevel
}
