package v1alpha1

import (
	"github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/contract"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=pfc

// +kubebuilder:subresource:status
type PFConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Spec      `json:"spec"`
	Data   Data      `json:"data"`
	Status PFCStatus `json:"status,omitempty"`
}

func (pfc *PFConfig) IsGlobalConfiguration() bool {
	return pfc.Namespace == ""
}

func (pfc *PFConfig) IsForPipeline(pipeline contract.PipelineInfo) bool {
	// no selector = matches all
	if len(pfc.Spec.JobDiscovery.LabelSelector) == 0 {
		return true
	}
	for _, labelSelector := range pfc.Spec.JobDiscovery.LabelSelector {
		selector, _ := metav1.LabelSelectorAsSelector(&labelSelector)
		if selector.Matches(pipeline.GetLabels()) {
			return true
		}
	}
	return false
}

func (pfc *PFConfig) HasLabelSelector() bool {
	return len(pfc.Spec.JobDiscovery.LabelSelector) > 0
}

// NewPFConfig is making a new instance of a resource making sure that defaults will be respected
func NewPFConfig() PFConfig {
	return PFConfig{
		Spec: Spec{
			JobDiscovery: JobDiscovery{
				LabelSelector: []metav1.LabelSelector{},
			},
		},
		Data:   Data{},
		Status: PFCStatus{},
	}
}

// Spec represents .spec
type Spec struct {
	JobDiscovery JobDiscovery `json:"jobDiscovery"`
}

// JobDiscovery represents .spec.jobDiscovery
type JobDiscovery struct {
	// .spec.jobDiscovery.labelSelector
	LabelSelector []metav1.LabelSelector `json:"labelSelector,omitempty"`
}

// Data represents similar field like "data" in ConfigMap, a simple key-value store
// with a convention of lowercase entries with dots as groups e.g. `scm.token-secret-name: "my-secret"`
type Data map[string]string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// PFConfigList represents a list
type PFConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PFConfig `json:"items"`
}

type PFCStatus struct {
}
