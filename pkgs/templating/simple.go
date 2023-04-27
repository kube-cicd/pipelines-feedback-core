package templating

import (
	"bytes"
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
import "text/template"

func TemplateDashboardUrl(templateStr string, kubeObject v1.Object, typeMeta v1.TypeMeta) (string, error) {
	return render(templateStr, "dashboard-url", map[string]interface{}{
		"job":        kubeObject,
		"name":       kubeObject.GetName(),
		"namespace":  kubeObject.GetNamespace(),
		"kind":       typeMeta.Kind,
		"apiVersion": typeMeta.APIVersion,
		"apiGroup":   typeMeta.GroupVersionKind().Group,
		"gvk":        typeMeta.GroupVersionKind(),
	})
}

func TemplateProgressComment(templateStr string, pipeline contract.PipelineInfo, buildId string) (string, error) {
	return render(templateStr, "progress-comment-template", map[string]interface{}{
		"pipeline":  pipeline,
		"commentId": buildId,
	})
}

func TemplateSummaryComment(templateStr string, pipeline contract.PipelineInfo, buildId string) (string, error) {
	return render(templateStr, "summary-comment-template", map[string]interface{}{
		"pipeline":  pipeline,
		"commentId": buildId,
	})
}

func render(templateStr string, name string, variables map[string]interface{}) (string, error) {
	t, err := template.New(name).Parse(templateStr)
	if err != nil {
		return "", errors.Wrap(err, "cannot create a "+name+". Please check your template")
	}
	var result bytes.Buffer
	execErr := t.Execute(&result, variables)
	if execErr != nil {
		return "", errors.Wrap(execErr, "cannot create a "+name+". Please check your template")
	}
	return result.String(), nil
}
