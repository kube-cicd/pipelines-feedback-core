package config

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

type FakeValidator struct {

}

func (f *FakeValidator) ValidateRequestedEntry(group string, key string) error {
	return nil
}

func (f *FakeValidator) ValidateConfig(data v1alpha1.Data) error {
	return nil
}

func (f *FakeValidator) Add(schema Schema) {

}

type FakeCollector struct {
}

func (f *FakeCollector) SetLogger(logger *logging.InternalLogger) {
}

func (f *FakeCollector) InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error {
	return nil
}

func (f *FakeCollector) CanHandle(adapterName string) bool {
	return true
}

func (f *FakeCollector) GetImplementationName() string {
	return "fake"
}

func (f *FakeCollector) CollectInitially() ([]*v1alpha1.PFConfig, error) {
	return []*v1alpha1.PFConfig{
		{
			TypeMeta: v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{
				Name:      "bread",
				Namespace: "books",
			},
			Spec: v1alpha1.Spec{},
			Data: v1alpha1.Data{
				"rating": "What a tasty book!",
			},
			Status: v1alpha1.PFCStatus{},
		},
	}, nil
}

func (f *FakeCollector) CollectOnRequest(info contract.PipelineInfo) ([]*v1alpha1.PFConfig, error) {
	return []*v1alpha1.PFConfig{
		{
			TypeMeta: v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{
				Name:      "mutual-aid-a-factor-of-revolution",
				Namespace: "books",
			},
			Spec: v1alpha1.Spec{},
			Data: v1alpha1.Data{
				"rating": "What a fantastic book!",
			},
			Status: v1alpha1.PFCStatus{},
		},
	}, nil
}
