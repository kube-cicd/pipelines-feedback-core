package contract

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type KubernetesDependencies interface {
	InjectKubernetesContext(ctx KubernetesContext)
}

type KubernetesContext struct {
	recorder   record.EventRecorder
	kubeConfig *rest.Config
}

type WithInitialization interface {
	Initialize() error
}
