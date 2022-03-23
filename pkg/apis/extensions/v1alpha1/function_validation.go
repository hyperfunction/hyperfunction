package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable
func (in *Function) Validate(ctx context.Context) *apis.FieldError {
	return in.Spec.Validate(ctx).ViaField("spec")
}

// Validate implements apis.Validatable
func (in *FunctionSpec) Validate(ctx context.Context) *apis.FieldError {
	return nil
}
