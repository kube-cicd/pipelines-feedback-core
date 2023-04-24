package logging

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Logger is a universal interface valid for both global and contextual logger
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// CreateLogger is creating a global logger instance
func CreateLogger(isDebugLevel bool) Logger {
	instance := logrus.New()
	instance.SetLevel(logrus.InfoLevel)
	if isDebugLevel {
		instance.SetLevel(logrus.DebugLevel)
	}
	return instance
}

// CreateK8sContextualLogger creates a logger instance with a Kubernetes controller's reconciliation request context
func CreateK8sContextualLogger(ctx context.Context, req ctrl.Request) Logger {
	id, _ := uuid.NewUUID()
	return logrus.WithContext(ctx).WithFields(map[string]interface{}{
		"request": id,
		"name":    req.NamespacedName,
	})
}
