package k8s

import (
	"github.com/kube-cicd/pipelines-feedback-core/pkgs/contract"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HasUsableAnnotations is checking if Kubernetes object is usable at all
func HasUsableAnnotations(meta metav1.ObjectMeta) (bool, error) {
	scm, err := CreateJobContextFromKubernetesAnnotations(meta)
	if err != nil {
		return false, err
	}
	return scm.IsValid(), nil
}

// CreateJobContextFromKubernetesAnnotations translates any Kubernetes object into contract.JobContext
func CreateJobContextFromKubernetesAnnotations(meta metav1.ObjectMeta) (contract.JobContext, error) {
	isTechnicalJob := false
	techJob := ""
	if val, exists := meta.Annotations[contract.GetTechnicalJobAnnotation()]; exists {
		logrus.Debugf("Has '%s'", contract.GetTechnicalJobAnnotation())
		isTechnicalJob = true
		techJob = val
	}

	repoHttpsUrl := ""
	if val, exists := meta.Annotations[contract.GetHttpsRepoUrlAnnotation()]; exists {
		logrus.Debugf("Has '%s'", contract.GetHttpsRepoUrlAnnotation())
		repoHttpsUrl = val
	}

	scm, err := contract.NewSCMContext(repoHttpsUrl)
	if err != nil && !isTechnicalJob {
		return scm, errors.Wrap(err, "cannot create JobContext")
	}

	scm.TechnicalJob = techJob

	if val, exists := meta.Annotations[contract.GetPrIdAnnotation()]; exists {
		logrus.Debugf("Has '%s'", contract.GetPrIdAnnotation())
		scm.PrId = val
	}
	if val, exists := meta.Annotations[contract.GetRefAnnotation()]; exists {
		logrus.Debugf("Has '%s'", contract.GetRefAnnotation())
		scm.Reference = val
	}
	if val, exists := meta.Annotations[contract.GetCommmitAnnotation()]; exists {
		logrus.Debugf("Has '%s'", contract.GetCommmitAnnotation())
		scm.Commit = val
	}
	return scm, nil
}
