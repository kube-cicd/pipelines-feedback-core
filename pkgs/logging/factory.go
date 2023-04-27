package logging

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

// CreateLogger is creating a global logger instance
func CreateLogger(isDebugLevel bool) Logger {
	instance := logrus.New()
	instance.SetLevel(logrus.InfoLevel)
	logrus.SetLevel(logrus.InfoLevel)
	if isDebugLevel {
		logrus.SetLevel(logrus.DebugLevel)
		instance.SetLevel(logrus.DebugLevel)
	}
	return instance
}

// CreateK8sContextualLogger creates a logger instance with a Kubernetes controller's reconciliation request context
func CreateK8sContextualLogger(ctx context.Context, req ctrl.Request) Logger {
	// todo: set log level
	id, _ := uuid.NewUUID()
	return logrus.WithContext(ctx).WithFields(map[string]interface{}{
		"request": id,
		"name":    req.NamespacedName,
	})
}
