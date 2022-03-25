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

package main

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/certificates"
	"knative.dev/pkg/webhook/resourcesemantics"
	"knative.dev/pkg/webhook/resourcesemantics/defaulting"
	"knative.dev/pkg/webhook/resourcesemantics/validation"

	corev1alpha1 "github.com/hyperfunction/hyperfunction/pkg/apis/core/v1alpha1"
)

func main() {
	webhookName := webhook.NameFromEnv()

	// Set up a signal context with our webhook options
	ctx := webhook.WithOptions(signals.NewContext(), webhook.Options{
		ServiceName: webhook.NameFromEnv(),
		Port:        webhook.PortFromEnv(8443),
		SecretName:  webhookName + "-certs",
	})

	sharedmain.WebhookMainWithContext(ctx, webhook.NameFromEnv(),
		certificates.NewController,
		newDefaultingAdmissionController,
		newValidationAdmissionController,
	)
}

var types = map[schema.GroupVersionKind]resourcesemantics.GenericCRD{
	corev1alpha1.SchemeGroupVersion.WithKind("Function"): &corev1alpha1.Function{},
}

var callbacks = map[schema.GroupVersionKind]validation.Callback{}

func newDefaultingAdmissionController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return defaulting.NewAdmissionController(ctx,
		// Name of the resource webhook.
		"webhook.hyperfunction.dev",
		// The path on which to serve the webhook.
		"/defaulting",
		// The resources to default.
		types,
		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		nil,
		// Whether to disallow unknown fields.
		true,
	)
}

func newValidationAdmissionController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	return validation.NewAdmissionController(ctx,
		// Name of the resource webhook.
		"validation.webhook.serving.knative.dev",
		// The path on which to serve the webhook.
		"/resource-validation",
		// The resources to validate.
		types,
		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		nil,
		// Whether to disallow unknown fields.
		true,
		// Extra validating callbacks to be applied to resources.
		callbacks,
	)
}
