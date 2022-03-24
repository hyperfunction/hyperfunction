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

package function

import (
	"context"

	"k8s.io/client-go/kubernetes"
	pkgreconciler "knative.dev/pkg/reconciler"

	"github.com/hyperfunction/hyperfunction/pkg/apis/extensions/v1alpha1"
	functionreconciler "github.com/hyperfunction/hyperfunction/pkg/client/injection/reconciler/extensions/v1alpha1/function"
)

// Reconciler implements functionreconciler for Function resources.
type Reconciler struct {
	kubeclient kubernetes.Interface
}

// Check that our Reconciler implements Interface
var _ functionreconciler.Interface = (*Reconciler)(nil)

// Check that our Reconciler implements Interface
var _ functionreconciler.Finalizer = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *v1alpha1.Function) pkgreconciler.Event {
	panic("NYI")
}

func (r *Reconciler) FinalizeKind(ctx context.Context, o *v1alpha1.Function) pkgreconciler.Event {
	panic("NYI")
}
