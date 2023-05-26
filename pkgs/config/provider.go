package config

import (
	"context"
	"fmt"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/internal/config"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/logging"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/store"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"strings"
)

// NewConfigurationProvider is a constructor
func NewConfigurationProvider(docStore config.IndexedDocumentStore, logger logging.Logger,
	kubeConfig *rest.Config, kvStore store.Operator, cfgSchema *SchemaValidator) (ConfigurationProvider, error) {

	client, err := v1.NewForConfig(kubeConfig)
	if err != nil {
		return ConfigurationProvider{}, errors.Wrap(err, "cannot construct ConfigurationProvider, "+
			"Kubernetes Core API v1 construction error")
	}
	return ConfigurationProvider{
		docStore:      docStore,
		logger:        logger,
		secretsClient: client,
		stateStore:    kvStore,
		cfgSchema:     cfgSchema,
	}, nil
}

// ConfigurationProvider is serving already collected configuration. Served configuration is already merged from various sources
type ConfigurationProvider struct {
	docStore      config.IndexedDocumentStore
	logger        logging.Logger
	secretsClient *v1.CoreV1Client
	stateStore    store.Operator
	cfgSchema     *SchemaValidator
}

// todo: implement CollectOnRequest()

// FetchContextual is retrieving a final configuration in context of a given contract.PipelineInfo
func (cp *ConfigurationProvider) FetchContextual(component string, namespace string, pipeline contract.PipelineInfo) Data {
	cp.logger.Debugf("fetchContextual(%s, %s)", namespace, pipeline.GetFullName())
	endMap := make(map[string]string)
	for _, doc := range cp.docStore.GetForNamespace(namespace) {
		cp.logger.Debugf("fetchContextual => config '%s' available for this namespace, checking if matches", doc.Name)
		if doc.IsForPipeline(pipeline) {
			cp.logger.Debugf("fetchContextual(%s, %s) => using config '%s'", namespace, pipeline.GetFullName(), doc.Name)
			endMap = mergeMaps(endMap, doc.Data)
		}
	}
	return NewData(component, transformMapByComponent(endMap, component), cp.cfgSchema)
}

// FetchGlobal is fetching a global configuration for given component (without a context of a Pipeline)
func (cp *ConfigurationProvider) FetchGlobal(component string) Data {
	endMap := make(map[string]string)
	for _, doc := range cp.docStore.GetForNamespace("") {
		endMap = mergeMaps(endMap, doc.Data)
	}
	return NewData(component, transformMapByComponent(endMap, component), cp.cfgSchema)
}

// transformMapByComponent is stripping map out of other component keys and removing the component prefixes
func transformMapByComponent(input map[string]string, component string) map[string]string {
	output := make(map[string]string)
	for key, val := range input {
		if strings.HasPrefix(key, component+".") {
			newKey := key[len(component+"."):]
			output[newKey] = val
		}
	}
	return output
}

// FetchSecretKey is fetching a key from .data section from a Kubernetes Secret, directly from the Cluster's API
// Use this method in pair with FetchContextual() to at first fetch the key name, then to fetch the secret key
func (cp *ConfigurationProvider) FetchSecretKey(ctx context.Context, name string,
	namespace string, key string, cache bool) (string, error) {

	if cache {
		if cachedSecret := cp.stateStore.GetConfigSecretKey(namespace, name, key); cachedSecret != "" {
			return cachedSecret, nil
		}
	}

	secret, err := cp.secretsClient.Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "cannot fetch Kubernetes secret '%s/%s'", namespace, name)
	}
	val, exists := secret.Data[key]
	if !exists {
		return "", errors.New(fmt.Sprintf("the secret '%s/%s' does not contain key '%s'", namespace, name, key))
	}

	cp.stateStore.PushConfigSecretKey(namespace, name, key, string(val)) // Update cache
	return string(val), nil
}

// FetchFromFieldOrSecret allows to use an inline secret from configuration file (if present), fallbacks to fetching a secret key from a Kubernetes secret
func (cp *ConfigurationProvider) FetchFromFieldOrSecret(ctx context.Context, data *Data, namespace string, fieldKey string, referenceKey string, referenceSecretNameKey string) (string, error) {
	if data.HasKey(fieldKey) {
		return data.GetOrDefault(fieldKey, ""), nil
	}
	if referenceSecretNameKey != "" {
		// When the field with `kind: Secret` reference name is empty in the config data.
		// So we do not know which Kubernetes `kind: Secret` to open
		referenceSecretName := data.GetOrDefault(referenceSecretNameKey, "")
		if referenceSecretName == "" {
			return "", errors.New(fmt.Sprintf("'%s' should contain a valid Kubernetes secret name", referenceSecretNameKey))
		}
		// Try to fetch a .data.${referenceKey} from `kind: Secret` named ${referenceSecretNameKey}
		referenceKeyVal := data.GetOrDefault(referenceKey, "")
		val, err := cp.FetchSecretKey(ctx, referenceSecretName, namespace, referenceKeyVal, true)
		if err != nil {
			return "", errors.Wrap(err, "cannot fetch secret from Kubernetes")
		}
		return val, nil
	}
	// no ${referenceSecretNameKey} field specified neither ${fieldKey}
	return "", nil
}

func mergeMaps(m1 map[string]string, m2 map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range m1 {
		merged[k] = v
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}
