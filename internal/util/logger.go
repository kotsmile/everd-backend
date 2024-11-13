package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LoggerFields map[string]any

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)

	WithField(key string, value any) Logger
	WithFields(fields LoggerFields) Logger
	WithError(err error) Logger
}

type logger struct {
	internal *logrus.Entry
}

func NewLoggerTest() Logger {
	return NewLogger(true)
}

func NewLogger(debug bool) Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	l.SetLevel(logrus.InfoLevel)
	l.SetOutput(os.Stdout)

	if debug {
		l.SetLevel(logrus.DebugLevel)
	}

	return &logger{
		internal: logrus.NewEntry(l),
	}
}

func (l *logger) WithField(key string, value any) Logger {
	return &logger{
		internal: l.internal.WithField(key, value),
	}
}

func (l *logger) WithFields(fields LoggerFields) Logger {
	return &logger{
		internal: l.internal.WithFields(logrus.Fields(fields)),
	}
}

func (l *logger) WithError(err error) Logger {
	return &logger{
		internal: l.internal.WithError(err),
	}
}

func (l *logger) Debug(args ...any) {
	l.internal.Debug(args...)
}

func (l *logger) Info(args ...any) {
	l.internal.Info(args...)
}

func (l *logger) Warn(args ...any) {
	l.internal.Warn(args...)
}

func (l *logger) Error(args ...any) {
	l.internal.Error(args...)
}

func (l *logger) Fatal(args ...any) {
	l.internal.Fatal(args...)
}

func (l *logger) Debugf(format string, args ...any) {
	l.internal.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.internal.Infof(format, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.internal.Warnf(format, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.internal.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.internal.Fatalf(format, args...)
}
