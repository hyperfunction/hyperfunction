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

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"

	functioninformer "github.com/hyperfunction/hyperfunction/pkg/client/injection/informers/core/v1alpha1/function"
	functionreconciler "github.com/hyperfunction/hyperfunction/pkg/client/injection/reconciler/core/v1alpha1/function"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	functionInformer := functioninformer.Get(ctx)

	r := &Reconciler{
		Tracker: nil,

		coreClientSet: kubeclient.Get(ctx),
	}
	impl := functionreconciler.NewImpl(ctx, r)
	r.Tracker = tracker.New(impl.EnqueueSlowKey, controller.GetTrackerLease(ctx))

	logger.Info("Setting up event handlers.")

	functionInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))
	return impl
}
