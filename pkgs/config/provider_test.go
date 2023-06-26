package config_test

import (
	"context"
	internalConfig "github.com/kube-cicd/pipelines-feedback-core/internal/config"
	store2 "github.com/kube-cicd/pipelines-feedback-core/internal/store"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/apis/pipelinesfeedback.keskad.pl/v1alpha1"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/config"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/logging"
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/store"
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"testing"
)

func TestConfigurationProvider_FetchContextual(t *testing.T) {
	pfc := v1alpha1.NewPFConfig()
	pfc.ObjectMeta.Namespace = "books"
	pfc.ObjectMeta.Name = "cfg"
	pfc.Data = v1alpha1.Data{
		"jxscm.kropotkin": "The Bread Book",
	}

	validator := internalConfig.NullValidator{}
	docStore := internalConfig.CreateIndexedDocumentStore(&validator)
	_ = docStore.Push(&pfc)

	cp, err := config.NewConfigurationProvider(
		docStore,
		logging.NewInternalLogger(),
		&v1.CoreV1Client{},
		store.Operator{},
		&validator,
	)

	data := cp.FetchContextual("jxscm", "books", contract.PipelineInfo{})

	// Case: When key exists, then its value is returned
	assert.Nil(t, err)
	assert.True(t, data.HasKey("kropotkin"))
	assert.Equal(t, "The Bread Book", data.Get("kropotkin"))

	// Case: Key that is not defined is returning default value
	assert.Equal(t, "yup", data.GetOrDefault("this-key-does-not-exists", "yup"))
}

func TestConfigurationProvider_FetchGlobal(t *testing.T) {
	pfc := v1alpha1.NewPFConfig()
	pfc.ObjectMeta.Namespace = "books"
	pfc.ObjectMeta.Name = "cfg"
	pfc.Data = v1alpha1.Data{
		"jxscm.kropotkin": "The Bread Book",
	}

	validator := internalConfig.NullValidator{}
	docStore := internalConfig.CreateIndexedDocumentStore(&validator)
	_ = docStore.Push(&pfc)

	cp, err := config.NewConfigurationProvider(
		docStore,
		logging.NewInternalLogger(),
		&v1.CoreV1Client{},
		store.Operator{},
		&validator,
	)
	assert.Nil(t, err)

	// Case: No global (cluster-scope) resource created
	f := cp.FetchGlobal("jxscm")
	assert.Equal(t, "", f.Get("kropotkin"))

	// Case: cluster-scope resource created
	pfcG := v1alpha1.NewPFConfig()
	pfcG.ObjectMeta.Namespace = ""
	pfcG.ObjectMeta.Name = "global-cfg"
	pfcG.Data = v1alpha1.Data{
		"jxscm.kropotkin": "A global baking book!",
	}
	docStore.Push(&pfcG)

	data := cp.FetchGlobal("jxscm")
	assert.Equal(t, "A global baking book!", data.Get("kropotkin"))
}

func TestConfigurationProvider_FetchSecretKey_FetchesKeyFromKubernetesSecret(t *testing.T) {
	validator := internalConfig.NullValidator{}
	docStore := internalConfig.CreateIndexedDocumentStore(&validator)

	// mock Kubernetes kind: Secrets
	secret := v12.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "scm-secret",
			Namespace: "books",
		},
		Data: map[string][]byte{
			"token": []byte("blehblehbleh"),
		},
	}

	kubernetesClient := fake.NewSimpleClientset(&secret)
	coreV1 := kubernetesClient.CoreV1()

	cp, err := config.NewConfigurationProvider(
		docStore,
		logging.NewInternalLogger(),
		coreV1,
		store.Operator{Store: store2.NewMemory()},
		&validator,
	)
	assert.Nil(t, err)

	// FetchSecretKey() should fetch a Kubernetes secret and return value of it's "token" key (.Data.token)
	val, fetchErr := cp.FetchSecretKey(context.TODO(), "scm-secret", "books", "token", false)
	assert.Nil(t, fetchErr)
	assert.Equal(t, "blehblehbleh", val)
}

func TestConfigurationProvider_FetchFromFieldOrSecret_FetchesFromSecret(t *testing.T) {
	validator := internalConfig.NullValidator{}
	docStore := internalConfig.CreateIndexedDocumentStore(&validator)

	// mock Kubernetes kind: Secrets
	secret := v12.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "scm-secret",
			Namespace: "books",
		},
		Data: map[string][]byte{
			"token": []byte("blehblehbleh"),
		},
	}

	kubernetesClient := fake.NewSimpleClientset(&secret)
	coreV1 := kubernetesClient.CoreV1()

	cp, err := config.NewConfigurationProvider(
		docStore,
		logging.NewInternalLogger(),
		coreV1,
		store.Operator{Store: store2.NewMemory()},
		&validator,
	)
	assert.Nil(t, err)

	data := config.NewData("jxscm", map[string]string{
		// no "token": "..." defined there
		// so the token will be taken from kind: Secret as defined there:
		"token-key":   "token",      // .data.token
		"secret-name": "scm-secret", // kind: Secret, name=scm-secret
	}, &validator, logging.NewInternalLogger())

	val, fetchErr := cp.FetchFromFieldOrSecret(context.TODO(), &data, "books", "token", "token-key", "secret-name")
	assert.Nil(t, fetchErr)
	assert.Equal(t, "blehblehbleh", val)
}

func TestConfigurationProvider_FetchFromFieldOrSecret_FetchesFromConfigKeyFirst(t *testing.T) {
	validator := internalConfig.NullValidator{}
	docStore := internalConfig.CreateIndexedDocumentStore(&validator)

	kubernetesClient := fake.NewSimpleClientset()
	coreV1 := kubernetesClient.CoreV1()

	cp, err := config.NewConfigurationProvider(
		docStore,
		logging.NewInternalLogger(),
		coreV1,
		store.Operator{Store: store2.NewMemory()},
		&validator,
	)
	assert.Nil(t, err)

	data := config.NewData("jxscm", map[string]string{
		"token":       "hello-world",
		"token-key":   "token",      // .data.token
		"secret-name": "scm-secret", // kind: Secret, name=scm-secret
	}, &validator, logging.NewInternalLogger())

	val, fetchErr := cp.FetchFromFieldOrSecret(context.TODO(), &data, "books", "token", "token-key", "secret-name")
	assert.Nil(t, fetchErr)
	assert.Equal(t, "hello-world", val)
}
