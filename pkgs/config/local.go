package config

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"os"
)

type LocalFileConfigurationCollector struct {
}

func (lf *LocalFileConfigurationCollector) InjectDependencies(recorder record.EventRecorder, kubeconfig *rest.Config) error {
	return nil
}

func (lf *LocalFileConfigurationCollector) CanHandle(adapterName string) bool {
	return adapterName == lf.GetImplementationName()
}

func (lf *LocalFileConfigurationCollector) GetImplementationName() string {
	return "local"
}

// CollectInitially is looking for a JSON file specified in CONFIG_PATH environment variable
// it falls back to "pipelines-feedback.json" local file path. If the file does not exist it is ignored.
func (lf *LocalFileConfigurationCollector) CollectInitially() ([]*v1alpha1.PFConfig, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "pipelines-feedback.json"
	}
	stat, err := os.Stat(path)

	// the file is optional
	if os.IsNotExist(err) {
		return []*v1alpha1.PFConfig{}, nil
	}
	// unknown error
	if err != nil {
		return []*v1alpha1.PFConfig{}, errors.Wrap(err, "cannot load configuration file")
	}
	// not a file - a directory
	if stat.IsDir() {
		return []*v1alpha1.PFConfig{}, errors.New("is a directory, not a file")
	}

	// the file is valid, so lets parse it
	data, openErr := os.ReadFile(path)
	if openErr != nil {
		return []*v1alpha1.PFConfig{}, errors.Wrap(openErr, "cannot read configuration file")
	}
	pfc := v1alpha1.NewPFConfig()
	if unmarshalErr := json.Unmarshal(data, &pfc.Data); unmarshalErr != nil {
		return []*v1alpha1.PFConfig{}, errors.Wrap(unmarshalErr, "cannot unmarshal file from JSON into struct")
	}
	return []*v1alpha1.PFConfig{&pfc}, nil
}

// CollectOnRequest is not implemented as the "local" does not allow dynamic resolution
func (lf *LocalFileConfigurationCollector) CollectOnRequest(info contract.PipelineInfo) ([]*v1alpha1.PFConfig, error) {
	return []*v1alpha1.PFConfig{}, nil
}
