/*
Copyright 2019 The Kubernetes Authors.

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

package experiment

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	experimentv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
	trialv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/trial/v1alpha2"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/experiment/clientset"
	"github.com/kubeflow/katib/pkg/controller/v1alpha3/recorder"
)

const (
	controllerName = "experiment-controller"
)

var log = logf.Log.WithName(controllerName)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Experiment Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	r := &ReconcileExperiment{
		Client:   mgr.GetClient(),
		Recorder: recorder.New(mgr.GetRecorder(controllerName)),
		scheme:   mgr.GetScheme(),
		Trial:    clientset.New(mgr.GetClient()),
	}
	r.updateStatusHandler = r.updateStatus
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Experiment
	err = c.Watch(&source.Kind{Type: &experimentv1alpha2.Experiment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &trialv1alpha2.Trial{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &experimentv1alpha2.Experiment{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileExperiment{}

// ReconcileExperiment reconciles a Experiment object
type ReconcileExperiment struct {
	client.Client
	recorder.Recorder
	scheme *runtime.Scheme
	clientset.Trial
	// To allow injection of updateStatus for testing.
	updateStatusHandler func(e *experimentv1alpha2.Experiment) error
}

// Reconcile reads that state of the cluster for a Experiment object and makes changes based on the state read
// and what is in the Experiment.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=experiments.kubeflow.org,resources=experiments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=experiments.kubeflow.org,resources=experiments/status,verbs=get;update;patch
func (r *ReconcileExperiment) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Experiment instance
	instance := &experimentv1alpha2.Experiment{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return r.handle(instance)
}
