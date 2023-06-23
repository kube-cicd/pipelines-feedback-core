package logging

import (
	"context"
	"github.com/sirupsen/logrus"
)

type FatalLogger interface {
	Fatalf(format string, args ...interface{})
}

func NewInternalLogger() *InternalLogger {
	return &InternalLogger{
		parent: logrus.New().WithContext(context.TODO()),
	}
}

// InternalLogger is separating our codebase from logging implementation
type InternalLogger struct {
	parent *logrus.Entry
}

func (l *InternalLogger) Debug(args ...interface{}) {
	l.parent.Debug(args...)
}

func (l *InternalLogger) Info(args ...interface{}) {
	l.parent.Info(args...)
}

func (l *InternalLogger) Warning(args ...interface{}) {
	l.parent.Warning(args...)
}

func (l *InternalLogger) Error(args ...interface{}) {
	l.parent.Error(args...)
}

func (l *InternalLogger) Debugf(format string, args ...interface{}) {
	l.parent.Debugf(format, args...)
}

func (l *InternalLogger) Infof(format string, args ...interface{}) {
	l.parent.Infof(format, args...)
}

func (l *InternalLogger) Warningf(format string, args ...interface{}) {
	l.parent.Warningf(format, args...)
}

func (l *InternalLogger) Fatalf(format string, args ...interface{}) {
	l.parent.Fatalf(format, args...)
}

func (l *InternalLogger) Errorf(format string, args ...interface{}) {
	l.parent.Errorf(format, args...)
}

func (l *InternalLogger) ForkWithFields(ctx context.Context, fields map[string]interface{}) *InternalLogger {
	return &InternalLogger{
		l.parent.WithContext(ctx).WithFields(fields),
	}
}
