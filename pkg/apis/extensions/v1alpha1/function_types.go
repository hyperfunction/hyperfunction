package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmeta"
)

// Function
//
// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Function struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Spec FunctionSpec `json:"spec,omitempty"`

	// +optional
	Status FunctionStatus `json:"status"`
}

var (
	// Check that Function can be validated and defaulted.
	_ apis.Validatable   = (*Function)(nil)
	_ apis.Defaultable   = (*Function)(nil)
	_ kmeta.OwnerRefable = (*Function)(nil)
)

type FunctionSpec struct {
}

type FunctionStatus struct {
}

// FunctionList is a list of Function resources.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Function `json:"items"`
}
