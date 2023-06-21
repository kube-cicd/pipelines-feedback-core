package logging

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	ctrl "sigs.k8s.io/controller-runtime"
)

// CreateLogger is creating a global logger instance
func CreateLogger(isDebugLevel bool) *InternalLogger {
	instance := logrus.New()
	instance.SetLevel(logrus.InfoLevel)
	logrus.SetLevel(logrus.InfoLevel)
	if isDebugLevel {
		logrus.SetLevel(logrus.DebugLevel)
		instance.SetLevel(logrus.DebugLevel)
	}
	return &InternalLogger{instance.WithContext(context.TODO())}
}

// CreateK8sContextualLogger creates a logger instance with a Kubernetes controller's reconciliation request context
func CreateK8sContextualLogger(ctx context.Context, mainLogger *InternalLogger, req ctrl.Request) *InternalLogger {
	id, _ := uuid.NewUUID()
	return mainLogger.ForkWithFields(ctx, map[string]interface{}{
		"request": id,
		"name":    req.NamespacedName,
	})
}
