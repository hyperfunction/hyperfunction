package v1alpha1

import "k8s.io/apimachinery/pkg/runtime/schema"

// GetGroupVersionKind implements kmeta.OwnerRefable
func (*Function) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("Function")
}
