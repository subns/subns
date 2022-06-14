/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	hierarchyv1alpha1 "github.com/subns/subns/api/v1alpha1"
)

const subnamespaceFinalizer = "subns" + hierarchyv1alpha1.Group + "/finalizer"

// SubnamespaceReconciler reconciles a Subnamespace object
type SubnamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hierarchy.subns.org,resources=subnamespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hierarchy.subns.org,resources=subnamespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hierarchy.subns.org,resources=subnamespaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Subnamespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *SubnamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	sn := &hierarchyv1alpha1.Subnamespace{}
	err := r.Client.Get(ctx, req.NamespacedName, sn)
	if err != nil {
		// object may have been deleted already
		return ctrl.Result{}, nil
	}

	// create an hns hierarchy configuration sn.Name under the sn.Parent
	err = createHNS(*sn)
	if err != nil {
		r.setStatusUnready(ctx, err, sn)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *SubnamespaceReconciler) setStatusUnready(ctx context.Context, err error, sn *hierarchyv1alpha1.Subnamespace) {
	c := metav1.Condition{
		Type:    HNSSubnamespaceCreated,
		Status:  metav1.ConditionFalse,
		Reason:  err.Error(),
		Message: "",
	}
	p := client.MergeFrom(sn.DeepCopy())
	meta.SetStatusCondition(&sn.Status.Conditions, c)
	r.Status().Patch(ctx, sn, p)
}

var ErrHNSCreateFailure = errors.New("errHNSCreateFailure")
var HNSSubnamespaceCreated = "HNSSubNamespaceCreated"

func createHNS(sn hierarchyv1alpha1.Subnamespace) error {
	return ErrHNSCreateFailure
}

// SetupWithManager sets up the controller with the Manager.
func (r *SubnamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hierarchyv1alpha1.Subnamespace{}).
		Complete(r)
}
