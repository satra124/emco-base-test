// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package resourcebundlestate

import (
	"context"
	"log"

	"gitlab.com/project-emco/core/emco-base/src/monitor/pkg/apis/k8splugin/v1alpha1"

	pkgerrors "github.com/pkg/errors"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// AddIngressController the new controller to the controller manager
func AddIngressController(mgr manager.Manager) error {
	return addIngressController(mgr, newIngressReconciler(mgr))
}

func addIngressController(mgr manager.Manager, r *ingressReconciler) error {
	// Create a new controller
	c, err := controller.New("Ingress-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to secondar resource Ingress
	// Predicate filters Ingress which don't have the k8splugin label
	err = c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForObject{}, &ingressPredicate{})
	if err != nil {
		return err
	}

	return nil
}

func newIngressReconciler(m manager.Manager) *ingressReconciler {
	return &ingressReconciler{client: m.GetClient()}
}

type ingressReconciler struct {
	client client.Client
}

// Reconcile implements the loop that will update the ResourceBundleState CR
// whenever we get any updates from all the ingress we watch.
func (r *ingressReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log.Printf("Updating ResourceBundleState for Ingress: %+v\n", req)

	ing := &v1beta1.Ingress{}
	err := r.client.Get(context.TODO(), req.NamespacedName, ing)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Printf("Ingress not found: %+v. Remove from CR if it is stored there.\n", req.NamespacedName)
			// Remove the Ingress's status from StatusList
			// This can happen if we get the DeletionTimeStamp event
			// after the Ingress has been deleted.
			DeleteFromAllCRs(r, req.NamespacedName)
			return reconcile.Result{}, nil
		}
		log.Printf("Failed to get ingress: %+v\n", req.NamespacedName)
		return reconcile.Result{}, err
	}
	err = UpdateCR(r, ing)
	if err != nil {
		// Requeue the update
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ingressReconciler) GetClient() client.Client {
	return r.client
}

func (r *ingressReconciler) DeleteObj(cr *v1alpha1.ResourceBundleState, name string) bool {
	var found bool
	length := len(cr.Status.IngressStatuses)
	for i, rstatus := range cr.Status.IngressStatuses {
		if rstatus.Name == name {
			found = true
			//Delete that status from the array
			cr.Status.IngressStatuses[i] = cr.Status.IngressStatuses[length-1]
			cr.Status.IngressStatuses[length-1].Status = v1beta1.IngressStatus{}
			cr.Status.IngressStatuses = cr.Status.IngressStatuses[:length-1]
			break
		}
	}
	return found
}

func (r *ingressReconciler) UpdateStatus(cr *v1alpha1.ResourceBundleState, obj runtime.Object) (bool, error) {

	var found bool
	ing, ok := obj.(*v1beta1.Ingress)
	if !ok {
		return found, pkgerrors.Errorf("Unknown resource %v", obj)
	}
	// Update status after searching for it in the list of resourceStatuses
	for i, rstatus := range cr.Status.IngressStatuses {
		// Look for the status if we already have it in the CR
		if rstatus.Name == ing.Name {
			ing.Status.DeepCopyInto(&cr.Status.IngressStatuses[i].Status)
			found = true
			break
		}
	}
	if !found {
		// Add it to CR
		c := v1beta1.Ingress{
			TypeMeta:   ing.TypeMeta,
			ObjectMeta: ing.ObjectMeta,
			Spec:       ing.Spec,
			Status:     ing.Status,
		}
		c.ObjectMeta.ManagedFields = []metav1.ManagedFieldsEntry{}
		cr.Status.IngressStatuses = append(cr.Status.IngressStatuses, c)
	}
	return found, nil
}
