package templating_test

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/templating"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestTemplateDashboardUrl(t *testing.T) {
	job := v1.Job{
		TypeMeta:   metav1.TypeMeta{Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{Name: "peppa", Namespace: "default"},
		Spec:       v1.JobSpec{},
		Status:     v1.JobStatus{},
	}

	result, err := templating.TemplateDashboardUrl("https://console-openshift-console.apps.my-cluster.org/k8s/ns/{{ .namespace }}/tekton.dev~v1beta1~PipelineRun/{{ .name }}", &job, job.TypeMeta)
	assert.Nil(t, err)
	assert.Equal(t, "https://console-openshift-console.apps.my-cluster.org/k8s/ns/default/tekton.dev~v1beta1~PipelineRun/peppa", result)
}
