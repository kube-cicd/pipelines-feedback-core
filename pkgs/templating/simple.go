package templating

import (
	"bytes"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
import "text/template"

func TemplateDashboardUrl(templateStr string, kubeObject v1.Object, typeMeta v1.TypeMeta) (string, error) {
	t, err := template.New("dashboard_url").Parse(templateStr)

	if err != nil {
		return "", errors.Wrap(err, "cannot create a dashboard URL. Please check your template")
	}
	var result bytes.Buffer
	execErr := t.Execute(&result, map[string]interface{}{
		"job":        kubeObject,
		"name":       kubeObject.GetName(),
		"namespace":  kubeObject.GetNamespace(),
		"kind":       typeMeta.Kind,
		"apiVersion": typeMeta.APIVersion,
		"apiGroup":   typeMeta.GroupVersionKind().Group,
		"gvk":        typeMeta.GroupVersionKind(),
	})
	if execErr != nil {
		return "", errors.Wrap(execErr, "cannot create a dashboard URL. Please check your template")
	}
	return result.String(), nil
}
