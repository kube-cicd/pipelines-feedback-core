package controller

import (
	"context"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ConfigurationController is reconciling CRD that provides configuration
type ConfigurationController struct {
}

func (cc *ConfigurationController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (cc *ConfigurationController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// For(). // todo crd
		Complete(cc)
}
