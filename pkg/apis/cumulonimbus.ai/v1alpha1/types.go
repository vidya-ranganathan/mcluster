package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

type Mcluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MclusterSpec   `json:"spec,omitempty"`
	Status MclusterStatus `json:"status,omitempty"`
}

type MclusterStatus struct {
	MclusterID string `json:"clusterID,omitempty"`
	// Progress string `json:"spec,omitempty"`
}

type MclusterSpec struct {
	Name string `json:"name,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MclusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Mcluster `json:"items,omitempty"`
}
