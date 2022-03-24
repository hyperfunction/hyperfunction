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
	"fmt"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"

	"github.com/hyperfunction/hyperfunction/pkg/apis/core/v1alpha1"
	functionreconciler "github.com/hyperfunction/hyperfunction/pkg/client/injection/reconciler/core/v1alpha1/function"
	"github.com/hyperfunction/hyperfunction/pkg/kube/controllers"
)

const (
	functionNameLabel = "core.hyperfunction.dev/function"
	functionCmKey     = "code"
)

// Reconciler implements functionreconciler for Function resources.
type Reconciler struct {
	// Tracker builds an index of what resources are watching other resources
	// so that we can immediately react to changes tracked resources.
	Tracker tracker.Interface

	coreClientSet kubernetes.Interface
}

var (
	// Check that our Reconciler implements Interface
	_ functionreconciler.Interface = (*Reconciler)(nil)
	// Check that our Reconciler implements Interface
	_ functionreconciler.Finalizer = (*Reconciler)(nil)
)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, f *v1alpha1.Function) pkgreconciler.Event {
	logger := logging.FromContext(ctx).
		With(zap.String("name", f.Name), zap.String("namespace", f.Namespace))
	ctx = logging.WithLogger(ctx, logger)

	if inlineCode := f.Spec.Code.Inline; inlineCode != nil {
		cm, err := r.reconcileCodeConfigmap(ctx, f)
		if err != nil {
			logger.Errorw("Error reconciling code configmap", zap.Error(err))
		}
		err = r.Tracker.TrackReference(tracker.Reference{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Name:       cm.Name,
			Namespace:  cm.Namespace,
		}, f)
		if err != nil {
			logger.Errorw("Error tracking code configmap", zap.Error(err))
			return err
		}
	}

	logger.Debug("Function reconciled")
	return nil
}

func (r *Reconciler) FinalizeKind(ctx context.Context, f *v1alpha1.Function) pkgreconciler.Event {
	// TODO(timonwong) Implements cleanup logic
	return nil
}

func (r *Reconciler) reconcileCodeConfigmap(ctx context.Context, f *v1alpha1.Function) (*corev1.ConfigMap, error) {
	logger := logging.FromContext(ctx)

	cmClient := r.coreClientSet.CoreV1().ConfigMaps(f.Namespace)
	cmList, err := cmClient.List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", functionNameLabel, f.Name),
	})
	if err != nil {
		return nil, err
	}

	expectedCm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-code-", f.Name),
			Namespace:    f.Namespace,
			Labels: map[string]string{
				functionNameLabel: f.Name,
			},
		},
		Data: map[string]string{
			functionCmKey: *f.Spec.Code.Inline,
		},
	}
	expectedCm.SetOwnerReferences([]metav1.OwnerReference{
		*kmeta.NewControllerRef(f),
	})

	if len(cmList.Items) == 0 {
		logger.Infow("Creating configmap")
		return cmClient.Create(ctx, expectedCm, metav1.CreateOptions{})
	}

	if cmLen := len(cmList.Items); cmLen > 1 {
		logger.Warnw("Deleting dangling configmaps", zap.Int("count", cmLen))
		for i := 1; i < cmLen; i++ {
			name := cmList.Items[i].Name
			if err := controllers.IgnoreNotFound(cmClient.Delete(ctx, name, metav1.DeleteOptions{})); err != nil {
				logger.Errorw("Error deleting dangling configmap", zap.String("cmName", name), zap.Error(err))
			}
		}
	}

	actualCm := &cmList.Items[0]
	if !reflect.DeepEqual(actualCm.Data, expectedCm.Data) {
		actualCm.Data = expectedCm.Data
		return cmClient.Update(ctx, actualCm, metav1.UpdateOptions{})
	}

	return actualCm, nil
}
