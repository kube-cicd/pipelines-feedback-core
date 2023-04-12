package v1alpha1

import (
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

	Spec   PFCSpec   `json:"spec"`
	Status PFCStatus `json:"status,omitempty"`
}

// NewPFConfig is making a new instance of a resource making sure that defaults will be respected
func NewPFConfig() PFConfig {
	return PFConfig{
		Spec: PFCSpec{},
	}
}

// PFCSpec represents .spec
type PFCSpec struct {
}

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
