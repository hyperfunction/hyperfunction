// Copyright 2022 The hyperfunction Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

type FunctionSpec struct {
	// +optional
	Engine  string       `json:"engine"`
	Runtime string       `json:"runtime"`
	Public  bool         `json:"public"`
	Code    FunctionCode `json:"code"`
}

type FunctionCode struct {
	Inline *string `json:"inline,omitempty"`
}

// FunctionStatus defines the observed state of Function
type FunctionStatus struct {
	duckv1.Status `json:",inline"`
}

// Function
//
// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Function struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the Function
	// +optional
	Spec FunctionSpec `json:"spec,omitempty"`

	// Status is the current state of the Function
	// +optional
	Status FunctionStatus `json:"status,omitempty"`
}

// Verify that Function adheres to the appropriate interfaces.
var (
	// Check that Function can be validated and defaulted.
	_ apis.Validatable   = (*Function)(nil)
	_ apis.Defaultable   = (*Function)(nil)
	_ kmeta.OwnerRefable = (*Function)(nil)
	// Check that Function conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Function)(nil)
)

// FunctionList is a list of Function resources.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Function `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the duckv1.KRShaped interface.
func (in *Function) GetStatus() *duckv1.Status {
	return &in.Status.Status
}
