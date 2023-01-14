package k8s

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HasUsableAnnotations is checking if Kubernetes object is usable at all
func HasUsableAnnotations(meta metav1.ObjectMeta) bool {
	return CreateSCMContextFromKubernetesAnnotations(meta).IsValid()
}

// CreateSCMContextFromKubernetesAnnotations translates any Kubernetes object into contract.SCMContext
func CreateSCMContextFromKubernetesAnnotations(meta metav1.ObjectMeta) contract.SCMContext {
	scm := contract.SCMContext{
		Commit:       "",
		Reference:    "",
		RepoHttpsUrl: "",
		PrId:         "",
	}
	if val, exists := meta.Annotations[contract.AnnotationPrId]; exists {
		scm.PrId = val
	}
	if val, exists := meta.Annotations[contract.AnnotationReference]; exists {
		scm.Reference = val
	}
	if val, exists := meta.Annotations[contract.AnnotationCommitHash]; exists {
		scm.Commit = val
	}
	if val, exists := meta.Annotations[contract.AnnotationHttpsRepo]; exists {
		scm.RepoHttpsUrl = val
	}
	return scm
}
